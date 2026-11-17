package test

import (
	"os"
	"strings"
	"testing"
)

func readInstaller(t *testing.T, name string) string {
	t.Helper()
	content, err := os.ReadFile("../scripts/" + name)
	if err != nil {
		t.Fatalf("read installer %s: %v", name, err)
	}
	return string(content)
}

func TestPowerShellInstallerFailClosedContract(t *testing.T) {
	script := readInstaller(t, "install.ps1")
	required := []string{
		"releases/latest",
		"releases/tags/$tag",
		"/commits/$tag",
		"windows_amd64.zip",
		"checksums.txt.sigstore.json",
		"verify-blob",
		"Get-Sha256",
		"bootstrap helper SHA-256 mismatch",
		"unsupported architecture",
		"LOCALAPPDATA",
		"[Environment]::GetEnvironmentVariable(\"Path\", \"User\")",
	}
	for _, value := range required {
		if !strings.Contains(script, value) {
			t.Errorf("PowerShell installer is missing %q", value)
		}
	}
	if strings.Contains(script, "http://") {
		t.Error("PowerShell installer contains an insecure HTTP URL")
	}
}

func TestShellInstallerPlatformAndFailClosedContract(t *testing.T) {
	script := readInstaller(t, "install.sh")
	required := []string{
		"releases/latest",
		"releases/tags/$TAG",
		"/commits/$TAG",
		"linux/amd64",
		"linux/arm64",
		"darwin/amd64",
		"darwin/arm64",
		"checksums.txt.sigstore.json",
		"verify-blob",
		"artifact SHA-256 mismatch",
		"bootstrap helper SHA-256 mismatch",
		"$HOME/.local/bin",
		".profile.darp.XXXXXX",
		"could not update the current-user profile; the previous installation was restored",
		"unsupported platform",
	}
	for _, value := range required {
		if !strings.Contains(script, value) {
			t.Errorf("Shell installer is missing %q", value)
		}
	}
	if strings.Contains(script, "http://") {
		t.Error("Shell installer contains an insecure HTTP URL")
	}
	if strings.Contains(script, "/git/tags/") || strings.Contains(script, "/git/ref/tags/") {
		t.Error("Shell installer must resolve annotated and lightweight tags through the commits endpoint")
	}
}
