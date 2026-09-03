package boilerplate

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

var ignoredWorkspaceDirectories = map[string]struct{}{
	".git":         {},
	".idea":        {},
	"build":        {},
	"dist":         {},
	"node_modules": {},
	"target":       {},
}

func DiscoverWorkspace(root string) (Workspace, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return Workspace{}, NewCLIError(ExitConfiguration, "nao foi possivel resolver o diretorio raiz", err)
	}
	absolute = filepath.Clean(absolute)
	info, err := os.Stat(absolute)
	if err != nil || !info.IsDir() {
		return Workspace{}, NewCLIError(ExitConfiguration, "o diretorio raiz nao existe ou nao e um diretorio", err)
	}

	javaProjects := map[string]struct{}{}
	vueProjects := map[string]struct{}{}
	err = filepath.WalkDir(absolute, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path != absolute {
				if _, ignored := ignoredWorkspaceDirectories[entry.Name()]; ignored {
					return filepath.SkipDir
				}
				if entry.Type()&os.ModeSymlink != 0 || nestedGitBoundary(path) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		switch entry.Name() {
		case "pom.xml":
			javaProjects[filepath.Dir(path)] = struct{}{}
		case "package.json":
			vueProjects[filepath.Dir(path)] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return Workspace{}, NewCLIError(ExitConfiguration, "nao foi possivel descobrir os projetos do workspace", err)
	}

	return Workspace{
		Root:         absolute,
		JavaProjects: orderedKeys(javaProjects),
		VueProjects:  orderedKeys(vueProjects),
	}, nil
}

func nestedGitBoundary(path string) bool {
	_, err := os.Lstat(filepath.Join(path, ".git"))
	return err == nil
}

func orderedKeys(values map[string]struct{}) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
