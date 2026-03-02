package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

var (
	listen       string
	upstream     string
	bidiUpstream string
)

func init() {
	flag.StringVar(&listen, "listen", ":4444", "Network address to accept connections")
	flag.StringVar(&upstream, "upstream", "127.0.0.1:4445", "Network address of inner WebDriver upstream")
	flag.StringVar(&bidiUpstream, "bidi-upstream", "127.0.0.1:9222", "Network address of Firefox BiDi websocket upstream")
}

func main() {
	flag.Parse()

	webdriverURL := &url.URL{Scheme: "http", Host: upstream}
	webdriverProxy := httputil.NewSingleHostReverseProxy(webdriverURL)
	defaultDirector := webdriverProxy.Director
	webdriverProxy.Director = func(req *http.Request) {
		host := req.Host
		scheme := "ws"
		if req.TLS != nil || strings.EqualFold(req.Header.Get("X-Forwarded-Proto"), "https") {
			scheme = "wss"
		}
		req.Header.Set("X-Bidi-Proxy-Host", host)
		req.Header.Set("X-Bidi-Proxy-Scheme", scheme)
		defaultDirector(req)
	}
	webdriverProxy.ModifyResponse = rewriteSessionWebSocketURL
	webdriverProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[HTTP_PROXY_ERROR] [%v]", err)
		w.WriteHeader(http.StatusBadGateway)
	}

	bidiURL := &url.URL{Scheme: "http", Host: bidiUpstream}
	bidiProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = bidiURL.Scheme
			req.URL.Host = bidiURL.Host
			req.Host = "localhost"
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("[BIDI_PROXY_ERROR] [%v]", err)
			w.WriteHeader(http.StatusBadGateway)
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isWebSocketUpgrade(r) && strings.HasPrefix(r.URL.Path, "/session/") {
			bidiProxy.ServeHTTP(w, r)
			return
		}
		webdriverProxy.ServeHTTP(w, r)
	})

	log.Printf("[INIT] [Listening on %s]", listen)
	log.Fatal(http.ListenAndServe(listen, handler))
}

func rewriteSessionWebSocketURL(resp *http.Response) error {
	req := resp.Request
	if req.Method != http.MethodPost || !isSessionCreatePath(req.URL.Path) {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	rewritten := body
	host := req.Header.Get("X-Bidi-Proxy-Host")
	scheme := req.Header.Get("X-Bidi-Proxy-Scheme")
	if host != "" {
		updated, err := replaceWebSocketURL(body, scheme, host)
		if err == nil {
			rewritten = updated
		}
	}

	resp.Body = io.NopCloser(bytes.NewReader(rewritten))
	resp.ContentLength = int64(len(rewritten))
	resp.Header.Set("Content-Length", strconv.Itoa(len(rewritten)))
	return nil
}

func replaceWebSocketURL(body []byte, scheme string, host string) ([]byte, error) {
	payload := map[string]interface{}{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	valueRaw, ok := payload["value"]
	if !ok {
		return nil, fmt.Errorf("no value field")
	}
	value, ok := valueRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value is not an object")
	}

	capsRaw, ok := value["capabilities"]
	if !ok {
		return nil, fmt.Errorf("no capabilities field")
	}
	caps, ok := capsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("capabilities is not an object")
	}

	if _, ok = caps["webSocketUrl"]; !ok {
		return nil, fmt.Errorf("no webSocketUrl field")
	}

	sessionID, _ := value["sessionId"].(string)
	if sessionID == "" {
		sessionID, _ = payload["sessionId"].(string)
	}
	if sessionID == "" {
		return nil, fmt.Errorf("no session id")
	}
	if scheme == "" {
		scheme = "ws"
	}

	caps["webSocketUrl"] = fmt.Sprintf("%s://%s/session/%s", scheme, host, sessionID)
	return json.Marshal(payload)
}

func isSessionCreatePath(path string) bool {
	return path == "/session" || path == "/wd/hub/session"
}

func isWebSocketUpgrade(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade") &&
		strings.EqualFold(r.Header.Get("Upgrade"), "websocket")
}
