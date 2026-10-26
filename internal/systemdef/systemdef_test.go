package systemdef

import "testing"

func TestBuiltinTieneLosPerfiles(t *testing.T) {
	d, err := Builtin()
	if err != nil {
		t.Fatalf("builtin: %v", err)
	}
	for _, name := range []string{"desktop", "dev", "gaming", "tools"} {
		if _, ok := d.Profiles[name]; !ok {
			t.Fatalf("falta el perfil %s", name)
		}
	}
}

func TestDockerVaEnDesktop(t *testing.T) {
	d, _ := Builtin()
	desktop := d.Profiles["desktop"]
	found := false
	for _, pkg := range desktop.Repo {
		if pkg == "docker" {
			found = true
		}
	}
	if !found {
		t.Fatal("docker tiene que estar en desktop")
	}
	// y no duplicado en dev
	for _, pkg := range d.Profiles["dev"].Repo {
		if pkg == "docker" {
			t.Fatal("docker está duplicado en dev")
		}
	}
}

func TestMergeSinDuplicados(t *testing.T) {
	d, _ := Builtin()
	merged, err := d.Merge([]string{"desktop", "desktop"})
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	seen := map[string]int{}
	for _, pkg := range merged.Repo {
		seen[pkg]++
	}
	for pkg, n := range seen {
		if n > 1 {
			t.Fatalf("%s aparece %d veces", pkg, n)
		}
	}
}

func TestMergePerfilDesconocido(t *testing.T) {
	d, _ := Builtin()
	if _, err := d.Merge([]string{"no-existe"}); err == nil {
		t.Fatal("esperaba error con un perfil que no existe")
	}
}

func TestMergeJuntaPerfiles(t *testing.T) {
	d, _ := Builtin()
	merged, err := d.Merge([]string{"desktop", "dev"})
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	has := func(list []string, want string) bool {
		for _, item := range list {
			if item == want {
				return true
			}
		}
		return false
	}
	if !has(merged.Repo, "docker") || !has(merged.Repo, "go") {
		t.Fatalf("faltan paquetes de los dos perfiles: %v", merged.Repo)
	}
}
