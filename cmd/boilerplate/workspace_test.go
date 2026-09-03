package boilerplate

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDiscoverWorkspaceFindsOrderedMixedMonorepoAndHonorsBoundaries(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "pom.xml"), "<project/>")
	mustWrite(t, filepath.Join(root, "apps", "web", "package.json"), `{ "name": "web" }`)
	mustWrite(t, filepath.Join(root, "services", "api", "pom.xml"), "<project/>")
	mustWrite(t, filepath.Join(root, "node_modules", "ignored", "package.json"), `{}`)
	mustWrite(t, filepath.Join(root, "target", "ignored", "pom.xml"), "<project/>")
	mustWrite(t, filepath.Join(root, "nested-repo", ".git", "HEAD"), "ref: refs/heads/main")
	mustWrite(t, filepath.Join(root, "nested-repo", "package.json"), `{}`)

	workspace, err := DiscoverWorkspace(root)
	if err != nil {
		t.Fatal(err)
	}

	wantJava := []string{root, filepath.Join(root, "services", "api")}
	wantVue := []string{filepath.Join(root, "apps", "web")}
	if !reflect.DeepEqual(workspace.JavaProjects, wantJava) {
		t.Fatalf("java projects = %#v, want %#v", workspace.JavaProjects, wantJava)
	}
	if !reflect.DeepEqual(workspace.VueProjects, wantVue) {
		t.Fatalf("vue projects = %#v, want %#v", workspace.VueProjects, wantVue)
	}
}

func TestDiscoverWorkspaceRejectsMissingOrNonDirectoryRoot(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{filepath.Join(root, "missing"), file} {
		if _, err := DiscoverWorkspace(path); ExitCodeFor(err) != ExitConfiguration {
			t.Fatalf("DiscoverWorkspace(%q) error = %v, code = %d", path, err, ExitCodeFor(err))
		}
	}
}
