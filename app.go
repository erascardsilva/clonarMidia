package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"clonarmidia/internal/tools"
	"clonarmidia/pkg/disk"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	diskService  *disk.Service
	toolsManager *tools.Manager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		diskService:  disk.NewService(),
		toolsManager: tools.NewManager(),
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	_ = a.toolsManager.Init()
}

// IsRoot verifica se a aplicação está rodando com privilégios de root
func (a *App) IsRoot() bool {
	return os.Geteuid() == 0
}

// GetDisks retorna a lista de discos disponíveis
func (a *App) GetDisks() ([]disk.Info, error) {
	return a.diskService.GetDisks()
}

// StartClone inicia o processo de clonagem usando ferramentas externas com privilégios
func (a *App) StartClone(opts disk.CloneOptions, password string) string {
	go func() {
		bs := fmt.Sprintf("%d", opts.BufferSize)
		if opts.BufferSize <= 0 {
			bs = "64M"
		}

		// Usamos sudo -S para injetar a senha via entrada padrão
		cmd := exec.Command("sudo", "-S", "dd", 
			fmt.Sprintf("if=%s", opts.Source), 
			fmt.Sprintf("of=%s", opts.Destination), 
			fmt.Sprintf("bs=%s", bs),
			"status=progress",
			"conv=fsync",
		)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			runtime.EventsEmit(a.ctx, "clone_error", "Falha ao criar pipe de entrada: "+err.Error())
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			runtime.EventsEmit(a.ctx, "clone_error", "Falha ao iniciar pipe: "+err.Error())
			return
		}

		if err := cmd.Start(); err != nil {
			runtime.EventsEmit(a.ctx, "clone_error", "Falha ao disparar comando: "+err.Error())
			return
		}

		// Injeta a senha no stdin
		fmt.Fprintln(stdin, password)
		stdin.Close()

		// Goroutine para parsear o progresso do dd
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stderr.Read(buf)
				if n > 0 {
					output := string(buf[:n])
					runtime.EventsEmit(a.ctx, "clone_log", output)
				}
				if err != nil {
					break
				}
			}
		}()

		err = cmd.Wait()
		if err != nil {
			runtime.EventsEmit(a.ctx, "clone_error", "Clonagem interrompida ou falhou: "+err.Error())
			return
		}
		
		runtime.EventsEmit(a.ctx, "clone_complete", "Sucesso")
	}()

	return "Clonagem iniciada"
}

// ScanPartitions realiza uma busca por partições perdidas usando o TestDisk
func (a *App) ScanPartitions(device string) string {
	go func() {
		// testdisk /cmd /dev/sdX analyze,search
		cmd := a.toolsManager.RunWithPrivileges("testdisk", "/cmd", device, "analyze,search")
		
		output, err := cmd.CombinedOutput()
		if err != nil {
			runtime.EventsEmit(a.ctx, "recovery_error", "Falha no Scan: "+err.Error())
			return
		}
		runtime.EventsEmit(a.ctx, "recovery_result", string(output))
	}()
	return "Escaneamento iniciado"
}

// RecoverFiles inicia a extração de arquivos usando o PhotoRec
func (a *App) RecoverFiles(device string, outputDir string) string {
	go func() {
		// photorec /d [output] /cmd [device] options...
		cmd := a.toolsManager.RunWithPrivileges("photorec", "/d", outputDir, "/cmd", device, "search")
		
		// O photorec emite muito log, capturaremos o progresso básico
		stdout, _ := cmd.StdoutPipe()
		_ = cmd.Start()
		
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				runtime.EventsEmit(a.ctx, "recovery_log", string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
		_ = cmd.Wait()
		runtime.EventsEmit(a.ctx, "recovery_complete", "Recuperação finalizada")
	}()
	return "Recuperação de arquivos iniciada"
}

// RepairFS tenta reparar um sistema de arquivos usando fsck
func (a *App) RepairFS(device string) string {
	go func() {
		// fsck -y /dev/sdX
		cmd := a.toolsManager.RunWithPrivileges("fsck", "-y", device)
		
		output, err := cmd.CombinedOutput()
		if err != nil {
			runtime.EventsEmit(a.ctx, "recovery_error", "Falha no reparo: "+err.Error())
			return
		}
		runtime.EventsEmit(a.ctx, "recovery_complete", "Reparo concluído:\n"+string(output))
	}()
	return "Reparo de disco iniciado"
}

// ElevatePrivileges tenta validar a senha de root para mudar o estado da aplicação
func (a *App) ElevatePrivileges(password string) bool {
	// Tenta executar um comando simples com sudo -S
	cmd := exec.Command("sudo", "-S", "true")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return false
	}

	go func() {
		defer stdin.Close()
		fmt.Fprintln(stdin, password)
	}()

	err = cmd.Run()
	return err == nil
}

// Erasmo Cardoso - Software Engineer | Electronics Technician
