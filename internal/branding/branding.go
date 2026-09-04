package branding

import (
	"encoding/json"
	"os"
	"strings"
)

const (
	DefaultName      = "boilerplate"
	DefaultNpmScope  = "@aplicacoesBoilerplate"
	DefaultMavenRepo = "github-boilerplate"
)

type Config struct {
	Name       string `json:"name"`
	Binary     string `json:"binary"`
	NpmScope   string `json:"npmScope"`
	MavenRepo  string `json:"mavenRepo"`
	GitHubOrg  string `json:"githubOrg"`
	MavenGroup string `json:"mavenGroupId"`
}

func Default() Config {
	return Config{
		Name:       DefaultName,
		Binary:     DefaultName,
		NpmScope:   DefaultNpmScope,
		MavenRepo:  DefaultMavenRepo,
		GitHubOrg:  "aplicacoesBoilerplate",
		MavenGroup: "br.com.aplicacoesBoilerplate",
	}
}

func Load() Config {
	config := Default()
	loadFromFile(&config, ".boilerplate-cli.json")
	applyEnv(&config)
	normalize(&config)
	return config
}

func loadFromFile(config *Config, path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	_ = json.NewDecoder(file).Decode(config)
}

func applyEnv(config *Config) {
	if value := os.Getenv("BOILERPLATE_CLI_NAME"); value != "" {
		config.Name = value
	}
	if value := os.Getenv("BOILERPLATE_CLI_BINARY"); value != "" {
		config.Binary = value
	}
	if value := os.Getenv("BOILERPLATE_NPM_SCOPE"); value != "" {
		config.NpmScope = value
	}
	if value := os.Getenv("BOILERPLATE_MAVEN_REPO"); value != "" {
		config.MavenRepo = value
	}
	if value := os.Getenv("BOILERPLATE_GITHUB_ORG"); value != "" {
		config.GitHubOrg = value
	}
	if value := os.Getenv("BOILERPLATE_MAVEN_GROUP_ID"); value != "" {
		config.MavenGroup = value
	}
}

func normalize(config *Config) {
	if strings.TrimSpace(config.Name) == "" {
		config.Name = DefaultName
	}
	if strings.TrimSpace(config.Binary) == "" {
		config.Binary = config.Name
	}
	if strings.TrimSpace(config.NpmScope) == "" {
		config.NpmScope = DefaultNpmScope
	}
	if strings.TrimSpace(config.MavenRepo) == "" {
		config.MavenRepo = DefaultMavenRepo
	}
}
