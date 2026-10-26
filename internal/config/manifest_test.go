package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSinArchivoDevuelveDefault(t *testing.T) {
	m, err := Load(filepath.Join(t.TempDir(), "no-existe.json"))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(m.Programs) == 0 {
		t.Fatal("el default no trae programas")
	}
	if m.Get("boot-preview") == nil || m.Get("kitty-shapes") == nil {
		t.Fatalf("faltan los programas por defecto: %+v", m.Programs)
	}
}

func TestSaveLoadRoundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "programs.json")
	m := Default()
	m.Add("otro", "github:algo/otro")
	if err := m.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("no se creó el archivo: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Programs) != len(m.Programs) {
		t.Fatalf("programas = %d, quiero %d", len(loaded.Programs), len(m.Programs))
	}
	if loaded.Get("otro") == nil {
		t.Fatal("se perdió el programa agregado")
	}
}

func TestAddEsIdempotente(t *testing.T) {
	m := Default()
	if m.Add("boot-preview", "x") {
		t.Fatal("agregar uno existente devolvió true")
	}
	if m.Add("nuevo", "x") != true {
		t.Fatal("agregar uno nuevo devolvió false")
	}
	if m.Add("nuevo", "y") {
		t.Fatal("agregar dos veces devolvió true la segunda")
	}
}

func TestRemove(t *testing.T) {
	m := Default()
	if !m.Remove("kitty-shapes") {
		t.Fatal("no borró un programa existente")
	}
	if m.Get("kitty-shapes") != nil {
		t.Fatal("el programa sigue ahí")
	}
	if m.Remove("kitty-shapes") {
		t.Fatal("borrar de nuevo devolvió true")
	}
}
