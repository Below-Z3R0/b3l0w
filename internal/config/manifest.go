// Package config lee y escribe el manifiesto de programas.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Schema es la versión del formato del manifiesto.
const Schema = 1

// Program es un programa propio que el core sabe consultar y actualizar.
type Program struct {
	Name   string `json:"name"`
	Source string `json:"source,omitempty"`
}

// Manifest es la lista de programas registrados.
type Manifest struct {
	Version  int       `json:"version"`
	Programs []Program `json:"programs"`
}

// Default es el manifiesto que se usa la primera vez.
func Default() *Manifest {
	return &Manifest{
		Version: Schema,
		Programs: []Program{
			{Name: "boot-preview", Source: "github:Below-Z3R0/boot-preview"},
			{Name: "kitty-shapes", Source: "github:Below-Z3R0/kitty-shapes"},
		},
	}
}

// Load lee el manifiesto. Si no existe devuelve el default: el CLI
// tiene que funcionar recién instalado, sin un paso de configuración.
func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Default(), nil
	}
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if m.Version == 0 {
		m.Version = Schema
	}
	return &m, nil
}

// Save escribe el manifiesto creando el directorio si hace falta.
func (m *Manifest) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

// Get devuelve el programa con ese nombre, o nil.
func (m *Manifest) Get(name string) *Program {
	for i := range m.Programs {
		if m.Programs[i].Name == name {
			return &m.Programs[i]
		}
	}
	return nil
}

// Add registra un programa. Devuelve false si ya estaba: es idempotente.
func (m *Manifest) Add(name, source string) bool {
	if m.Get(name) != nil {
		return false
	}
	m.Programs = append(m.Programs, Program{Name: name, Source: source})
	return true
}

// Remove borra un programa. Devuelve false si no estaba.
func (m *Manifest) Remove(name string) bool {
	for i := range m.Programs {
		if m.Programs[i].Name == name {
			m.Programs = append(m.Programs[:i], m.Programs[i+1:]...)
			return true
		}
	}
	return false
}
