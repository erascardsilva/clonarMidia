/*
Clonar Mídia - Motor de Clonagem Bit-a-Bit
Criado por Erasmo Cardoso - Software Engineer | Electronics Technician
*/

package disk

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"clonarmidia/pkg/udisks"
)

// Progress representa o estado atual da clonagem
type Progress struct {
	BytesCopied uint64  `json:"bytesCopied"`
	TotalBytes  uint64  `json:"totalBytes"`
	Percentage  float64 `json:"percentage"`
	Speed       float64 `json:"speed"` // Bytes por segundo
}

// CloneOptions define as configurações para a clonagem
type CloneOptions struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	BufferSize  int    `json:"bufferSize"` // em Bytes
}

// isSnapEnvironment verifica se o app está rodando dentro de um container Snap
func isSnapEnvironment() bool {
	return os.Getenv("SNAP") != ""
}

// openSource abre o dispositivo de origem pelo método adequado ao ambiente
func openSource(path string) (*os.File, func(), error) {
	if isSnapEnvironment() {
		client, err := udisks.NewClient()
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao conectar ao UDisks2: %w", err)
		}
		f, err := client.OpenForRead(path)
		if err != nil {
			client.Close()
			return nil, nil, fmt.Errorf("UDisks2 falhou ao abrir origem: %w", err)
		}
		return f, func() { f.Close(); client.Close() }, nil
	}
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao abrir origem: %w", err)
	}
	return f, func() { f.Close() }, nil
}

// openDest abre o dispositivo de destino pelo método adequado ao ambiente
func openDest(path string) (*os.File, func(), error) {
	if isSnapEnvironment() {
		client, err := udisks.NewClient()
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao conectar ao UDisks2: %w", err)
		}
		f, err := client.OpenForWrite(path)
		if err != nil {
			client.Close()
			return nil, nil, fmt.Errorf("UDisks2 falhou ao abrir destino: %w", err)
		}
		return f, func() { f.Close(); client.Close() }, nil
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao abrir destino: %w", err)
	}
	return f, func() { f.Close() }, nil
}

// Clone realiza a cópia bruta de um dispositivo para outro
func (s *Service) Clone(ctx context.Context, opts CloneOptions, progressChan chan<- Progress) error {
	src, closeSrc, err := openSource(opts.Source)
	if err != nil {
		return err
	}
	defer closeSrc()

	totalSize, _ := src.Seek(0, io.SeekEnd)
	src.Seek(0, io.SeekStart)

	dst, closeDst, err := openDest(opts.Destination)
	if err != nil {
		return err
	}
	defer closeDst()

	if opts.BufferSize <= 0 {
		opts.BufferSize = 1024 * 1024 // Default 1MB
	}

	buffer := make([]byte, opts.BufferSize)
	var bytesCopied uint64
	startTime := time.Now()
	lastUpdate := time.Now()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := src.Read(buffer)
			if n > 0 {
				_, wErr := dst.Write(buffer[:n])
				if wErr != nil {
					return fmt.Errorf("erro na escrita: %w", wErr)
				}
				bytesCopied += uint64(n)

				if time.Since(lastUpdate) >= 100*time.Millisecond || err == io.EOF {
					if progressChan != nil {
						elapsed := time.Since(startTime).Seconds()
						var speed float64
						if elapsed > 0 {
							speed = float64(bytesCopied) / elapsed
						}
						p := Progress{
							BytesCopied: bytesCopied,
							TotalBytes:  uint64(totalSize),
							Speed:       speed,
						}
						if totalSize > 0 {
							p.Percentage = (float64(bytesCopied) / float64(totalSize)) * 100
						}
						progressChan <- p
					}
					lastUpdate = time.Now()
				}
			}
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("erro na leitura: %w", err)
			}
		}
	}
}

// Erasmo Cardoso - Software Engineer | Electronics Technician
