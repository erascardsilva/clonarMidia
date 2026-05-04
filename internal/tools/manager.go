package tools

import (
	"embed"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

//go:embed all:bin
var embeddedTools embed.FS

type Manager struct {
	basePath string
}

func NewManager() *Manager {
	// Usamos um diretório no home do usuário para persistência entre sessões
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".clonarmidia", "bin")
	return &Manager{basePath: path}
}

// Init extrai os binários embutidos para o sistema de arquivos
func (m *Manager) Init() error {
	if err := os.MkdirAll(m.basePath, 0755); err != nil {
		return err
	}

	// Lista arquivos no binário embutido
	entries, err := embeddedTools.ReadDir("bin")
	if err != nil {
		return nil // Se a pasta bin estiver vazia, não faz nada
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		srcPath := filepath.Join("bin", entry.Name())
		dstPath := filepath.Join(m.basePath, entry.Name())

		// Extrai apenas se não existir ou se quisermos forçar atualização
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			if err := m.extractFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Manager) extractFile(srcPath, dstPath string) error {
	src, err := embeddedTools.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// GetToolPath retorna o caminho absoluto da ferramenta
func (m *Manager) GetToolPath(name string) string {
	// Verifica primeiro na nossa pasta interna
	localPath := filepath.Join(m.basePath, name)
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}
	// Fallback para o PATH do sistema
	return name
}

// RunWithPrivileges executa um comando usando pkexec se necessário
func (m *Manager) RunWithPrivileges(name string, args ...string) *exec.Cmd {
	if runtime.GOOS == "linux" {
		fullArgs := append([]string{m.GetToolPath(name)}, args...)
		return exec.Command("pkexec", fullArgs...)
	}
	return exec.Command(m.GetToolPath(name), args...)
}

// Erasmo Cardoso - Software Engineer | Electronics Technician
