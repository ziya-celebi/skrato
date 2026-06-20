# skrato

Linux utility (GUI) for maintaining bootloader + initramfs.

## Requirements

- `love` (LÖVE 11.x)
- `pkexec` (polkit)

## Install

```bash
bash installer.sh
```

This installs into:
- `~/.skrato/`
- launcher: `~/.local/bin/skrato`
- desktop entry: `~/.local/share/applications/skrato.desktop`

Make sure `~/.local/bin` is in your PATH (once):

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Log out/in or reload your shell after changing PATH.

## Run

```bash
skrato
```

## What it does

- Detects available maintenance tools on your system.
- Updates GRUB (`grub-mkconfig` / `grub2-mkconfig`) when present.
- Updates systemd-boot (`bootctl update`) when present.
- Refreshes rEFInd (`refind-install`) when present.
- Rebuilds initramfs (`mkinitcpio`, `dracut`, or `update-initramfs`) when present.

Maintenance commands are executed via `pkexec`.

## Uninstall

```bash
bash installer.sh uninstall
```

