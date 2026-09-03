package boilerplate

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	backupSuffix    = ".boilerplate-cli.bak"
	npmManagedBegin = "# boilerplate-cli:begin"
	npmManagedEnd   = "# boilerplate-cli:end"
	npmScopeLine    = "@aplicacoesBoilerplate:registry=https://npm.pkg.github.com"
	npmTokenPrefix  = "//npm.pkg.github.com/:_authToken="
	ghScopesQuery   = `.hosts["github.com"][0].scopes`
)

type ProcessRunner interface {
	Run(context.Context, string, ...string) (string, error)
}

type execProcessRunner struct{}

func (execProcessRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, name, args...)
	var stdout bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = io.Discard
	if err := command.Run(); err != nil {
		return "", err
	}
	return stdout.String(), nil
}

type authService struct {
	runner  ProcessRunner
	homeDir func() (string, error)
	output  io.Writer
}

func newAuthService(runner ProcessRunner, homeDir func() (string, error), output io.Writer) *authService {
	return &authService{runner: runner, homeDir: homeDir, output: output}
}

func (s *authService) Auth(ctx context.Context, request AuthRequest) error {
	switch request.Action {
	case AuthLogin:
		return s.login(ctx, request.DryRun)
	case AuthLogout:
		return s.logout(request.DryRun)
	case AuthStatus:
		return s.status(ctx)
	default:
		return usageError("acao de autenticacao invalida")
	}
}

func (s *authService) login(ctx context.Context, dryRun bool) error {
	if err := s.validatePackagesScope(ctx); err != nil {
		return err
	}
	tokenOutput, err := s.runner.Run(ctx, "gh", "auth", "token", "--hostname", "github.com")
	if err != nil {
		return NewCLIError(ExitAuthentication, "GitHub CLI nao esta autenticado em github.com", err)
	}
	token := strings.TrimSpace(tokenOutput)
	if token == "" || strings.ContainsAny(token, "\r\n") {
		return NewCLIError(ExitAuthentication, "GitHub CLI nao retornou uma credencial valida", nil)
	}

	userOutput, err := s.runner.Run(ctx, "gh", "api", "--hostname", "github.com", "user", "--jq", ".login")
	if err != nil {
		return NewCLIError(ExitAuthentication, "nao foi possivel identificar a conta ativa do GitHub", err)
	}
	username := strings.TrimSpace(userOutput)
	if username == "" || strings.ContainsAny(username, "\r\n") {
		return NewCLIError(ExitAuthentication, "GitHub CLI nao retornou uma conta valida", nil)
	}

	home, err := s.resolveHome()
	if err != nil {
		return err
	}
	settingsPath := filepath.Join(home, ".m2", "settings.xml")
	npmPath := filepath.Join(home, ".npmrc")
	settingsBefore, settingsExists, err := readOptional(settingsPath)
	if err != nil {
		return NewCLIError(ExitConfiguration, "nao foi possivel ler o settings.xml", err)
	}
	npmBefore, npmExists, err := readOptional(npmPath)
	if err != nil {
		return NewCLIError(ExitConfiguration, "nao foi possivel ler o .npmrc", err)
	}

	settingsAfter, err := upsertMavenServer(settingsBefore, username, token)
	if err != nil {
		return err
	}
	npmAfter, err := upsertNPMCredentials(npmBefore, token)
	if err != nil {
		return err
	}
	changes := []fileChange{
		{path: settingsPath, before: settingsBefore, after: settingsAfter, existed: settingsExists},
		{path: npmPath, before: npmBefore, after: npmAfter, existed: npmExists},
	}
	if dryRun {
		_, _ = fmt.Fprintln(s.output, "dry-run: configuracoes Maven e npm seriam atualizadas")
		return nil
	}
	if err := applyFileChanges(changes); err != nil {
		return NewCLIError(ExitConfiguration, "nao foi possivel gravar as configuracoes de packages", err)
	}
	_, _ = fmt.Fprintln(s.output, "autenticacao configurada para Maven e npm")
	return nil
}

func (s *authService) logout(dryRun bool) error {
	home, err := s.resolveHome()
	if err != nil {
		return err
	}
	settingsPath := filepath.Join(home, ".m2", "settings.xml")
	npmPath := filepath.Join(home, ".npmrc")
	settingsBefore, settingsExists, err := readOptional(settingsPath)
	if err != nil {
		return NewCLIError(ExitConfiguration, "nao foi possivel ler o settings.xml", err)
	}
	npmBefore, npmExists, err := readOptional(npmPath)
	if err != nil {
		return NewCLIError(ExitConfiguration, "nao foi possivel ler o .npmrc", err)
	}
	settingsAfter, err := removeMavenServer(settingsBefore)
	if err != nil {
		return err
	}
	npmAfter, err := removeNPMManagedBlock(npmBefore)
	if err != nil {
		return err
	}
	changes := []fileChange{
		{path: settingsPath, before: settingsBefore, after: settingsAfter, existed: settingsExists},
		{path: npmPath, before: npmBefore, after: npmAfter, existed: npmExists},
	}
	if dryRun {
		_, _ = fmt.Fprintln(s.output, "dry-run: configuracoes gerenciadas seriam removidas")
		return nil
	}
	if err := applyFileChanges(changes); err != nil {
		return NewCLIError(ExitConfiguration, "nao foi possivel remover as configuracoes gerenciadas", err)
	}
	_, _ = fmt.Fprintln(s.output, "configuracoes gerenciadas removidas")
	return nil
}

func (s *authService) status(ctx context.Context) error {
	if err := s.validatePackagesScope(ctx); err != nil {
		return err
	}
	home, err := s.resolveHome()
	if err != nil {
		return err
	}
	settings, settingsExists, err := readOptional(filepath.Join(home, ".m2", "settings.xml"))
	if err != nil {
		return NewCLIError(ExitConfiguration, "nao foi possivel ler o settings.xml", err)
	}
	npm, npmExists, err := readOptional(filepath.Join(home, ".npmrc"))
	if err != nil {
		return NewCLIError(ExitConfiguration, "nao foi possivel ler o .npmrc", err)
	}
	if !settingsExists || !containsMavenServer(settings) || !npmExists || !containsNPMManagedBlock(npm) {
		return NewCLIError(ExitConfiguration, "as configuracoes Maven e npm ainda nao estao completas", nil)
	}
	_, _ = fmt.Fprintln(s.output, "GitHub CLI, Maven e npm configurados")
	return nil
}

func (s *authService) validatePackagesScope(ctx context.Context) error {
	scopesOutput, err := s.runner.Run(ctx, "gh", "auth", "status", "--hostname", "github.com",
		"--active", "--json", "hosts", "--jq", ghScopesQuery)
	if err != nil {
		return NewCLIError(ExitAuthentication, "GitHub CLI nao esta autenticado em github.com", err)
	}
	for _, scope := range strings.Split(strings.TrimSpace(scopesOutput), ",") {
		normalized := strings.TrimSpace(scope)
		if normalized == "read:packages" || normalized == "write:packages" {
			return nil
		}
	}
	return NewCLIError(ExitAuthentication,
		"a credencial ativa do GitHub precisa do escopo read:packages", nil)
}

func (s *authService) resolveHome() (string, error) {
	home, err := s.homeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return "", NewCLIError(ExitConfiguration, "nao foi possivel localizar o diretorio do usuario", err)
	}
	absolute, err := filepath.Abs(home)
	if err != nil {
		return "", NewCLIError(ExitConfiguration, "nao foi possivel resolver o diretorio do usuario", err)
	}
	return filepath.Clean(absolute), nil
}

type fileChange struct {
	path    string
	before  []byte
	after   []byte
	existed bool
}

func applyFileChanges(changes []fileChange) error {
	pending := make([]fileChange, 0, len(changes))
	for _, change := range changes {
		if !bytes.Equal(change.before, change.after) {
			pending = append(pending, change)
		}
	}

	var backups []fileChange
	for _, change := range pending {
		if change.existed {
			backupPath := change.path + backupSuffix
			backupBefore, backupExists, err := readOptional(backupPath)
			if err != nil {
				return errors.Join(err, rollbackChanges(backups))
			}
			backup := fileChange{path: backupPath, before: backupBefore, after: change.before, existed: backupExists}
			if bytes.Equal(backup.before, backup.after) {
				continue
			}
			if err := writeAtomic(backup.path, backup.after, 0o600); err != nil {
				return errors.Join(err, rollbackChanges(backups))
			}
			backups = append(backups, backup)
		}
	}

	var changed []fileChange
	for _, change := range pending {
		if err := writeAtomic(change.path, change.after, 0o600); err != nil {
			return errors.Join(err, rollbackChanges(changed), rollbackChanges(backups))
		}
		changed = append(changed, change)
	}
	return nil
}

func rollbackChanges(changes []fileChange) error {
	var rollbackErrors []error
	for index := len(changes) - 1; index >= 0; index-- {
		change := changes[index]
		if change.existed {
			if err := writeAtomic(change.path, change.before, 0o600); err != nil {
				rollbackErrors = append(rollbackErrors, err)
			}
		} else if err := os.Remove(change.path); err != nil && !errors.Is(err, os.ErrNotExist) {
			rollbackErrors = append(rollbackErrors, err)
		}
	}
	return errors.Join(rollbackErrors...)
}

func writeAtomic(path string, content []byte, permission os.FileMode) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(directory, ".boilerplate-cli-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()
	if err := temporary.Chmod(permission); err != nil {
		_ = temporary.Close()
		return err
	}
	if _, err := temporary.Write(content); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, path)
}

func readOptional(path string) ([]byte, bool, error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	return content, err == nil, err
}

func upsertMavenServer(content []byte, username, token string) ([]byte, error) {
	if len(bytes.TrimSpace(content)) == 0 {
		content = []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<settings>\n</settings>\n")
	}
	if err := validateSettingsXML(content); err != nil {
		return nil, NewCLIError(ExitConfiguration, "settings.xml nao contem XML valido", err)
	}
	selfClosing, err := rootElementSelfClosing(content, "settings")
	if err != nil {
		return nil, NewCLIError(ExitConfiguration, "nao foi possivel interpretar o settings.xml", err)
	}
	if selfClosing {
		return nil, NewCLIError(ExitConflict, "settings.xml usa o elemento settings autocontido", nil)
	}
	sections, err := directChildRanges(content, "settings", "servers")
	if err != nil {
		return nil, NewCLIError(ExitConfiguration, "nao foi possivel interpretar o settings.xml", err)
	}
	if len(sections) > 1 {
		return nil, NewCLIError(ExitConflict, "settings.xml contem secoes servers ambiguas", nil)
	}
	newline := detectNewline(content)
	server := mavenServerXML(username, token, newline, "    ")
	if len(sections) == 1 {
		location := sections[0]
		section := content[location.start:location.end]
		selfClosing, err := rootElementSelfClosing(section, "servers")
		if err != nil {
			return nil, NewCLIError(ExitConfiguration, "nao foi possivel interpretar os servidores Maven", err)
		}
		if selfClosing {
			return nil, NewCLIError(ExitConflict, "settings.xml usa o elemento servers autocontido", nil)
		}
		updated, err := upsertServerInSection(section, server, newline)
		if err != nil {
			return nil, NewCLIError(ExitConfiguration, "nao foi possivel interpretar os servidores Maven", err)
		}
		return replaceRange(content, location.slice(), updated), nil
	}
	closing, err := rootClosingOffset(content, "settings")
	if err != nil {
		return nil, NewCLIError(ExitConfiguration, "nao foi possivel interpretar o settings.xml", err)
	}
	if closing < 0 {
		return nil, NewCLIError(ExitConflict, "settings.xml nao possui o elemento raiz settings", nil)
	}
	insertion := []byte("  <servers>" + newline + string(server) + newline + "  </servers>" + newline)
	return replaceRange(content, []int{closing, closing}, insertion), nil
}

func removeMavenServer(content []byte) ([]byte, error) {
	if len(bytes.TrimSpace(content)) == 0 {
		return content, nil
	}
	if err := validateSettingsXML(content); err != nil {
		return nil, NewCLIError(ExitConfiguration, "settings.xml nao contem XML valido", err)
	}
	sections, err := directChildRanges(content, "settings", "servers")
	if err != nil {
		return nil, NewCLIError(ExitConfiguration, "nao foi possivel interpretar o settings.xml", err)
	}
	if len(sections) > 1 {
		return nil, NewCLIError(ExitConflict, "settings.xml contem secoes servers ambiguas", nil)
	}
	if len(sections) == 0 {
		return content, nil
	}
	location := sections[0]
	section := content[location.start:location.end]
	updated, err := removeServerFromSection(section)
	if err != nil {
		return nil, NewCLIError(ExitConfiguration, "nao foi possivel interpretar os servidores Maven", err)
	}
	return replaceRange(content, location.slice(), updated), nil
}

func validateSettingsXML(content []byte) error {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	root := ""
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if start, ok := token.(xml.StartElement); ok && root == "" {
			root = start.Name.Local
		}
	}
	if root != "settings" {
		return errors.New("unexpected XML root")
	}
	return nil
}

func upsertServerInSection(section, server []byte, newline string) ([]byte, error) {
	matches, err := directChildRanges(section, "servers", "server")
	if err != nil {
		return nil, err
	}
	var result bytes.Buffer
	cursor := 0
	found := false
	for _, match := range matches {
		block := section[match.start:match.end]
		id, err := mavenServerID(block)
		if err != nil {
			return nil, err
		}
		if id != "github-boilerplate" {
			continue
		}
		result.Write(section[cursor:match.start])
		if !found {
			result.Write(bytes.TrimPrefix(server, []byte("    ")))
			found = true
		}
		cursor = match.end
	}
	if found {
		result.Write(section[cursor:])
		return result.Bytes(), nil
	}
	closing, err := rootClosingOffset(section, "servers")
	if err != nil {
		return nil, err
	}
	if closing < 0 {
		return nil, errors.New("servers closing element not found")
	}
	insertion := append([]byte(newline), server...)
	insertion = append(insertion, []byte(newline+"  ")...)
	return replaceRange(section, []int{closing, closing}, insertion), nil
}

func removeServerFromSection(section []byte) ([]byte, error) {
	matches, err := directChildRanges(section, "servers", "server")
	if err != nil {
		return nil, err
	}
	var result bytes.Buffer
	cursor := 0
	for _, match := range matches {
		block := section[match.start:match.end]
		id, err := mavenServerID(block)
		if err != nil {
			return nil, err
		}
		if id != "github-boilerplate" {
			continue
		}
		result.Write(section[cursor:match.start])
		cursor = match.end
	}
	if cursor == 0 {
		return section, nil
	}
	result.Write(section[cursor:])
	return result.Bytes(), nil
}

type elementRange struct {
	start int
	end   int
}

func (r elementRange) slice() []int {
	return []int{r.start, r.end}
}

func directChildRanges(content []byte, parent, child string) ([]elementRange, error) {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	path := make([]string, 0, 3)
	activeStart := -1
	var ranges []elementRange
	for {
		start := int(decoder.InputOffset())
		token, err := decoder.RawToken()
		if errors.Is(err, io.EOF) {
			return ranges, nil
		}
		if err != nil {
			return nil, err
		}
		switch value := token.(type) {
		case xml.StartElement:
			path = append(path, value.Name.Local)
			if len(path) == 2 && path[0] == parent && path[1] == child {
				activeStart = start
			}
		case xml.EndElement:
			if len(path) == 2 && path[0] == parent && path[1] == child && activeStart >= 0 {
				ranges = append(ranges, elementRange{start: activeStart, end: int(decoder.InputOffset())})
				activeStart = -1
			}
			if len(path) > 0 {
				path = path[:len(path)-1]
			}
		}
	}
}

func rootClosingOffset(content []byte, root string) (int, error) {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	depth := 0
	for {
		start := int(decoder.InputOffset())
		token, err := decoder.RawToken()
		if errors.Is(err, io.EOF) {
			return -1, nil
		}
		if err != nil {
			return -1, err
		}
		switch value := token.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			if depth == 1 && value.Name.Local == root {
				return start, nil
			}
			depth--
		}
	}
}

func rootElementSelfClosing(content []byte, root string) (bool, error) {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	for {
		start := int(decoder.InputOffset())
		token, err := decoder.RawToken()
		if errors.Is(err, io.EOF) {
			return false, errors.New("root element not found")
		}
		if err != nil {
			return false, err
		}
		if element, ok := token.(xml.StartElement); ok {
			if element.Name.Local != root {
				return false, errors.New("unexpected root element")
			}
			opening := bytes.TrimSpace(content[start:int(decoder.InputOffset())])
			return bytes.HasSuffix(opening, []byte("/>")), nil
		}
	}
}

func mavenServerID(content []byte) (string, error) {
	var server struct {
		ID string `xml:"id"`
	}
	if err := xml.Unmarshal(content, &server); err != nil {
		return "", err
	}
	return strings.TrimSpace(server.ID), nil
}

func mavenServerXML(username, token, newline, indent string) []byte {
	return []byte(indent + "<server>" + newline +
		indent + "  <id>github-boilerplate</id>" + newline +
		indent + "  <username>" + escapeXML(username) + "</username>" + newline +
		indent + "  <password>" + escapeXML(token) + "</password>" + newline +
		indent + "</server>")
}

func escapeXML(value string) string {
	var escaped bytes.Buffer
	_ = xml.EscapeText(&escaped, []byte(value))
	return escaped.String()
}

func upsertNPMCredentials(content []byte, token string) ([]byte, error) {
	newline := detectNewline(content)
	lines, err := removeNPMLines(content, true)
	if err != nil {
		return nil, err
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) > 0 {
		lines = append(lines, "")
	}
	lines = append(lines, npmManagedBegin, npmScopeLine, npmTokenPrefix+token, npmManagedEnd, "")
	return []byte(strings.Join(lines, newline)), nil
}

func removeNPMManagedBlock(content []byte) ([]byte, error) {
	if len(content) == 0 {
		return content, nil
	}
	newline := detectNewline(content)
	lines, err := removeNPMLines(content, false)
	if err != nil {
		return nil, err
	}
	return []byte(strings.Join(lines, newline)), nil
}

func removeNPMLines(content []byte, removeLegacy bool) ([]string, error) {
	normalized := strings.ReplaceAll(string(content), "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	result := make([]string, 0, len(lines))
	inManagedBlock := false
	managedBlocks := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == npmManagedBegin {
			if inManagedBlock || managedBlocks > 0 {
				return nil, NewCLIError(ExitConflict, ".npmrc contem marcadores gerenciados ambiguos", nil)
			}
			inManagedBlock = true
			managedBlocks++
			continue
		}
		if trimmed == npmManagedEnd && !inManagedBlock {
			return nil, NewCLIError(ExitConflict, ".npmrc contem marcador final sem inicio", nil)
		}
		if inManagedBlock {
			if trimmed == npmManagedEnd {
				inManagedBlock = false
			}
			continue
		}
		if removeLegacy && (strings.HasPrefix(trimmed, "@aplicacoesBoilerplate:registry=") || strings.HasPrefix(trimmed, npmTokenPrefix)) {
			continue
		}
		result = append(result, line)
	}
	if inManagedBlock {
		return nil, NewCLIError(ExitConflict, ".npmrc contem bloco gerenciado incompleto", nil)
	}
	return result, nil
}

func containsMavenServer(content []byte) bool {
	if validateSettingsXML(content) != nil {
		return false
	}
	var settings struct {
		Servers []struct {
			ID       string `xml:"id"`
			Username string `xml:"username"`
			Password string `xml:"password"`
		} `xml:"servers>server"`
	}
	if err := xml.Unmarshal(content, &settings); err != nil {
		return false
	}
	for _, server := range settings.Servers {
		if strings.TrimSpace(server.ID) == "github-boilerplate" &&
			strings.TrimSpace(server.Username) != "" && strings.TrimSpace(server.Password) != "" {
			return true
		}
	}
	return false
}

func containsNPMManagedBlock(content []byte) bool {
	if _, err := removeNPMLines(content, false); err != nil {
		return false
	}
	normalized := strings.ReplaceAll(string(content), "\r\n", "\n")
	inManagedBlock := false
	hasScope := false
	hasToken := false
	for _, line := range strings.Split(normalized, "\n") {
		trimmed := strings.TrimSpace(line)
		switch trimmed {
		case npmManagedBegin:
			inManagedBlock = true
		case npmManagedEnd:
			inManagedBlock = false
		default:
			if inManagedBlock {
				hasScope = hasScope || trimmed == npmScopeLine
				hasToken = hasToken || (strings.HasPrefix(trimmed, npmTokenPrefix) && strings.TrimSpace(strings.TrimPrefix(trimmed, npmTokenPrefix)) != "")
			}
		}
	}
	return hasScope && hasToken
}

func detectNewline(content []byte) string {
	if bytes.Contains(content, []byte("\r\n")) {
		return "\r\n"
	}
	return "\n"
}

func replaceRange(content []byte, location []int, replacement []byte) []byte {
	result := make([]byte, 0, len(content)-location[1]+location[0]+len(replacement))
	result = append(result, content[:location[0]]...)
	result = append(result, replacement...)
	result = append(result, content[location[1]:]...)
	return result
}
