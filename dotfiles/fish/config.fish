source /usr/share/cachyos-fish-config/cachyos-config.fish

# Custom greeting — fastfetch con logo (imagen) + info del sistema
# Logo: ~/.config/fastfetch/logo.png (renderizado por fastfetch vía kitty/sixel)
# Info: ~/.config/fastfetch/config.jsonc
# Colores: sincronizados con Noctalia vía el template builtin "fastfetch"
function fish_greeting
    if test -x ~/.config/fastfetch/b3l0wz3r0-fastfetch
        ~/.config/fastfetch/b3l0wz3r0-fastfetch
    else if command -v fastfetch >/dev/null
        fastfetch --config ~/.config/fastfetch/config.jsonc
    end
end
# Rerender para Kitty
set -gx FASTFETCH_RERENDER 1

# Hermes Agent — ensure ~/.local/bin is on PATH
fish_add_path "$HOME/.local/bin"

# bun
set --export BUN_INSTALL "$HOME/.bun"
set --export PATH $BUN_INSTALL/bin $PATH

# pnpm
set -gx PNPM_HOME "$HOME/.local/share/pnpm"
if not string match -q -- "$PNPM_HOME/bin" $PATH
  set -gx PATH "$PNPM_HOME/bin" $PATH
end
# pnpm end

# bun
set --export BUN_INSTALL "$HOME/.bun"
set --export PATH $BUN_INSTALL/bin $PATH
set -gx ELECTRON_DISABLE_GPU 0
# kitty-shapes — greeting con emblema (instalado por kitty-shapes)
function fish_greeting
    if test -x "$HOME/.config/fastfetch/b3l0wz3r0-fastfetch"
        "$HOME/.config/fastfetch/b3l0wz3r0-fastfetch"
    else if command -v fastfetch >/dev/null
        fastfetch --config "$HOME/.config/fastfetch/config.jsonc"
    end
end
set -gx FASTFETCH_RERENDER 1
