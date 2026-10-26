package system

import "testing"

// Fixture del formato real de `cachy-update -l`.
const listado = `==> Looking for updates...

==> Packages:
7zip                                   26.04-1.1                   ->  26.05-1.1
alsa-card-profiles                     1:1.6.9-1                   ->  1:1.6.10-1
breeze                                 6.8.0-1.10

==> AUR packages:
visual-studio-code-bin                 1.99.0-1                    ->  1.100.0-1
`

func TestPendingSeparaRepoYAUR(t *testing.T) {
	repo, aur := Pending(listado)
	if repo != 3 {
		t.Fatalf("repo = %d, quiero 3", repo)
	}
	if aur != 1 {
		t.Fatalf("aur = %d, quiero 1", aur)
	}
}

func TestPendingSinUpdates(t *testing.T) {
	repo, aur := Pending("==> Looking for updates...\n\n==> Packages:\n")
	if repo != 0 || aur != 0 {
		t.Fatalf("esperaba 0 y 0, obtuve %d y %d", repo, aur)
	}
}

func TestPendingIgnoraEncabezados(t *testing.T) {
	repo, aur := Pending("==> Packages:\n==> AUR packages:\n")
	if repo != 0 || aur != 0 {
		t.Fatalf("contó encabezados: %d y %d", repo, aur)
	}
}
