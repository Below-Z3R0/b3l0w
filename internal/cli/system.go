package cli

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Run ejecuta un comando y devuelve su stdout.
func Run(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	return strings.TrimRight(string(out), "\n"), err
}

// RunCombined ejecuta un comando con stderr mezclado.
func RunCombined(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return strings.TrimRight(string(out), "\n"), err
}

// RunTimeout ejecuta un comando y lo mata si tarda más de lo indicado.
//
// Los clientes que hablan con un daemon pueden quedarse esperando una
// respuesta que no llega: sin deadline, el comando hereda esa espera.
func RunTimeout(d time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	return strings.TrimRight(string(out), "\n"), err
}

// Interactive ejecuta un comando heredando la terminal: es lo que
// necesita un actualizador que pregunta.
func Interactive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// HaveCommand dice si un binario está en el PATH.
func HaveCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// CommandPath devuelve la ruta de un binario, o "" si no está.
func CommandPath(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return path
}

// IsTTY dice si stdin es una terminal.
//
// Un actualizador interactivo lanzado sin terminal falla o se cuelga,
// y el fallo es silencioso: hay que preguntarlo antes.
func IsTTY() bool {
	fi, err := os.Stdin.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}

// PackageOwner devuelve el paquete que provee un archivo, o "".
//
// Sirve para la etapa en que los programas propios sean paquetes: un
// binario con dueño no se auto-actualiza, se actualiza con pacman.
func PackageOwner(path string) string {
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	out, err := RunCombined("pacman", "-Qo", path)
	if err != nil {
		return ""
	}
	// "ruta is owned by paquete version"
	fields := strings.Fields(out)
	for i, field := range fields {
		if field == "by" && i+1 < len(fields) {
			return fields[i+1]
		}
	}
	return ""
}

// ExpandHome cambia un ~ inicial por el home del usuario.
func ExpandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, strings.TrimPrefix(path, "~"))
}

// ReplaceFile reemplaza dest con src de forma atómica.
//
// Copia a un temporal en el mismo directorio y hace rename: así no se
// pisa un binario que está corriendo ("text file busy") y no queda un
// archivo a medias si algo falla.
func ReplaceFile(src, dest string) error {
	// Si el destino es un symlink (por ejemplo ~/.local/bin/prog ->
	// ~/.config/prog/prog), se reemplaza el archivo real: si no, quedaría
	// el symlink convertido en una copia y el binario duplicado.
	if resolved, err := filepath.EvalSymlinks(dest); err == nil {
		dest = resolved
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	tmp := dest + ".new"
	if err := os.WriteFile(tmp, data, 0o755); err != nil {
		return err
	}
	return os.Rename(tmp, dest)
}
