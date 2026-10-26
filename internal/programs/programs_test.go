package programs

import "testing"

func TestParseVersion(t *testing.T) {
	casos := []struct {
		in   string
		want string
	}{
		{"kitty-shapes v1.0.0", "v1.0.0"},
		{"boot-preview v0.2.2", "v0.2.2"},
		{"v1.2.3", "v1.2.3"},
		{"prog 10.20.30 extra", "v10.20.30"},
		{"sin version", ""},
		{"", ""},
	}
	for _, caso := range casos {
		if got := ParseVersion(caso.in); got != caso.want {
			t.Errorf("ParseVersion(%q) = %q, quiero %q", caso.in, got, caso.want)
		}
	}
}

func TestLastLine(t *testing.T) {
	if got := lastLine("uno\ndos\n\n"); got != "dos" {
		t.Fatalf("lastLine = %q, quiero %q", got, "dos")
	}
	if got := lastLine("   "); got != "sin salida" {
		t.Fatalf("lastLine vacío = %q", got)
	}
}
