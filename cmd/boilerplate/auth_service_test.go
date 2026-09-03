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

const testToken = "test-secret-1234567890<&"

type fakeProcessRunner struct {
	responses map[string]string
	errors    map[string]error
	calls     []string
}

func (f *fakeProcessRunner) Run(_ context.Context, name string, args ...string) (string, error) {
	key := name + " " + strings.Join(args, " ")
	f.calls = append(f.calls, key)
	if err := f.errors[key]; err != nil {
		return "", err
	}
	return f.responses[key], nil
}

func TestAuthLoginPreservesFixturesCreatesBackupAndIsIdempotent(t *testing.T) {
	fixtures := []struct {
		name string
		xml  string
		npm  string
	}{
		{
			name: "windows",
			xml:  "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n<settings>\r\n  <!-- keep-windows -->\r\n  <mirrors><mirror><id>corp</id></mirror></mirrors>\r\n  <servers><server><id>other</id><username>keep</username></server></servers>\r\n</settings>\r\n",
			npm:  "registry=https://registry.npmjs.org\r\n# keep-windows\r\n",
		},
		{
			name: "linux",
			xml:  "<?xml version=\"1.0\"?>\n<settings>\n  <!-- keep-linux -->\n  <profiles><profile><id>corp</id></profile></profiles>\n</settings>\n",
			npm:  "fund=false\n# keep-linux\n",
		},
		{
			name: "darwin",
			xml:  "<settings xmlns=\"http://maven.apache.org/SETTINGS/1.2.0\">\n  <!-- keep-darwin -->\n  <servers>\n    <server><id>github-boilerplate</id><username>old</username><password>old</password></server>\n  </servers>\n</settings>\n",
			npm:  "@other:registry=https://example.invalid\n//example.invalid/:_authToken=keep\n# keep-darwin\n",
		},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			home := t.TempDir()
			settingsPath := filepath.Join(home, ".m2", "settings.xml")
			npmPath := filepath.Join(home, ".npmrc")
			mustWrite(t, settingsPath, fixture.xml)
			mustWrite(t, npmPath, fixture.npm)
			runner := authenticatedRunner()
			var output bytes.Buffer
			service := newAuthService(runner, func() (string, error) { return home, nil }, &output)

			if err := service.Auth(context.Background(), AuthRequest{Action: AuthLogin}); err != nil {
				t.Fatal(err)
			}

			settings := mustRead(t, settingsPath)
			npm := mustRead(t, npmPath)
			if !strings.Contains(settings, "keep-"+fixture.name) || !strings.Contains(npm, "keep-"+fixture.name) {
				t.Fatalf("fixture content was not preserved\nsettings=%s\nnpm=%s", settings, npm)
			}
			if strings.Count(settings, "<id>github-boilerplate</id>") != 1 || !strings.Contains(settings, "test-secret-1234567890&lt;&amp;") {
				t.Fatalf("Maven credentials not written exactly once: %s", settings)
			}
			if strings.Count(npm, "@aplicacoesBoilerplate:registry=") != 1 || strings.Count(npm, "//npm.pkg.github.com/:_authToken=") != 1 || !strings.Contains(npm, testToken) {
				t.Fatalf("npm credentials not written exactly once: %s", npm)
			}
			if strings.Contains(fixture.xml, "\r\n") && strings.Contains(strings.ReplaceAll(settings, "\r\n", ""), "\n") {
				t.Fatal("Windows newline convention was not preserved")
			}
			if got := mustRead(t, settingsPath+backupSuffix); got != fixture.xml {
				t.Fatalf("settings backup differs from original: %q", got)
			}
			if got := mustRead(t, npmPath+backupSuffix); got != fixture.npm {
				t.Fatalf("npm backup differs from original: %q", got)
			}

			if err := service.Auth(context.Background(), AuthRequest{Action: AuthLogin}); err != nil {
				t.Fatal(err)
			}
			if got := mustRead(t, settingsPath+backupSuffix); got != fixture.xml {
				t.Fatal("idempotent login rewrote the settings backup")
			}
			if got := mustRead(t, npmPath+backupSuffix); got != fixture.npm {
				t.Fatal("idempotent login rewrote the npm backup")
			}
			if strings.Contains(output.String(), testToken) {
				t.Fatal("token leaked to output")
			}
		})
	}
}

func TestAuthDryRunDoesNotCreateFilesOrDirectories(t *testing.T) {
	home := filepath.Join(t.TempDir(), "new-home")
	runner := authenticatedRunner()
	var output bytes.Buffer
	service := newAuthService(runner, func() (string, error) { return home, nil }, &output)

	if err := service.Auth(context.Background(), AuthRequest{Action: AuthLogin, DryRun: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(home); !os.IsNotExist(err) {
		t.Fatalf("dry-run created home or files: %v", err)
	}
	if strings.Contains(output.String(), testToken) {
		t.Fatal("token leaked during dry-run")
	}
}

func TestAuthLogoutRemovesOnlyManagedEntriesAndKeepsBackup(t *testing.T) {
	home := t.TempDir()
	settingsPath := filepath.Join(home, ".m2", "settings.xml")
	npmPath := filepath.Join(home, ".npmrc")
	mustWrite(t, settingsPath, "<settings><mirrors><!-- keep --></mirrors></settings>\n")
	mustWrite(t, npmPath, "fund=false\n")
	var output bytes.Buffer
	service := newAuthService(authenticatedRunner(), func() (string, error) { return home, nil }, &output)
	if err := service.Auth(context.Background(), AuthRequest{Action: AuthLogin}); err != nil {
		t.Fatal(err)
	}
	configuredSettings := mustRead(t, settingsPath)
	configuredNPM := mustRead(t, npmPath)

	if err := service.Auth(context.Background(), AuthRequest{Action: AuthLogout}); err != nil {
		t.Fatal(err)
	}

	settings := mustRead(t, settingsPath)
	npm := mustRead(t, npmPath)
	if !strings.Contains(settings, "<!-- keep -->") || strings.Contains(settings, "github-boilerplate") || strings.Contains(settings, testToken) {
		t.Fatalf("unexpected settings after logout: %s", settings)
	}
	if !strings.Contains(npm, "fund=false") || strings.Contains(npm, "boilerplate-cli") || strings.Contains(npm, "npm.pkg.github.com") || strings.Contains(npm, testToken) {
		t.Fatalf("unexpected npmrc after logout: %s", npm)
	}
	if mustRead(t, settingsPath+backupSuffix) != configuredSettings || mustRead(t, npmPath+backupSuffix) != configuredNPM {
		t.Fatal("logout backup must contain the immediately previous configuration")
	}
	if strings.Contains(output.String(), testToken) {
		t.Fatal("token leaked during logout")
	}
}

func TestAuthStatusNeverPrintsTokenAndRequiresBothConfigurations(t *testing.T) {
	home := t.TempDir()
	var output bytes.Buffer
	runner := authenticatedRunner()
	service := newAuthService(runner, func() (string, error) { return home, nil }, &output)
	if err := service.Auth(context.Background(), AuthRequest{Action: AuthLogin}); err != nil {
		t.Fatal(err)
	}
	if err := service.Auth(context.Background(), AuthRequest{Action: AuthStatus}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), testToken) {
		t.Fatal("status leaked token")
	}

	if err := os.Remove(filepath.Join(home, ".npmrc")); err != nil {
		t.Fatal(err)
	}
	err := service.Auth(context.Background(), AuthRequest{Action: AuthStatus})
	if ExitCodeFor(err) != ExitConfiguration {
		t.Fatalf("status without npmrc error = %v, code = %d", err, ExitCodeFor(err))
	}
}

func TestAuthSanitizesRunnerFailureAndRejectsInvalidXML(t *testing.T) {
	home := t.TempDir()
	runner := authenticatedRunner()
	runner.errors["gh auth token --hostname github.com"] = errors.New("failed with " + testToken)
	var output bytes.Buffer
	service := newAuthService(runner, func() (string, error) { return home, nil }, &output)

	err := service.Auth(context.Background(), AuthRequest{Action: AuthLogin})
	if ExitCodeFor(err) != ExitAuthentication || strings.Contains(err.Error(), testToken) || strings.Contains(output.String(), testToken) {
		t.Fatalf("runner failure was not sanitized: err=%v output=%q", err, output.String())
	}

	runner = authenticatedRunner()
	settingsPath := filepath.Join(home, ".m2", "settings.xml")
	mustWrite(t, settingsPath, "<settings><servers>")
	service = newAuthService(runner, func() (string, error) { return home, nil }, &output)
	err = service.Auth(context.Background(), AuthRequest{Action: AuthLogin})
	if ExitCodeFor(err) != ExitConfiguration {
		t.Fatalf("invalid XML error = %v, code = %d", err, ExitCodeFor(err))
	}
	if _, statErr := os.Stat(settingsPath + backupSuffix); !os.IsNotExist(statErr) {
		t.Fatal("invalid XML must not create a backup or mutate the file")
	}
}

func TestAuthRejectsMalformedNPMManagedBlockWithoutMutation(t *testing.T) {
	home := t.TempDir()
	settingsPath := filepath.Join(home, ".m2", "settings.xml")
	npmPath := filepath.Join(home, ".npmrc")
	mustWrite(t, settingsPath, "<settings></settings>\n")
	malformed := "fund=false\n" + npmManagedBegin + "\nkeep=true\n"
	mustWrite(t, npmPath, malformed)
	service := newAuthService(authenticatedRunner(), func() (string, error) { return home, nil }, &bytes.Buffer{})

	err := service.Auth(context.Background(), AuthRequest{Action: AuthLogin})
	if ExitCodeFor(err) != ExitConflict {
		t.Fatalf("malformed npmrc error = %v, code = %d", err, ExitCodeFor(err))
	}
	if got := mustRead(t, npmPath); got != malformed {
		t.Fatalf("malformed npmrc was changed: %q", got)
	}
	if _, statErr := os.Stat(npmPath + backupSuffix); !os.IsNotExist(statErr) {
		t.Fatal("malformed npmrc must not create a backup")
	}
}

func TestAuthDoesNotPartiallyUpdateWhenASecondBackupFails(t *testing.T) {
	home := t.TempDir()
	settingsPath := filepath.Join(home, ".m2", "settings.xml")
	npmPath := filepath.Join(home, ".npmrc")
	settingsBefore := "<settings><!-- original --></settings>\n"
	npmBefore := "fund=false\n"
	mustWrite(t, settingsPath, settingsBefore)
	mustWrite(t, npmPath, npmBefore)
	mustWrite(t, settingsPath+backupSuffix, "previous settings backup")
	if err := os.Mkdir(npmPath+backupSuffix, 0o700); err != nil {
		t.Fatal(err)
	}
	service := newAuthService(authenticatedRunner(), func() (string, error) { return home, nil }, &bytes.Buffer{})

	err := service.Auth(context.Background(), AuthRequest{Action: AuthLogin})
	if ExitCodeFor(err) != ExitConfiguration {
		t.Fatalf("backup failure error = %v, code = %d", err, ExitCodeFor(err))
	}
	if got := mustRead(t, settingsPath); got != settingsBefore {
		t.Fatalf("settings were partially updated: %q", got)
	}
	if got := mustRead(t, npmPath); got != npmBefore {
		t.Fatalf("npmrc was partially updated: %q", got)
	}
	if got := mustRead(t, settingsPath+backupSuffix); got != "previous settings backup" {
		t.Fatalf("previous backup was not restored: %q", got)
	}
}

func TestAuthStatusRejectsIncompleteManagedCredentials(t *testing.T) {
	home := t.TempDir()
	mustWrite(t, filepath.Join(home, ".m2", "settings.xml"), "<settings><servers><server><id>github-boilerplate</id><username>dev</username></server></servers></settings>\n")
	mustWrite(t, filepath.Join(home, ".npmrc"), npmManagedBegin+"\n"+npmScopeLine+"\n"+npmManagedEnd+"\n")
	service := newAuthService(authenticatedRunner(), func() (string, error) { return home, nil }, &bytes.Buffer{})

	err := service.Auth(context.Background(), AuthRequest{Action: AuthStatus})
	if ExitCodeFor(err) != ExitConfiguration {
		t.Fatalf("incomplete status error = %v, code = %d", err, ExitCodeFor(err))
	}
}

func TestAuthRejectsAmbiguousMavenServersWithoutMutation(t *testing.T) {
	home := t.TempDir()
	settingsPath := filepath.Join(home, ".m2", "settings.xml")
	ambiguous := "<settings><servers><server><id>other</id></server></servers><servers><server><id>github-boilerplate</id><username>old</username><password>old</password></server></servers></settings>\n"
	mustWrite(t, settingsPath, ambiguous)
	service := newAuthService(authenticatedRunner(), func() (string, error) { return home, nil }, &bytes.Buffer{})

	err := service.Auth(context.Background(), AuthRequest{Action: AuthLogin})
	if ExitCodeFor(err) != ExitConflict {
		t.Fatalf("ambiguous Maven error = %v, code = %d", err, ExitCodeFor(err))
	}
	if got := mustRead(t, settingsPath); got != ambiguous {
		t.Fatalf("ambiguous Maven settings were changed: %q", got)
	}
}

func TestAuthIgnoresMavenElementsInsideComments(t *testing.T) {
	home := t.TempDir()
	settingsPath := filepath.Join(home, ".m2", "settings.xml")
	original := "<settings>\n  <!-- example: <servers><server><id>github-boilerplate</id></server></servers> -->\n  <servers><server><id>other</id></server></servers>\n</settings>\n"
	mustWrite(t, settingsPath, original)
	service := newAuthService(authenticatedRunner(), func() (string, error) { return home, nil }, &bytes.Buffer{})

	if err := service.Auth(context.Background(), AuthRequest{Action: AuthLogin}); err != nil {
		t.Fatal(err)
	}
	settings := mustRead(t, settingsPath)
	if !strings.Contains(settings, "<!-- example: <servers><server><id>github-boilerplate</id></server></servers> -->") {
		t.Fatalf("commented Maven example was changed: %s", settings)
	}
	if strings.Count(settings, "<password>") != 1 || !containsMavenServer([]byte(settings)) {
		t.Fatalf("active Maven server was not configured: %s", settings)
	}
}

func TestDefaultDependenciesWireAuthWithoutPassingTokenAsArgument(t *testing.T) {
	home := t.TempDir()
	runner := authenticatedRunner()
	var stdout, stderr bytes.Buffer
	code := RunCLI(context.Background(), []string{"auth", "login"}, Dependencies{
		Runner:  runner,
		HomeDir: func() (string, error) { return home, nil },
		Stdout:  &stdout,
		Stderr:  &stderr,
	})
	if code != int(ExitSuccess) {
		t.Fatalf("exit=%d stderr=%q", code, stderr.String())
	}
	for _, call := range runner.calls {
		if strings.Contains(call, testToken) {
			t.Fatalf("token was passed as process argument: %s", call)
		}
	}
	if strings.Contains(stdout.String()+stderr.String(), testToken) {
		t.Fatal("token leaked through CLI streams")
	}
}

func authenticatedRunner() *fakeProcessRunner {
	return &fakeProcessRunner{
		responses: map[string]string{
			"gh auth token --hostname github.com":           testToken + "\n",
			"gh api --hostname github.com user --jq .login": "dev&ops\n",
			"gh auth status --hostname github.com":          "github.com\n",
		},
		errors: map[string]error{},
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
