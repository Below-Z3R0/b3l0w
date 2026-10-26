// b3l0w es el CLI de control del sistema: versiones, estado, updates y
// extensiones de los programas propios.
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Below-Z3R0/b3l0w/internal/cli"
	"github.com/Below-Z3R0/b3l0w/internal/config"
	"github.com/Below-Z3R0/b3l0w/internal/doctor"
	"github.com/Below-Z3R0/b3l0w/internal/plugins"
	"github.com/Below-Z3R0/b3l0w/internal/programs"
	"github.com/Below-Z3R0/b3l0w/internal/update"
	"github.com/Below-Z3R0/b3l0w/internal/version"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		usage()
		return 0
	}

	command, rest := args[0], args[1:]

	switch command {
	case "help", "-h", "--help":
		usage()
		return 0
	case "version", "--version", "-V":
		fmt.Printf("b3l0w %s\n", version.Version)
		return 0
	case "list":
		return list()
	case "status":
		return withManifest(func(m *config.Manifest) int {
			return update.Run(m, update.Options{Check: true})
		})
	case "doctor":
		return withManifest(doctor.Run)
	case "update":
		return updateCommand(rest)
	case "program":
		return programCommand(rest)
	case "plugins":
		return pluginsCommand()
	}

	// Extensión: b3l0w <nombre> -> b3l0w-<nombre> del directorio privado.
	if plugins.Find(command) != "" {
		if err := plugins.Run(command, rest); err != nil {
			cli.Fail(err.Error())
			return 1
		}
		return 0
	}

	cli.Fail(fmt.Sprintf("comando desconocido: %s", command))
	fmt.Println()
	usage()
	return 1
}

// withManifest carga el manifiesto una vez y se lo pasa al comando.
func withManifest(fn func(*config.Manifest) int) int {
	manifest, err := config.Load(cli.ManifestPath())
	if err != nil {
		cli.Fail(fmt.Sprintf("no se pudo leer %s: %v", cli.ManifestPath(), err))
		return 1
	}
	return fn(manifest)
}

func updateCommand(args []string) int {
	opts := update.Options{}
	for _, arg := range args {
		switch arg {
		case "--system":
			opts.System = true
		case "--programs":
			opts.Programs = true
		case "--check":
			opts.Check = true
		case "-h", "--help":
			fmt.Println("uso: b3l0w update [--system|--programs|--check]")
			return 0
		default:
			cli.Fail(fmt.Sprintf("opción desconocida: %s", arg))
			return 1
		}
	}
	if opts.System && opts.Programs {
		cli.Fail("--system y --programs son excluyentes; sin banderas se hacen las dos")
		return 1
	}
	return withManifest(func(m *config.Manifest) int { return update.Run(m, opts) })
}

func list() int {
	return withManifest(func(m *config.Manifest) int {
		cli.Step("Programas registrados")
		if len(m.Programs) == 0 {
			fmt.Println("  (ninguno)")
			return 0
		}
		for _, status := range programs.Statuses(m) {
			switch {
			case !status.Installed:
				fmt.Printf("  %-14s %s\n", status.Name, cli.Dim("no instalado"))
			case status.Err != "":
				fmt.Printf("  %-14s %s\n", status.Name, cli.Dim("error: "+status.Err))
			default:
				fmt.Printf("  %-14s %-8s %s\n", status.Name, status.Version, cli.Dim(status.Bin))
			}
		}
		return 0
	})
}

func programCommand(args []string) int {
	if len(args) == 0 {
		fmt.Println("uso: b3l0w program <add|remove> <nombre> [--source github:owner/repo]")
		return 1
	}

	action := args[0]
	name := ""
	source := ""
	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "--source=") {
			source = strings.TrimPrefix(arg, "--source=")
			continue
		}
		if name == "" {
			name = arg
		}
	}
	if name == "" {
		cli.Fail("falta el nombre del programa")
		return 1
	}

	return withManifest(func(m *config.Manifest) int {
		switch action {
		case "add":
			if !m.Add(name, source) {
				cli.Warn(fmt.Sprintf("%s ya estaba registrado", name))
				return 0
			}
		case "remove":
			if !m.Remove(name) {
				cli.Fail(fmt.Sprintf("%s no está registrado", name))
				return 1
			}
		default:
			cli.Fail(fmt.Sprintf("acción desconocida: %s", action))
			return 1
		}

		if err := m.Save(cli.ManifestPath()); err != nil {
			cli.Fail(fmt.Sprintf("no se pudo guardar el manifiesto: %v", err))
			return 1
		}
		cli.Ok(fmt.Sprintf("%s: %s", action, name))
		return 0
	})
}

func pluginsCommand() int {
	names := plugins.Discover()
	cli.Step("Extensiones")
	if len(names) == 0 {
		fmt.Printf("  (ninguna en %s)\n", cli.PluginsDir())
		return 0
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("  b3l0w-%s\n", name)
	}
	return 0
}

func usage() {
	fmt.Printf(`b3l0w %s — control del sistema

uso: b3l0w <comando> [opciones]

  status                 sistema base (novedades) + versión de cada programa
  list                   programas registrados, versión y binario
  update                 actualiza el sistema (interactivo) y los programas
  update --system        solo el sistema base (cachy-update / arch-update)
  update --programs      solo los programas propios
  update --check         reporta, no toca nada
  doctor                 dependencias, manifiesto y directorio de plugins
  program add <nombre>   registra un programa
  program remove <nom>   lo saca del manifiesto
  plugins                lista las extensiones instaladas
  version                imprime la versión

Las extensiones son b3l0w-<nombre> en %s
(no se ejecutan desde el PATH).
`, version.Version, cli.PluginsDir())
}
