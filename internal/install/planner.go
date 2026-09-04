package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Language string

const (
	LanguageJava Language = "java"
	LanguageVue  Language = "vue"
)

type Request struct {
	Language Language
	Package  string
	Preset   string
	WorkDir  string
	DryRun   bool
}

type Plan struct {
	Language     Language
	Package      string
	Preset       string
	WorkDir      string
	ProjectFiles []string
	ManifestHint string
	Actions      []string
}

func BuildPlan(request Request) (*Plan, error) {
	if strings.TrimSpace(request.Package) == "" {
		return nil, fmt.Errorf("package e obrigatorio")
	}
	if request.WorkDir == "" {
		workDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("nao foi possivel detectar diretorio atual: %w", err)
		}
		request.WorkDir = workDir
	}
	if request.Preset == "" {
		request.Preset = "default"
	}

	switch request.Language {
	case LanguageVue:
		return buildVuePlan(request)
	case LanguageJava:
		return buildJavaPlan(request)
	default:
		return nil, fmt.Errorf("linguagem nao suportada: %s", request.Language)
	}
}

func buildVuePlan(request Request) (*Plan, error) {
	packageJSON := filepath.Join(request.WorkDir, "package.json")
	if _, err := os.Stat(packageJSON); err != nil {
		return nil, fmt.Errorf("projeto Vue/npm invalido: package.json nao encontrado em %s", request.WorkDir)
	}

	manifestHint := filepath.Join(request.WorkDir, "node_modules", request.Package, ".boilerplate", "package.manifest.json")
	return &Plan{
		Language:     LanguageVue,
		Package:      request.Package,
		Preset:       request.Preset,
		WorkDir:      request.WorkDir,
		ProjectFiles: []string{"package.json"},
		ManifestHint: manifestHint,
		Actions: []string{
			fmt.Sprintf("instalar package npm %s", request.Package),
			fmt.Sprintf("ler manifesto DX em %s", manifestHint),
			fmt.Sprintf("aplicar preset %s com merge seguro", request.Preset),
			"sincronizar snippets e metadados de IDE declarados pelo manifesto",
		},
	}, nil
}

func buildJavaPlan(request Request) (*Plan, error) {
	pomXML := filepath.Join(request.WorkDir, "pom.xml")
	if _, err := os.Stat(pomXML); err != nil {
		return nil, fmt.Errorf("projeto Java/Maven invalido: pom.xml nao encontrado em %s", request.WorkDir)
	}

	manifestHint := filepath.Join("META-INF", "boilerplate", "package.manifest.json")
	return &Plan{
		Language:     LanguageJava,
		Package:      request.Package,
		Preset:       request.Preset,
		WorkDir:      request.WorkDir,
		ProjectFiles: []string{"pom.xml"},
		ManifestHint: manifestHint,
		Actions: []string{
			fmt.Sprintf("adicionar dependency Maven %s", request.Package),
			"resolver dependencias Maven",
			fmt.Sprintf("extrair manifesto DX do jar em %s", manifestHint),
			fmt.Sprintf("aplicar preset %s em application.yml/properties com merge seguro", request.Preset),
			"sincronizar snippets e live templates declarados pelo manifesto",
		},
	}, nil
}
