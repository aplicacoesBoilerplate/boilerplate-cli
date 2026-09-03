package boilerplate

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeService struct {
	calls []string
	last  any
	err   error
}

func (f *fakeService) Auth(_ context.Context, request AuthRequest) error {
	f.calls = append(f.calls, "auth:"+string(request.Action))
	f.last = request
	return f.err
}

func (f *fakeService) Init(_ context.Context, request InitRequest) error {
	f.calls = append(f.calls, "init")
	f.last = request
	return f.err
}

func (f *fakeService) New(_ context.Context, request NewRequest) error {
	f.calls = append(f.calls, "new")
	f.last = request
	return f.err
}

func (f *fakeService) Add(_ context.Context, request AddRequest) error {
	f.calls = append(f.calls, "add")
	f.last = request
	return f.err
}

func (f *fakeService) Update(_ context.Context, request UpdateRequest) error {
	f.calls = append(f.calls, "update")
	f.last = request
	return f.err
}

func (f *fakeService) Doctor(_ context.Context, request DoctorRequest) error {
	f.calls = append(f.calls, "doctor")
	f.last = request
	return f.err
}

func (f *fakeService) Audit(_ context.Context, request AuditRequest) error {
	f.calls = append(f.calls, "audit")
	f.last = request
	return f.err
}

func TestNewRootCommandExposesStableTreeAndExecutionFlags(t *testing.T) {
	root := NewRootCommand(Dependencies{Service: &fakeService{}})

	if !root.SilenceUsage || !root.SilenceErrors {
		t.Fatal("root command must silence Cobra usage and duplicate errors")
	}
	for _, flag := range []string{"root", "dry-run"} {
		if root.PersistentFlags().Lookup(flag) == nil {
			t.Fatalf("missing persistent flag %q", flag)
		}
	}

	want := []string{"add", "audit", "auth", "completion", "doctor", "help", "init", "new", "update"}
	var got []string
	for _, command := range root.Commands() {
		got = append(got, command.Name())
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("commands = %v, want %v", got, want)
	}
}

func TestRunCLIMapsRequestsWithoutLeakingInfrastructureIntoCommands(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "pom.xml"), "<project/>")
	mustWrite(t, filepath.Join(root, "web", "package.json"), `{ "name": "web" }`)

	tests := []struct {
		name     string
		args     []string
		wantCall string
		assert   func(*testing.T, any)
	}{
		{
			name: "auth login",
			args: []string{"auth", "login"}, wantCall: "auth:login",
			assert: func(t *testing.T, value any) {
				request := value.(AuthRequest)
				if request.DryRun {
					t.Fatal("dry-run unexpectedly enabled")
				}
			},
		},
		{
			name: "init monorepo dry-run",
			args: []string{"--root", root, "--dry-run", "init"}, wantCall: "init",
			assert: func(t *testing.T, value any) {
				request := value.(InitRequest)
				if !request.DryRun || len(request.Workspace.JavaProjects) != 1 || len(request.Workspace.VueProjects) != 1 {
					t.Fatalf("unexpected init request: %#v", request)
				}
			},
		},
		{
			name: "new private by default",
			args: []string{"--root", root, "new", "java", "orders", "--owner", "acme", "--directory", "generated"}, wantCall: "new",
			assert: func(t *testing.T, value any) {
				request := value.(NewRequest)
				if request.Platform != PlatformJava || request.Name != "orders" || request.Owner != "acme" || request.Visibility != VisibilityPrivate {
					t.Fatalf("unexpected new request: %#v", request)
				}
				if request.Directory != filepath.Join(root, "generated") {
					t.Fatalf("directory = %q", request.Directory)
				}
			},
		},
		{
			name: "add scoped package",
			args: []string{"--root", root, "add", "vue", "@aplicacoesBoilerplate/ui@1.2.3"}, wantCall: "add",
			assert: func(t *testing.T, value any) {
				request := value.(AddRequest)
				if request.Platform != PlatformVue || request.Package != "@aplicacoesBoilerplate/ui" || request.Version != "1.2.3" {
					t.Fatalf("unexpected add request: %#v", request)
				}
			},
		},
		{
			name: "update defaults to all",
			args: []string{"--root", root, "update"}, wantCall: "update",
			assert: func(t *testing.T, value any) {
				if request := value.(UpdateRequest); request.Platform != PlatformAll {
					t.Fatalf("platform = %q", request.Platform)
				}
			},
		},
		{name: "doctor", args: []string{"--root", root, "doctor"}, wantCall: "doctor"},
		{
			name: "audit defaults to stdout",
			args: []string{"--root", root, "audit"}, wantCall: "audit",
			assert: func(t *testing.T, value any) {
				request := value.(AuditRequest)
				if request.Format != AuditFormatText || request.Output != "" {
					t.Fatalf("unexpected audit defaults: %#v", request)
				}
			},
		},
		{
			name: "audit json",
			args: []string{"--root", root, "--dry-run", "audit", "--format", "json", "--output", "reports/audit.json"}, wantCall: "audit",
			assert: func(t *testing.T, value any) {
				request := value.(AuditRequest)
				if !request.DryRun || request.Format != AuditFormatJSON || request.Output != filepath.Join(root, "reports", "audit.json") {
					t.Fatalf("unexpected audit request: %#v", request)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &fakeService{}
			var stdout, stderr bytes.Buffer
			code := RunCLI(context.Background(), test.args, Dependencies{
				Service: service,
				Stdout:  &stdout,
				Stderr:  &stderr,
			})
			if code != int(ExitSuccess) {
				t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
			}
			if strings.Join(service.calls, ",") != test.wantCall {
				t.Fatalf("calls = %v, want %q", service.calls, test.wantCall)
			}
			if test.assert != nil {
				test.assert(t, service.last)
			}
		})
	}
}

func TestRunCLIUsesStableExitCodesAndPrintsErrorOnce(t *testing.T) {
	tests := []struct {
		name string
		args []string
		err  error
		want ExitCode
	}{
		{name: "invalid platform", args: []string{"new", "mobile", "app"}, want: ExitUsage},
		{name: "missing args", args: []string{"add", "java"}, want: ExitUsage},
		{name: "unknown command", args: []string{"destroy"}, want: ExitUsage},
		{name: "typed auth error", args: []string{"auth", "status"}, err: NewCLIError(ExitAuthentication, "autenticacao necessaria", errors.New("sensitive cause")), want: ExitAuthentication},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &fakeService{err: test.err}
			var stderr bytes.Buffer
			code := RunCLI(context.Background(), test.args, Dependencies{Service: service, Stderr: &stderr})
			if code != int(test.want) {
				t.Fatalf("exit = %d, want %d, stderr = %q", code, test.want, stderr.String())
			}
			if strings.Contains(stderr.String(), "Usage:") || strings.Count(strings.TrimSpace(stderr.String()), "\n") > 0 {
				t.Fatalf("error must be printed once without usage: %q", stderr.String())
			}
			if strings.Contains(stderr.String(), "sensitive cause") {
				t.Fatalf("internal cause leaked: %q", stderr.String())
			}
		})
	}
}

func TestUnavailableDependencyNeverReportsTODOAsSuccess(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "pom.xml"), "<project/>")
	var stdout, stderr bytes.Buffer

	code := RunCLI(context.Background(), []string{"--root", root, "init"}, Dependencies{Stdout: &stdout, Stderr: &stderr})

	if code != int(ExitConfiguration) {
		t.Fatalf("exit = %d, want %d", code, ExitConfiguration)
	}
	if strings.Contains(strings.ToUpper(stdout.String()), "TODO") || strings.TrimSpace(stdout.String()) != "" {
		t.Fatalf("unexpected successful placeholder output: %q", stdout.String())
	}
}

func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
