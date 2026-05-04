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

// Clone realiza a cópia bruta de um dispositivo para outro
func (s *Service) Clone(ctx context.Context, opts CloneOptions, progressChan chan<- Progress) error {
	src, err := os.OpenFile(opts.Source, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("falha ao abrir origem: %w", err)
	}
	defer src.Close()

	// Obtendo o tamanho total para progresso
	srcInfo, err := src.Stat()
	var totalSize uint64
	if err == nil && srcInfo.Size() > 0 {
		totalSize = uint64(srcInfo.Size())
	} else {
		// Alguns arquivos de dispositivo não retornam Stat() comum, tentaremos buscar via seek
		pos, _ := src.Seek(0, io.SeekEnd)
		totalSize = uint64(pos)
		src.Seek(0, io.SeekStart)
	}

	dst, err := os.OpenFile(opts.Destination, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("falha ao abrir destino: %w", err)
	}
	defer dst.Close()

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

				// Envia progresso a cada 100ms ou se terminar para não sobrecarregar o canal
				if time.Since(lastUpdate) >= 100*time.Millisecond || err == io.EOF {
					if progressChan != nil {
						elapsed := time.Since(startTime).Seconds()
						var speed float64
						if elapsed > 0 {
							speed = float64(bytesCopied) / elapsed
						}

						p := Progress{
							BytesCopied: bytesCopied,
							TotalBytes:  totalSize,
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
