package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func AddMavenDependency(groupID, artifactID, version string) error {
	// 1. Localiza o pom.xml no diretório onde a CLI foi executada
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %v", err)
	}
	pomPath := filepath.Join(cwd, "pom.xml")

	// 2. Lê o arquivo inteiro como texto
	fileBytes, err := os.ReadFile(pomPath)
	if err != nil {
		return fmt.Errorf("erro ao ler o pom.xml (ele existe neste diretório?): %v", err)
	}
	pomContent := string(fileBytes)

	// 3. Verifica se o pacote já está instalado para evitar duplicação
	artifactTag := fmt.Sprintf("<artifactId>%s</artifactId>", artifactID)
	if strings.Contains(pomContent, artifactTag) {
		fmt.Printf("O pacote %s já está no pom.xml!\n", artifactID)
		return nil
	}

	// 4. Monta o bloco XML (mantendo espaços para a indentação padrão)
	newDependency := fmt.Sprintf(`
		<dependency>
			<groupId>%s</groupId>
			<artifactId>%s</artifactId>
			<version>%s</version>
		</dependency>
	</dependencies>`, groupID, artifactID, version)

	// 5. Encontra a última tag </dependencies>
	targetTag := "</dependencies>"
	insertionIndex := strings.LastIndex(pomContent, targetTag)

	if insertionIndex == -1 {
		return fmt.Errorf("tag </dependencies> não encontrada no pom.xml")
	}

	// 6. Injeta a nova dependência no local correto usando slicing
	// Strings no Go podem ser concatenadas diretamente
	newPomContent := pomContent[:insertionIndex] + 
		strings.TrimLeft(newDependency, "\n\t ") + 
		pomContent[insertionIndex+len(targetTag):]

	// 7. Sobrescreve o arquivo com o conteúdo atualizado
	err = os.WriteFile(pomPath, []byte(newPomContent), 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar o pom.xml: %v", err)
	}

	fmt.Printf("✅ %s adicionado com sucesso!\n", artifactID)
	return nil
}

func main() {
	// Exemplo de uso da função dentro do seu comando
	err := AddMavenDependency("com.meus.pacotes", "minha-lib-core", "1.0.0")
	if err != nil {
		fmt.Println("❌ Falha na CLI:", err)
	}
}
