package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/darpbr/darp-cli/internal/project/doctor"
	projectinit "github.com/darpbr/darp-cli/internal/project/init"
	projectsync "github.com/darpbr/darp-cli/internal/project/sync"
)

const (
	cliName        = "DARP CLI"
	cliDescription = "Developer AI Resource Platform"
	defaultVersion = "dev"
)

type commandHelp struct {
	Name    string
	Summary string
}

var commandCatalog = []commandHelp{
	{Name: "init", Summary: "Initialize a DARP project"},
	{Name: "doctor", Summary: "Diagnose DARP project integrity"},
	{Name: "sync", Summary: "Reconcile the DARP asset inventory"},
	{Name: "help", Summary: "Show CLI help"},
	{Name: "version", Summary: "Show CLI version"},
}

// Run executes the CLI command flow and returns a process exit code.
func Run(args []string, stdout io.Writer, stderr io.Writer, lifecycleContent string) int {
	return RunWithVersion(args, stdout, stderr, lifecycleContent, defaultVersion)
}

// RunWithVersion executes the CLI command flow with the provided version string.
func RunWithVersion(args []string, stdout io.Writer, stderr io.Writer, lifecycleContent string, version string) int {
	if len(args) == 0 {
		if !writef(stderr, "%s", renderHelp()) {
			return 1
		}
		return 1
	}

	switch args[0] {
	case "help", "-h", "--help":
		if !writef(stdout, "%s", renderHelp()) {
			return 1
		}
		return 0
	case "version", "-v", "--version":
		if !writef(stdout, "%s\n", renderVersion(version)) {
			return 1
		}
		return 0
	case "init":
		service := projectinit.NewService(projectinit.NewOSFileSystem(), lifecycleContent)
		result, err := service.Initialize(".")
		if err != nil {
			if !writef(stderr, "error: %v\n", err) {
				return 1
			}
			return 1
		}

		for _, message := range result.Messages {
			if !writef(stdout, "%s\n", message) {
				return 1
			}
		}

		return 0
	case "doctor":
		result := doctor.NewService().Diagnose(".")
		if !writef(stdout, "Running DARP Doctor...\n\n") {
			return 1
		}
		for _, check := range result.Checks {
			if !writef(stdout, "%s\n", renderDoctorCheck(check)) {
				return 1
			}
		}
		passed, warnings, errors := result.Counts()
		if !writef(stdout, "\nSummary\n\nPassed: %d\nWarnings: %d\nErrors: %d\n\n", passed, warnings, errors) {
			return 1
		}
		if errors > 0 {
			if !writef(stdout, "%s\n", "Project has errors.") {
				return 1
			}
		} else if warnings > 0 {
			if !writef(stdout, "%s\n", "Project healthy with warnings.") {
				return 1
			}
		} else {
			if !writef(stdout, "%s\n", "Project healthy.") {
				return 1
			}
		}
		return result.ExitCode()
	case "sync":
		dryRun := false
		if len(args) > 2 || (len(args) == 2 && args[1] != "--dry-run") {
			if !writef(stderr, "error: unknown sync argument\n") {
				return 1
			}
			return 1
		}
		if len(args) == 2 {
			dryRun = true
		}
		result, err := (projectsync.Service{Root: "."}).Reconcile(dryRun)
		if err != nil {
			if !writef(stderr, "error: %v\n", err) {
				return 1
			}
			return 1
		}
		if !writef(stdout, "%s", renderSync(result, dryRun)) {
			return 1
		}
		return 0
	default:
		if !writef(stderr, "unknown command %q\n\n%s", args[0], renderHelp()) {
			return 1
		}
		return 1
	}
}

func renderSync(result projectsync.Result, dryRun bool) string {
	var b strings.Builder
	if dryRun {
		b.WriteString("DARP sync preview\n")
	} else {
		b.WriteString("DARP sync\n")
	}
	for _, group := range []struct {
		name  string
		items []projectsync.Item
	}{{"Added", result.Added}, {"Removed", result.Removed}, {"Kept", result.Kept}, {"Ignored", result.Ignored}, {"Warnings", result.Warnings}} {
		fmt.Fprintf(&b, "\n%s (%d)\n", group.name, len(group.items))
		for _, item := range group.items {
			identity := item.Identity
			if item.Detail != "" {
				identity += " — " + item.Detail
			}
			fmt.Fprintf(&b, "  %s\n", identity)
		}
	}
	return b.String()
}

func writef(writer io.Writer, format string, args ...any) bool {
	_, err := fmt.Fprintf(writer, format, args...)
	return err == nil
}

func renderDoctorCheck(check doctor.CheckResult) string {
	switch check.State {
	case doctor.Pass:
		return "✔ PASS " + check.Name
	case doctor.Warning:
		return "⚠ WARNING " + check.Message
	default:
		return "✖ FAIL " + check.Message
	}
}

func renderHelp() string {
	var builder strings.Builder

	builder.WriteString(cliName)
	builder.WriteString(" - ")
	builder.WriteString(cliDescription)
	builder.WriteString("\n\n")
	builder.WriteString("Manage reusable AI assets and DARP project structure from the command line.\n\n")
	builder.WriteString("Usage:\n")
	builder.WriteString("  darp <command> [arguments]\n\n")
	builder.WriteString("Useful commands:\n")

	for _, command := range commandCatalog {
		_, _ = fmt.Fprintf(&builder, "  %-8s %s\n", command.Name, command.Summary)
	}

	builder.WriteString("\nOptions:\n")
	builder.WriteString("  -h, --help       Show help\n")
	builder.WriteString("  -v, --version    Show version\n")

	return builder.String()
}

func renderVersion(version string) string {
	return fmt.Sprintf("darp version %s", version)
}
