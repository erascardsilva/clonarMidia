/*
Clonar Mídia - Cliente UDisks2 via D-Bus
Criado por Erasmo Cardoso - Software Engineer | Electronics Technician
*/

package udisks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/godbus/dbus/v5"
)

const (
	udisks2Service    = "org.freedesktop.UDisks2"
	udisks2BlockBase  = "/org/freedesktop/UDisks2/block_devices"
	udisks2BlockIface = "org.freedesktop.UDisks2.Block"
)

// Client gerencia a conexão com o UDisks2 via D-Bus
type Client struct {
	conn *dbus.Conn
}

// NewClient cria um novo cliente conectado ao D-Bus do sistema
func NewClient() (*Client, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao D-Bus do sistema: %w", err)
	}
	return &Client{conn: conn}, nil
}

// Close encerra a conexão com o D-Bus
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// deviceToObjectPath converte /dev/sda para o caminho D-Bus do UDisks2
func deviceToObjectPath(devicePath string) dbus.ObjectPath {
	name := filepath.Base(devicePath)
	// UDisks2 usa underscores para hifens e outros caracteres especiais
	name = strings.ReplaceAll(name, "-", "_")
	return dbus.ObjectPath(udisks2BlockBase + "/" + name)
}

// OpenForRead abre um dispositivo para leitura via UDisks2
// Retorna um *os.File pronto para leitura sequencial
func (c *Client) OpenForRead(devicePath string) (*os.File, error) {
	objPath := deviceToObjectPath(devicePath)
	obj := c.conn.Object(udisks2Service, objPath)

	var fd dbus.UnixFD
	options := map[string]dbus.Variant{}

	err := obj.Call(udisks2BlockIface+".OpenForBackup", 0, options).Store(&fd)
	if err != nil {
		return nil, fmt.Errorf("UDisks2 OpenForBackup falhou em %s: %w", devicePath, err)
	}

	file := os.NewFile(uintptr(fd), devicePath)
	if file == nil {
		return nil, fmt.Errorf("falha ao criar os.File a partir do file descriptor UDisks2")
	}
	return file, nil
}

// OpenForWrite abre um dispositivo para escrita via UDisks2
// Retorna um *os.File pronto para escrita sequencial
func (c *Client) OpenForWrite(devicePath string) (*os.File, error) {
	objPath := deviceToObjectPath(devicePath)
	obj := c.conn.Object(udisks2Service, objPath)

	var fd dbus.UnixFD
	options := map[string]dbus.Variant{}

	err := obj.Call(udisks2BlockIface+".OpenForRestore", 0, options).Store(&fd)
	if err != nil {
		return nil, fmt.Errorf("UDisks2 OpenForRestore falhou em %s: %w", devicePath, err)
	}

	file := os.NewFile(uintptr(fd), devicePath)
	if file == nil {
		return nil, fmt.Errorf("falha ao criar os.File a partir do file descriptor UDisks2")
	}
	return file, nil
}

// Erasmo Cardoso - Software Engineer | Electronics Technician
