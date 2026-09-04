package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

const SchemaVersion = "1.0"

type Manifest struct {
	SchemaVersion     string            `json:"schemaVersion"`
	Package           string            `json:"package"`
	Language          string            `json:"language"`
	Aliases           []string          `json:"aliases,omitempty"`
	MinimumCLIVersion string            `json:"minimumCliVersion"`
	GroupID           string            `json:"groupId,omitempty"`
	ArtifactID        string            `json:"artifactId,omitempty"`
	Docs              Docs              `json:"docs,omitempty"`
	IDE               IDE               `json:"ide,omitempty"`
	Presets           map[string]Preset `json:"presets,omitempty"`
}

type Docs struct {
	Remote  string `json:"remote,omitempty"`
	Local   string `json:"local,omitempty"`
	Javadoc string `json:"javadoc,omitempty"`
}

type IDE struct {
	VSCode    *VSCodeIDE    `json:"vscode,omitempty"`
	JetBrains *JetBrainsIDE `json:"jetbrains,omitempty"`
}

type VSCodeIDE struct {
	Snippets []string `json:"snippets,omitempty"`
}

type JetBrainsIDE struct {
	WebTypes      string   `json:"webTypes,omitempty"`
	LiveTemplates []string `json:"liveTemplates,omitempty"`
}

type Preset struct {
	Description string   `json:"description,omitempty"`
	Templates   []string `json:"templates,omitempty"`
}

func Decode(r io.Reader) (*Manifest, error) {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return nil, fmt.Errorf("manifesto DX invalido: %w", err)
	}

	if err := manifest.Validate(); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (m Manifest) Validate() error {
	if strings.TrimSpace(m.SchemaVersion) == "" {
		return fmt.Errorf("manifesto DX invalido: schemaVersion e obrigatorio")
	}
	if m.SchemaVersion != SchemaVersion {
		return fmt.Errorf("manifesto DX invalido: schemaVersion %q nao suportado; use %q", m.SchemaVersion, SchemaVersion)
	}
	if strings.TrimSpace(m.Package) == "" {
		return fmt.Errorf("manifesto DX invalido: package e obrigatorio")
	}
	if strings.TrimSpace(m.MinimumCLIVersion) == "" {
		return fmt.Errorf("manifesto DX invalido: minimumCliVersion e obrigatorio")
	}

	switch m.Language {
	case "java", "vue":
	default:
		return fmt.Errorf("manifesto DX invalido: language deve ser java ou vue")
	}

	for name, preset := range m.Presets {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("manifesto DX invalido: preset com nome vazio")
		}
		for _, path := range preset.Templates {
			if err := validatePackagePath(path); err != nil {
				return fmt.Errorf("manifesto DX invalido: preset %q referencia caminho inseguro %q: %w", name, path, err)
			}
		}
	}

	for _, path := range m.IDE.VSCodeSnippets() {
		if err := validatePackagePath(path); err != nil {
			return fmt.Errorf("manifesto DX invalido: snippet VS Code referencia caminho inseguro %q: %w", path, err)
		}
	}
	for _, path := range m.IDE.JetBrainsArtifacts() {
		if err := validatePackagePath(path); err != nil {
			return fmt.Errorf("manifesto DX invalido: artefato JetBrains referencia caminho inseguro %q: %w", path, err)
		}
	}

	return nil
}

func (m Manifest) PresetNames() []string {
	names := make([]string, 0, len(m.Presets))
	for name := range m.Presets {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (m Manifest) IDENames() []string {
	var names []string
	if m.IDE.VSCode != nil {
		names = append(names, "vscode")
	}
	if m.IDE.JetBrains != nil {
		names = append(names, "jetbrains")
	}
	return names
}

func (i IDE) VSCodeSnippets() []string {
	if i.VSCode == nil {
		return nil
	}
	return i.VSCode.Snippets
}

func (i IDE) JetBrainsArtifacts() []string {
	if i.JetBrains == nil {
		return nil
	}

	artifacts := append([]string{}, i.JetBrains.LiveTemplates...)
	if i.JetBrains.WebTypes != "" {
		artifacts = append(artifacts, i.JetBrains.WebTypes)
	}
	return artifacts
}

func validatePackagePath(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("caminho vazio")
	}
	if filepath.IsAbs(path) {
		return fmt.Errorf("caminho absoluto nao permitido")
	}

	clean := filepath.Clean(path)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return fmt.Errorf("caminho fora da raiz do package")
	}
	if strings.Contains(clean, "\x00") {
		return fmt.Errorf("caminho contem byte nulo")
	}

	return nil
}
