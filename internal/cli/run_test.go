package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/darpbr/darp-cli/internal/project/doctor"
)

func TestSyncRejectsUnknownArgumentsWithoutWriting(t *testing.T) {
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(old); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	}()
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"sync", "--unknown"}, &stdout, &stderr, ""); code != 1 {
		t.Fatalf("exit code=%d", code)
	}
	if _, err := os.Stat("darp.yml"); !os.IsNotExist(err) {
		t.Fatalf("unexpected manifest write: %v", err)
	}
}

func TestRunSyncDryRunShowsPreviewWithoutBootstrapWrites(t *testing.T) {
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(old); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	}()
	if err := os.MkdirAll(filepath.Join(".agents", "skills", "custom"), 0o755); err != nil {
		t.Fatal(err)
	}
	skill := "---\nname: custom\ndescription: A valid test skill for previewing sync behavior.\n---\n# Instructions\n"
	if err := os.WriteFile(filepath.Join(".agents", "skills", "custom", "SKILL.md"), []byte(skill), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"sync", "--dry-run"}, &stdout, &stderr, ""); code != 0 {
		t.Fatalf("exit code=%d, stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "DARP sync preview") || !strings.Contains(stdout.String(), "custom") {
		t.Fatalf("unexpected report: %s", stdout.String())
	}
	if _, err := os.Stat("darp.yml"); !os.IsNotExist(err) {
		t.Fatalf("dry run created darp.yml: %v", err)
	}
	if _, err := os.Stat(".darp"); !os.IsNotExist(err) {
		t.Fatalf("dry run created .darp: %v", err)
	}
}

func TestRunSyncReturnsSuccessForPartialDiscoveryWarnings(t *testing.T) {
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(old); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	}()
	if err := os.MkdirAll(filepath.Join(".agents", "skills", "invalid-name"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(".agents", "skills", "invalid-name", "SKILL.md"), []byte("---\nname: different\ndescription: Invalid identity.\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"sync", "--dry-run"}, &stdout, &stderr, ""); code != 0 {
		t.Fatalf("partial discovery should succeed: code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Ignored (1)") || !strings.Contains(stdout.String(), "Warnings (1)") {
		t.Fatalf("partial result was not reported: %s", stdout.String())
	}
}

func TestRunSyncReturnsFailureForInvalidManifest(t *testing.T) {
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(old); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	}()
	if err := os.WriteFile("darp.yml", []byte("skills: [\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"sync"}, &stdout, &stderr, ""); code != 1 {
		t.Fatalf("invalid manifest should fail: code=%d stdout=%s", code, stdout.String())
	}
	if !strings.Contains(stderr.String(), "invalid darp.yml") {
		t.Fatalf("invalid manifest failure not reported: %s", stderr.String())
	}
}

func TestRunInitInitializesProject(t *testing.T) {
	projectDir := t.TempDir()
	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir to temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previousDir)
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{"init"}, &stdout, &stderr, "# lifecycle\n")
	if exitCode != 0 {
		t.Fatalf("expected success, got exit code %d with stderr %q", exitCode, stderr.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	if !strings.Contains(stdout.String(), "Project initialized") {
		t.Fatalf("expected success output, got %q", stdout.String())
	}

	if _, err := os.Stat(filepath.Join(projectDir, "darp.yml")); err != nil {
		t.Fatalf("expected darp.yml to exist: %v", err)
	}
}

func TestRunWithoutCommandReturnsUsageError(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(nil, &stdout, &stderr, "# lifecycle\n")
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "DARP CLI - Developer AI Resource Platform") {
		t.Fatalf("expected usage output, got %q", stderr.String())
	}
}

func TestRunHelpVariants(t *testing.T) {
	testCases := [][]string{
		{"help"},
		{"--help"},
		{"-h"},
	}

	for _, args := range testCases {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			exitCode := Run(args, &stdout, &stderr, "# lifecycle\n")
			if exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d", exitCode)
			}

			if stderr.Len() != 0 {
				t.Fatalf("expected empty stderr, got %q", stderr.String())
			}

			if !strings.Contains(stdout.String(), "Useful commands:") {
				t.Fatalf("expected help output, got %q", stdout.String())
			}
		})
	}
}

func TestRunVersionVariants(t *testing.T) {
	testCases := [][]string{
		{"version"},
		{"--version"},
		{"-v"},
	}

	for _, args := range testCases {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			exitCode := RunWithVersion(args, &stdout, &stderr, "# lifecycle\n", "0.1.0")
			if exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d", exitCode)
			}

			if stderr.Len() != 0 {
				t.Fatalf("expected empty stderr, got %q", stderr.String())
			}

			if strings.TrimSpace(stdout.String()) != "darp version 0.1.0" {
				t.Fatalf("unexpected version output %q", stdout.String())
			}
		})
	}
}

func TestRunUnknownCommandShowsHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{"nope"}, &stdout, &stderr, "# lifecycle\n")
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), `unknown command "nope"`) {
		t.Fatalf("expected unknown command error, got %q", stderr.String())
	}

	if !strings.Contains(stderr.String(), "Useful commands:") {
		t.Fatalf("expected help text with unknown command, got %q", stderr.String())
	}
}

func TestRunDoctorReportsMissingProject(t *testing.T) {
	projectDir := t.TempDir()
	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(previousDir) })

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if exitCode := Run([]string{"doctor"}, &stdout, &stderr, "# lifecycle\n"); exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Running DARP Doctor...") || !strings.Contains(stdout.String(), "Errors:") {
		t.Fatalf("unexpected output: %q", stdout.String())
	}
}

func TestRenderDoctorCheckIncludesState(t *testing.T) {
	tests := []struct {
		name  string
		check doctor.CheckResult
		state string
	}{
		{name: "pass", check: doctor.CheckResult{Name: "Configuration", State: doctor.Pass}, state: "PASS"},
		{name: "warning", check: doctor.CheckResult{Name: "Version", State: doctor.Warning, Message: "old"}, state: "WARNING"},
		{name: "fail", check: doctor.CheckResult{Name: "Workflows", State: doctor.Fail, Message: "broken"}, state: "FAIL"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output := renderDoctorCheck(test.check); !strings.Contains(output, test.state) {
				t.Fatalf("expected %q in %q", test.state, output)
			}
		})
	}
}
