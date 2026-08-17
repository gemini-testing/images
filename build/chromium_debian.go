package build

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	debianChromiumRelease = "trixie"
	metasnapURL           = "https://metasnap.debian.net/cgi-bin/api"
)

var chromiumMajorPattern = regexp.MustCompile(`^[0-9]+$`)

type debianChromiumPackage struct {
	Version   string
	Suite     string
	Component string
	Snapshot  string
}

type metasnapParameter struct {
	name  string
	value string
}

type DebianChromium struct {
	Requirements
	Architecture string
}

func (c *DebianChromium) Build() error {
	major := string(c.BrowserSource)
	if !chromiumMajorPattern.MatchString(major) {
		return fmt.Errorf("browser must be a Chromium major version when using --debian")
	}

	architecture, err := normalizeArchitecture(c.Architecture)
	if err != nil {
		return err
	}
	dockerArchitecture, err := currentDockerArchitecture()
	if err != nil {
		return err
	}
	if architecture != dockerArchitecture {
		return fmt.Errorf("requested %s, but the Docker daemon is %s; native builds only", architecture, dockerArchitecture)
	}

	pkg, err := resolveDebianChromium(context.Background(), nil, metasnapURL, major, architecture, debianChromiumRelease)
	if err != nil {
		return fmt.Errorf("resolve Debian Chromium: %v", err)
	}
	pkgTagVersion := extractVersion(pkg.Version)
	platform := "linux/" + architecture

	devDestDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create dev temporary dir: %v", err)
	}
	devImageTag := fmt.Sprintf("selenoid/dev_chromium:%s", pkgTagVersion)
	devImageRequirements := Requirements{NoCache: c.NoCache, Tags: []string{devImageTag}}
	devImage, err := NewImage("chromium/debian", devDestDir, devImageRequirements)
	if err != nil {
		return fmt.Errorf("init dev image: %v", err)
	}
	if err := copyDebianChromiumBaseFiles(devImage.Dir); err != nil {
		return fmt.Errorf("copy Debian Chromium base files: %v", err)
	}
	devImage.BuildArgs = []string{
		fmt.Sprintf("VERSION=%s", pkg.Version),
		fmt.Sprintf("DEBIAN_RELEASE=%s", debianChromiumRelease),
		fmt.Sprintf("DEBIAN_SUITE=%s", pkg.Suite),
		fmt.Sprintf("SNAPSHOT=%s", pkg.Snapshot),
	}
	if err := devImage.buildPlatform(platform); err != nil {
		return fmt.Errorf("build dev image: %v", err)
	}

	destDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create temporary dir: %v", err)
	}
	image, err := NewImage("chromium", destDir, c.Requirements)
	if err != nil {
		return fmt.Errorf("init image: %v", err)
	}
	image.BuildArgs = append(image.BuildArgs, fmt.Sprintf("VERSION=%s", pkgTagVersion))
	if err := image.buildPlatform(platform); err != nil {
		return fmt.Errorf("build image: %v", err)
	}
	if err := image.Test(c.TestsDir, "chrome", pkgTagVersion); err != nil {
		return fmt.Errorf("test image: %v", err)
	}
	if err := image.Push(); err != nil {
		return fmt.Errorf("push image: %v", err)
	}

	return nil
}

func copyDebianChromiumBaseFiles(destDir string) error {
	baseDir := filepath.Join(destDir, "base")
	for _, name := range []string{"xseld", "fileserver", "fluxbox", "aerokube.png"} {
		if _, err := copyEmbeddedFiles("/selenium/base", name, baseDir); err != nil {
			return err
		}
	}
	return nil
}

func normalizeArchitecture(architecture string) (string, error) {
	switch architecture {
	case "amd64", "x86_64":
		return "amd64", nil
	case "arm64", "aarch64":
		return "arm64", nil
	case "":
		return "", fmt.Errorf("--architecture is required with --debian")
	default:
		return "", fmt.Errorf("architecture must be amd64 or arm64")
	}
}

func currentDockerArchitecture() (string, error) {
	output, err := exec.Command("docker", "info", "--format", "{{.Architecture}}").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("get Docker architecture: %s", strings.TrimSpace(string(output)))
	}
	architecture, err := normalizeArchitecture(strings.TrimSpace(string(output)))
	if err != nil {
		return "", fmt.Errorf("unsupported Docker architecture: %s", strings.TrimSpace(string(output)))
	}
	return architecture, nil
}

func resolveDebianChromium(ctx context.Context, client *http.Client, baseURL string, major string, architecture string, release string) (debianChromiumPackage, error) {
	lines, err := queryMetasnap(ctx, client, baseURL, []metasnapParameter{
		{name: "archive", value: "debian"},
		{name: "pkg", value: "chromium"},
		{name: "arch", value: architecture},
	})
	if err != nil {
		return debianChromiumPackage{}, err
	}

	var pkg debianChromiumPackage
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		if strings.HasPrefix(fields[0], major+".") && (fields[1] == release || fields[1] == release+"-proposed-updates") {
			pkg.Version = fields[0]
			pkg.Suite = fields[1]
			pkg.Component = fields[2]
		}
	}
	if pkg.Version == "" {
		return debianChromiumPackage{}, fmt.Errorf("Debian %s has no Chromium %s package for %s", release, major, architecture)
	}

	pair := fmt.Sprintf("chromium=%s,chromium-driver=%s", pkg.Version, pkg.Version)
	lines, err = queryMetasnap(ctx, client, baseURL, []metasnapParameter{
		{name: "archive", value: "debian"},
		{name: "pkgs", value: pair},
		{name: "arch", value: architecture},
		{name: "suite", value: pkg.Suite},
		{name: "comp", value: pkg.Component},
	})
	if err != nil {
		return debianChromiumPackage{}, err
	}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == architecture {
			pkg.Snapshot = fields[1]
			break
		}
	}
	if pkg.Snapshot == "" {
		return debianChromiumPackage{}, fmt.Errorf("no snapshot contains %s for %s", pair, architecture)
	}

	return pkg, nil
}

func queryMetasnap(ctx context.Context, client *http.Client, baseURL string, parameters []metasnapParameter) ([]string, error) {
	endpoint, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse Metasnap URL: %v", err)
	}
	query := make([]string, 0, len(parameters))
	for _, parameter := range parameters {
		query = append(query, url.QueryEscape(parameter.name)+"="+url.QueryEscape(parameter.value))
	}
	endpoint.RawQuery = strings.Join(query, "&")

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create Metasnap request: %v", err)
	}
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("query Metasnap: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query Metasnap: unexpected status %s", response.Status)
	}

	var lines []string
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read Metasnap response: %v", err)
	}
	return lines, nil
}
