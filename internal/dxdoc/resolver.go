package dxdoc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/install"
)

type Request struct {
	Language install.Language
	Package  string
	WorkDir  string
	Local    bool
}

type Resolution struct {
	Language     install.Language
	Package      string
	ManifestHint string
	Source       string
}

func Resolve(request Request) (*Resolution, error) {
	if request.Package == "" {
		return nil, fmt.Errorf("package e obrigatorio")
	}
	if request.WorkDir == "" {
		workDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("nao foi possivel detectar diretorio atual: %w", err)
		}
		request.WorkDir = workDir
	}

	switch request.Language {
	case install.LanguageVue:
		source := "docs.remote"
		if request.Local {
			source = "docs.local"
		}
		return &Resolution{
			Language:     request.Language,
			Package:      request.Package,
			ManifestHint: filepath.Join(request.WorkDir, "node_modules", request.Package, ".boilerplate", "package.manifest.json"),
			Source:       source,
		}, nil
	case install.LanguageJava:
		source := "docs.remote"
		if request.Local {
			source = "docs.javadoc"
		}
		return &Resolution{
			Language:     request.Language,
			Package:      request.Package,
			ManifestHint: filepath.Join("META-INF", "boilerplate", "package.manifest.json"),
			Source:       source,
		}, nil
	default:
		return nil, fmt.Errorf("linguagem nao suportada: %s", request.Language)
	}
}
