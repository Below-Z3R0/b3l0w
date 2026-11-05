# b3l0w

Distribución de **b3l0wz3r0-os**: la definición del sistema, los dotfiles y el manifiesto de plugins.

| Ruta | Qué |
|---|---|
| `system.json` | perfiles (`desktop`, `dev`, `gaming`, `tools`, `docs`) |
| `dotfiles/` | configs versionadas (`fish`, `noctalia`, `umbriel`, `fastfetch`) |
| `plugins/plugins.json` | manifiesto de los 8 plugins de Noctalia |

El CLI vive en [b3l0w-cli](https://github.com/Below-Z3R0/b3l0w-cli).

## Uso

```bash
b3l0w install system --dry-run
b3l0w update --dotfiles --dry-run
```
