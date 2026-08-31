package build

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func metasnapResponse(body string) *http.Response {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestResolveDebianChromium(t *testing.T) {
	for _, architecture := range []string{"amd64", "arm64"} {
		t.Run(architecture, func(t *testing.T) {
			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					query := r.URL.Query()
					if query.Get("arch") != architecture {
						t.Errorf("unexpected architecture: %s", query.Get("arch"))
					}
					switch {
					case query.Get("pkg") == "chromium":
						if r.URL.RawQuery != "archive=debian&pkg=chromium&arch="+architecture {
							t.Errorf("unexpected parameter order: %s", r.URL.RawQuery)
						}
						return metasnapResponse(strings.Join([]string{
							"125.0.0-1 trixie main 20240101T000000Z 20240102T000000Z",
							"126.0.6478.100-1 trixie main 20240601T000000Z 20240602T000000Z",
							"126.0.6478.126-1~deb13u1 trixie-proposed-updates main 20240716T000000Z 20240717T000000Z",
						}, "\n")), nil
					case query.Get("pkgs") == "chromium=126.0.6478.126-1~deb13u1,chromium-driver=126.0.6478.126-1~deb13u1":
						if !strings.HasPrefix(r.URL.RawQuery, "archive=debian&pkgs=") {
							t.Errorf("unexpected parameter order: %s", r.URL.RawQuery)
						}
						return metasnapResponse(fmt.Sprintf("%s 20240717T083710Z\n", architecture)), nil
					default:
						return nil, fmt.Errorf("unexpected query: %s", r.URL.RawQuery)
					}
				}),
			}

			pkg, err := resolveDebianChromium(context.Background(), client, "https://metasnap.test/api", "126", architecture, "trixie")
			if err != nil {
				t.Fatalf("resolve Debian Chromium: %v", err)
			}
			if pkg.Version != "126.0.6478.126-1~deb13u1" {
				t.Errorf("unexpected version: %s", pkg.Version)
			}
			if pkg.Suite != "trixie-proposed-updates" {
				t.Errorf("unexpected suite: %s", pkg.Suite)
			}
			if pkg.Component != "main" {
				t.Errorf("unexpected component: %s", pkg.Component)
			}
			if pkg.Snapshot != "20240717T083710Z" {
				t.Errorf("unexpected snapshot: %s", pkg.Snapshot)
			}
		})
	}
}

func TestResolveDebianChromiumMissingMajor(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return metasnapResponse("126.0.6478.126-1~deb13u1 trixie main 20240716T000000Z 20240717T000000Z\n"), nil
	})}

	_, err := resolveDebianChromium(context.Background(), client, "https://metasnap.test/api", "127", "arm64", "trixie")
	if err == nil {
		t.Fatal("expected missing major error")
	}
}

func TestNormalizeArchitecture(t *testing.T) {
	for input, expected := range map[string]string{
		"amd64":   "amd64",
		"x86_64":  "amd64",
		"arm64":   "arm64",
		"aarch64": "arm64",
	} {
		actual, err := normalizeArchitecture(input)
		if err != nil {
			t.Fatalf("normalize %s: %v", input, err)
		}
		if actual != expected {
			t.Errorf("normalize %s: got %s, want %s", input, actual, expected)
		}
	}

	if _, err := normalizeArchitecture("386"); err == nil {
		t.Fatal("expected unsupported architecture error")
	}
}
