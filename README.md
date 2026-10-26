# b3l0w

CLI de control del sistema: actualiza **el sistema base** y **los programas
propios** desde un solo comando, y muestra qué versión tiene cada cosa.

```
b3l0w status                 # sistema base (novedades) + versión de cada programa
b3l0w list                   # programas registrados, versión y binario
b3l0w update                 # sistema (interactivo) y después los programas
b3l0w update --system        # solo cachy-update / arch-update
b3l0w update --programs      # solo los programas propios
b3l0w update --check         # reporta, no toca nada
b3l0w update --self          # actualiza el propio b3l0w
b3l0w install <nombre>       # baja el programa del release a ~/.local/bin
b3l0w doctor                 # dependencias, programas y directorio de plugins
```

## Cómo está armado

Un core chico que **delega**. No sabe de GitHub, de assets ni de pacman:
ejecuta el `update` de cada programa y reporta. Así cada programa sigue
siendo dueño de sus operaciones y no hay dos implementaciones que se
separen en el primer release.

```
b3l0w/
├── cmd/b3l0w/          # dispatcher
├── internal/cli/       # salida, rutas, ejecución de comandos
├── internal/config/    # manifiesto de programas (JSON)
├── internal/programs/  # versión y update de cada programa
├── internal/system/    # adaptador de cachy-update / arch-update
├── internal/update/    # orquestación y código de salida
├── internal/doctor/    # chequeos del entorno
└── internal/plugins/   # extensiones b3l0w-<nombre>
```

Sin dependencias: solo la biblioteca estándar de Go.

## Instalar un programa

```
b3l0w install kitty-shapes
```

Baja el asset del último release del repo que dice el manifiesto y lo
deja en `~/.local/bin`. Valida con `version` antes de pisar el que está.

El manifiesto declara el asset porque cada repo publica el suyo:

```json
{ "name": "kitty-shapes", "source": "github:Below-Z3R0/kitty-shapes", "asset": "kitty-shapes-linux-x64" }
```

`install` deja el **binario**. La configuración (configs de fastfetch,
enganche de fish) la hace el propio programa con su `install`.

## Actualizarse a sí mismo

```
b3l0w update --self
```

Es el caso especial del contrato: el core no puede delegar su propio
update. Baja su asset, valida la versión y se reemplaza con un rename
atómico. Si el binario pertenece a un paquete (`pacman -Qo`), no lo toca:
eso lo hace pacman.

## El contrato de cada programa

Un programa entra al core cuando cumple esto:

| Comando | Contrato |
|---|---|
| `prog version` | una línea: `prog vX.Y.Z` |
| `prog update` | se actualiza solo, sin preguntar, exit 0/1 |
| `prog doctor` | chequeos, exit != 0 si algo falta |
| `prog install [--dry-run] [--uninstall]` | idempotente, con backup |

Hoy lo cumplen `boot-preview` y `kitty-shapes`.

## El sistema base

`cachy-update` en CachyOS es un **symlink a `arch-update`**, y pertenece a
un paquete: no se renombra ni se reescribe. El CLI lo llama por nombre,
con fallback al upstream:

```
cachy-update -> arch-update
```

Su ejecución principal es **interactiva**: muestra las news, pide
confirmación y después ofrece huérfanos, `pacnew`, kernel pendiente y
servicios a reiniciar. No tiene `--yes`. Por eso:

- `status` / `update --check` usan `-l`, que no pregunta.
- `b3l0w update` le **entrega la terminal** en vez de fingir que lo
  automatiza.
- sin terminal, la mitad del sistema se rechaza con un mensaje claro y
  se sugiere `--programs`.

## Manifiesto

`~/.config/b3l0w/programs.json`:

```json
{
  "version": 1,
  "programs": [
    { "name": "boot-preview", "source": "github:Below-Z3R0/boot-preview" },
    { "name": "kitty-shapes", "source": "github:Below-Z3R0/kitty-shapes" }
  ]
}
```

Si no existe, se usa ese default. Se edita a mano o con
`b3l0w program add|remove`.

## Extensiones

Un binario `b3l0w-<nombre>` en `~/.local/lib/b3l0w/plugins/` (0700).
`b3l0w <nombre> ...` lo ejecuta.

**No se buscan en el PATH a propósito.** Cualquier cosa que pueda escribir
un `b3l0w-update` en un directorio del PATH lo haría correr con este
usuario, y este CLI pide sudo. `git` acepta ese riesgo; acá no.

## Roadmap

- **Etapa 2**: repo pacman `[b3l0w]` firmado y PKGBUILDs. Los programas
  pasan a ser paquetes y su `update` detecta el dueño con `pacman -Qo`.
- **Rollback y snapshots** (snapper): pendiente.
- **Checksums de los assets** (`SHA256SUMS`): pendiente.
- **Extensiones**: `b3l0w-containers`, cuando el setup de containers esté
  terminado.

## Compilar

```bash
go build -ldflags "-X github.com/Below-Z3R0/b3l0w/internal/version.Version=v0.1.0" -o b3l0w ./cmd/b3l0w
go test ./...
bash scripts/smoke-test.sh ./b3l0w
```
