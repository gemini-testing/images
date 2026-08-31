package build

import "testing"

func TestShouldResolveChromeForTesting(t *testing.T) {
	tests := []struct {
		name       string
		sourceDir  string
		pkgSrcPath string
		version    string
		expected   bool
	}{
		{name: "major", sourceDir: "chrome/for-testing", version: "148", expected: true},
		{name: "full version", sourceDir: "chrome/for-testing", version: "148.0.7339.207"},
		{name: "apt source", sourceDir: "chrome/apt", version: "148"},
		{name: "local package", sourceDir: "chrome/for-testing", pkgSrcPath: "/tmp/google-chrome.deb", version: "148"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := shouldResolveChromeForTesting(test.sourceDir, test.pkgSrcPath, test.version)
			if actual != test.expected {
				t.Errorf("got %t, want %t", actual, test.expected)
			}
		})
	}
}

func TestResolveChromeForTestingVersion(t *testing.T) {
	versions := map[string]string{
		"147.0.7480.148": "wrong-major",
		"148.0.7339.80":  "older",
		"148.0.7339.207": "latest",
		"149.0.7480.12":  "newer-major",
	}

	actual, err := resolveChromeForTestingVersion("148", versions)
	if err != nil {
		t.Fatalf("resolve Chrome for Testing: %v", err)
	}
	if actual != "148.0.7339.207" {
		t.Errorf("got %s, want 148.0.7339.207", actual)
	}
}

func TestResolveChromeForTestingVersionMissingMajor(t *testing.T) {
	_, err := resolveChromeForTestingVersion("148", map[string]string{
		"149.0.7480.12": "newer-major",
	})
	if err == nil {
		t.Fatal("expected missing major error")
	}
}

func TestResolveChromeForTestingVersionInvalidMajor(t *testing.T) {
	_, err := resolveChromeForTestingVersion("148.0", nil)
	if err == nil {
		t.Fatal("expected invalid major error")
	}
}
