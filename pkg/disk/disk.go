/*
Clonar Mídia - Gerenciamento de Discos e Partições
Criado por Erasmo Cardoso - Software Engineer | Electronics Technician
*/

package disk

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

// Info representa as informações de um disco ou partição
type Info struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        uint64 `json:"size"`
	Type        string `json:"type"` // "disk", "part", "loop", etc
	MountPoint  string `json:"mountpoint"`
	FSType      string `json:"fstype"`
	Model       string `json:"model"`
	Serial      string `json:"serial"`
	Health      string `json:"health"`
	Partitions  []Info `json:"partitions,omitempty"`
}

// Service gerencia as operações de disco
type Service struct{}

// NewService cria uma nova instância do serviço de disco
func NewService() *Service {
	return &Service{}
}

// GetDisks retorna a lista de discos disponíveis no sistema
func (s *Service) GetDisks() ([]Info, error) {
	switch runtime.GOOS {
	case "linux":
		return s.getLinuxDisks()
	case "windows":
		return s.getWindowsDisks()
	default:
		return nil, fmt.Errorf("sistema operacional não suportado: %s", runtime.GOOS)
	}
}

// getLinuxDisks usa o comando lsblk para obter informações dos discos
func (s *Service) getLinuxDisks() ([]Info, error) {
	// lsblk -J -b -o NAME,PATH,SIZE,TYPE,MOUNTPOINT,MODEL,SERIAL,FSTYPE
	cmd := exec.Command("lsblk", "-J", "-b", "-o", "NAME,PATH,SIZE,TYPE,MOUNTPOINT,MODEL,SERIAL,FSTYPE")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar lsblk: %w", err)
	}

	var result struct {
		BlockDevices []struct {
			Name       string `json:"name"`
			Path       string `json:"path"`
			Size       uint64 `json:"size"`
			Type       string `json:"type"`
			MountPoint string `json:"mountpoint"`
			FSType     string `json:"fstype"`
			Model      string `json:"model"`
			Serial     string `json:"serial"`
			Children   []struct {
				Name       string `json:"name"`
				Path       string `json:"path"`
				Size       uint64 `json:"size"`
				Type       string `json:"type"`
				MountPoint string `json:"mountpoint"`
				FSType     string `json:"fstype"`
			} `json:"children"`
		} `json:"blockdevices"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("erro ao processar JSON do lsblk: %w", err)
	}

	var disks []Info
	for _, dev := range result.BlockDevices {
		// Filtramos para pegar apenas discos reais (não loops ou roms se preferir)
		if dev.Type == "disk" {
			disk := Info{
				Name:       dev.Name,
				Path:       dev.Path,
				Size:       dev.Size,
				Type:       dev.Type,
				MountPoint: dev.MountPoint,
				FSType:     dev.FSType,
				Model:      dev.Model,
				Serial:     dev.Serial,
				Health:     "N/A", // Implementaremos análise depois
			}

			for _, part := range dev.Children {
				disk.Partitions = append(disk.Partitions, Info{
					Name:       part.Name,
					Path:       part.Path,
					Size:       part.Size,
					Type:       part.Type,
					MountPoint: part.MountPoint,
					FSType:     part.FSType,
				})
			}
			disks = append(disks, disk)
		}
	}

	return disks, nil
}

// getWindowsDisks (Placeholder para implementação futura)
func (s *Service) getWindowsDisks() ([]Info, error) {
	// Aqui usaremos Get-PhysicalDisk e Get-Partition via PowerShell ou WMI
	return []Info{}, nil
}

// Erasmo Cardoso - Software Engineer | Electronics Technician
