package dxartifacts

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/manifest"
)

type Request struct {
	Manifest    *manifest.Manifest
	Preset      string
	PackageRoot string
	ProjectRoot string
}

type ArtifactKind string

const (
	ArtifactTemplate      ArtifactKind = "template"
	ArtifactVSCodeSnippet ArtifactKind = "vscode-snippet"
	ArtifactWebTypes      ArtifactKind = "web-types"
	ArtifactLiveTemplate  ArtifactKind = "jetbrains-live-template"
)

type Artifact struct {
	Kind        ArtifactKind
	Source      string
	Destination string
}

type Plan struct {
	Package   string
	Language  string
	Preset    string
	Artifacts []Artifact
}

func BuildPlan(request Request) (*Plan, error) {
	if request.Manifest == nil {
		return nil, fmt.Errorf("manifesto DX e obrigatorio")
	}
	if request.Preset == "" {
		request.Preset = "default"
	}
	if request.PackageRoot == "" {
		request.PackageRoot = "."
	}
	if request.ProjectRoot == "" {
		request.ProjectRoot = "."
	}

	plan := &Plan{
		Package:  request.Manifest.Package,
		Language: request.Manifest.Language,
		Preset:   request.Preset,
	}

	if preset, ok := request.Manifest.Presets[request.Preset]; ok {
		for _, source := range preset.Templates {
			plan.Artifacts = append(plan.Artifacts, Artifact{
				Kind:        ArtifactTemplate,
				Source:      filepath.Join(request.PackageRoot, source),
				Destination: templateDestination(request.ProjectRoot, source),
			})
		}
	}

	for _, source := range request.Manifest.IDE.VSCodeSnippets() {
		plan.Artifacts = append(plan.Artifacts, Artifact{
			Kind:        ArtifactVSCodeSnippet,
			Source:      filepath.Join(request.PackageRoot, source),
			Destination: filepath.Join(request.ProjectRoot, ".vscode", "boilerplate.code-snippets"),
		})
	}

	if request.Manifest.IDE.JetBrains != nil {
		safeName := safeArtifactName(request.Manifest.Package)
		if request.Manifest.IDE.JetBrains.WebTypes != "" {
			plan.Artifacts = append(plan.Artifacts, Artifact{
				Kind:        ArtifactWebTypes,
				Source:      filepath.Join(request.PackageRoot, request.Manifest.IDE.JetBrains.WebTypes),
				Destination: filepath.Join(request.ProjectRoot, ".idea", "webTypes", safeName+".json"),
			})
		}
		for _, source := range request.Manifest.IDE.JetBrains.LiveTemplates {
			plan.Artifacts = append(plan.Artifacts, Artifact{
				Kind:        ArtifactLiveTemplate,
				Source:      filepath.Join(request.PackageRoot, source),
				Destination: filepath.Join(request.ProjectRoot, ".idea", "templates", safeName+".xml"),
			})
		}
	}

	return plan, nil
}

func templateDestination(projectRoot string, source string) string {
	base := filepath.Base(source)
	switch base {
	case "application.yml", "application.yaml", "application.properties":
		return filepath.Join(projectRoot, "src", "main", "resources", base)
	case "plugin.ts":
		return filepath.Join(projectRoot, "src", "plugins", base)
	case "index.scss":
		return filepath.Join(projectRoot, "src", "styles", base)
	default:
		return filepath.Join(projectRoot, base)
	}
}

func safeArtifactName(name string) string {
	replacer := strings.NewReplacer(
		"@", "",
		"/", "-",
		"\\", "-",
		":", "-",
		" ", "-",
	)
	safe := replacer.Replace(strings.TrimSpace(name))
	if safe == "" {
		return "package"
	}
	return safe
}