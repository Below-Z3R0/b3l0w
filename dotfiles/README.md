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

## Noctalia: config vs state (y secretos)

**`dotfiles/noctalia/config.toml` es el export de la config efectiva**
(`noctalia config export`): incluye lo que la UI guarda en el state, no solo los
defaults. Se regenera con:

```bash
noctalia config export > dotfiles/noctalia/config.toml
```

Antes de commitear, **redactar credenciales**: el export puede traer
`[plugin_settings."…"] client_secret = "…"`. Se deja vacío y cada máquina lo
configura desde la UI. Nunca commitear el valor.

**`~/.local/state/noctalia/settings.toml` NO se versiona**: contiene estado de la
máquina (monitores, avatar, wallpapers, apps del dock) y credenciales de plugins.

## Noctalia: config vs state

`~/.config/noctalia/config.toml` son los **defaults**. La UI del shell (Settings)
no escribe ahí: guarda lo que cambiás en `~/.local/state/noctalia/settings.toml`,
que **pisa** al config. Por eso el archivo puede decir una cosa y la UI otra.

Acá se versiona solo `config.toml`, con los valores deseados como default. El
`settings.toml` es estado de la máquina (avatar, monitores, wallpapers) y no se
versiona. Si querés que un cambio de la UI quede como default de la distro, hay
que copiarlo a este `config.toml`.

## Rutas

Las rutas absolutas del autor se reemplazaron por `$HOME` donde el programa lo expande
(fish, keybinds de Umbriel). Las de `noctalia/config.toml` (carpetas de wallpapers)
quedan como ejemplo: cada máquina las ajusta.
