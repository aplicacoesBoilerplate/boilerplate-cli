package boilerplate

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var projectNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

func newAuthCommand(state *commandContext) *cobra.Command {
	command := &cobra.Command{Use: "auth", Short: "Gerencia autenticacao do GitHub Packages"}
	for _, action := range []AuthAction{AuthLogin, AuthLogout, AuthStatus} {
		action := action
		command.AddCommand(&cobra.Command{
			Use:   string(action),
			Short: authDescription(action),
			Args:  noArgs,
			RunE: func(command *cobra.Command, _ []string) error {
				root, err := state.absoluteRoot()
				if err != nil {
					return err
				}
				return state.service.Auth(command.Context(), AuthRequest{Root: root, DryRun: state.dryRun, Action: action})
			},
		})
	}
	return command
}

func authDescription(action AuthAction) string {
	switch action {
	case AuthLogin:
		return "Configura credenciais reutilizando gh auth"
	case AuthLogout:
		return "Remove somente configuracoes gerenciadas pela CLI"
	default:
		return "Valida ferramentas e acesso ao GitHub Packages"
	}
}

func newInitCommand(state *commandContext) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Inicializa um projeto ou monorepo existente",
		Args:  noArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			workspace, err := state.workspace()
			if err != nil {
				return err
			}
			return state.service.Init(command.Context(), InitRequest{Workspace: workspace, DryRun: state.dryRun})
		},
	}
}

func newNewCommand(state *commandContext) *cobra.Command {
	var owner, directory, visibility string
	command := &cobra.Command{
		Use:   "new <java|vue> <nome>",
		Short: "Cria uma aplicacao a partir do template canonico",
		Args:  exactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			platform, err := parsePlatform(args[0], false)
			if err != nil {
				return err
			}
			if !projectNamePattern.MatchString(args[1]) {
				return usageError("nome de projeto invalido")
			}
			parsedVisibility, err := parseVisibility(visibility)
			if err != nil {
				return err
			}
			root, err := state.absoluteRoot()
			if err != nil {
				return err
			}
			if directory == "" {
				directory = args[1]
			}
			if !filepath.IsAbs(directory) {
				directory = filepath.Join(root, directory)
			}
			return state.service.New(command.Context(), NewRequest{
				Root: root, DryRun: state.dryRun, Platform: platform, Name: args[1],
				Owner: owner, Directory: filepath.Clean(directory), Visibility: parsedVisibility,
			})
		},
	}
	command.Flags().StringVar(&owner, "owner", "", "organizacao ou usuario de destino")
	command.Flags().StringVar(&directory, "directory", "", "diretorio local de destino")
	command.Flags().StringVar(&visibility, "visibility", string(VisibilityPrivate), "visibilidade: private, internal ou public")
	return command
}

func newAddCommand(state *commandContext) *cobra.Command {
	return &cobra.Command{
		Use:   "add <java|vue> <pacote>@<versao>",
		Short: "Adiciona um package com versao explicita",
		Args:  exactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			platform, err := parsePlatform(args[0], false)
			if err != nil {
				return err
			}
			name, version, err := parsePackageSpec(args[1])
			if err != nil {
				return err
			}
			workspace, err := state.workspace()
			if err != nil {
				return err
			}
			return state.service.Add(command.Context(), AddRequest{
				Workspace: workspace, DryRun: state.dryRun, Platform: platform, Package: name, Version: version,
			})
		},
	}
}

func newUpdateCommand(state *commandContext) *cobra.Command {
	return &cobra.Command{
		Use:   "update [java|vue|all]",
		Short: "Atualiza somente packages do ecossistema",
		Args:  maximumArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			value := string(PlatformAll)
			if len(args) == 1 {
				value = args[0]
			}
			platform, err := parsePlatform(value, true)
			if err != nil {
				return err
			}
			workspace, err := state.workspace()
			if err != nil {
				return err
			}
			return state.service.Update(command.Context(), UpdateRequest{Workspace: workspace, DryRun: state.dryRun, Platform: platform})
		},
	}
}

func newDoctorCommand(state *commandContext) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnostica ferramentas, credenciais e projetos",
		Args:  noArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			workspace, err := state.workspace()
			if err != nil {
				return err
			}
			return state.service.Doctor(command.Context(), DoctorRequest{Workspace: workspace})
		},
	}
}

func newAuditCommand(state *commandContext) *cobra.Command {
	var format, output string
	command := &cobra.Command{
		Use:   "audit",
		Short: "Relata defasagem dos packages sem modificar o workspace",
		Args:  noArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			parsedFormat, err := parseAuditFormat(format)
			if err != nil {
				return err
			}
			workspace, err := state.workspace()
			if err != nil {
				return err
			}
			if output != "" && !filepath.IsAbs(output) {
				output = filepath.Join(workspace.Root, output)
			}
			if output != "" {
				output = filepath.Clean(output)
			}
			return state.service.Audit(command.Context(), AuditRequest{
				Workspace: workspace, DryRun: state.dryRun, Format: parsedFormat, Output: output,
			})
		},
	}
	command.Flags().StringVar(&format, "format", string(AuditFormatText), "formato de saida: text ou json")
	command.Flags().StringVar(&output, "output", "", "arquivo de saida; vazio escreve no stdout")
	return command
}

func parsePlatform(value string, allowAll bool) (Platform, error) {
	platform := Platform(strings.ToLower(strings.TrimSpace(value)))
	if platform == PlatformJava || platform == PlatformVue || (allowAll && platform == PlatformAll) {
		return platform, nil
	}
	return "", usageError("plataforma invalida: use java, vue%s", map[bool]string{true: " ou all", false: ""}[allowAll])
}

func parseVisibility(value string) (Visibility, error) {
	visibility := Visibility(strings.ToLower(strings.TrimSpace(value)))
	if visibility == VisibilityPrivate || visibility == VisibilityInternal || visibility == VisibilityPublic {
		return visibility, nil
	}
	return "", usageError("visibilidade invalida: use private, internal ou public")
}

func parseAuditFormat(value string) (AuditFormat, error) {
	format := AuditFormat(strings.ToLower(strings.TrimSpace(value)))
	if format == AuditFormatText || format == AuditFormatJSON {
		return format, nil
	}
	return "", usageError("formato invalido: use text ou json")
}

func parsePackageSpec(value string) (string, string, error) {
	separator := strings.LastIndex(value, "@")
	if separator <= 0 || separator == len(value)-1 {
		return "", "", usageError("pacote deve usar o formato <pacote>@<versao>")
	}
	name := strings.TrimSpace(value[:separator])
	version := strings.TrimSpace(value[separator+1:])
	if name == "" || version == "" || strings.ContainsAny(name+version, "\r\n\t ") {
		return "", "", usageError("pacote e versao nao podem conter espacos")
	}
	return name, version, nil
}
