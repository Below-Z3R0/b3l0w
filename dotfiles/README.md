# Dotfiles del perfil desktop

Lo que `b3l0w install system --profile=desktop --dotfiles .` copia al home.
Se aplican **con backup**: si el destino existe, queda como `<destino>.bak-<fecha>`.

## Qué hay

| Carpeta | Destino | Qué es |
|---|---|---|
| `fish/` | `~/.config/fish` | `config.fish` (prompt, aliases, greeting) |
| `fastfetch/` | `~/.config/fastfetch` | las tres configs del greeting (completa, min, micro) |
| `noctalia/` | `~/.config/noctalia` | `config.toml` del shell |
| `umbriel/` | `~/.config/umbriel` | `config.toml` (keybinds) + `noctalia.toml` |

## Qué NO se versiona a propósito

- `fish_variables` (estado universal: puede tener tokens y credenciales).
- `fish/completions/` (los generan bun, copilot, etc.).
- Los `*.bak-*`.
- El binario de `kitty-shapes` (172 MB) y el de `boot-preview`: los instala `b3l0w install`.
- `fastfetch/emblem/` (salida en runtime del greeting).
- El script `b3l0wz3r0-fastfetch`: lo instala `kitty-shapes install`.

## Rutas

Las rutas absolutas del autor se reemplazaron por `$HOME` donde el programa lo expande
(fish, keybinds de Umbriel). Las de `noctalia/config.toml` (carpetas de wallpapers)
quedan como ejemplo: cada máquina las ajusta.
