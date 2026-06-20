# skrato

Linux GUI for maintaining your bootloader and initramfs.

## What it does

- Detects available maintenance tools on your system.
- Can update:
  - GRUB (`grub-mkconfig` / `grub2-mkconfig`)
  - systemd-boot (`bootctl update`)
  - rEFInd (`refind-install`)
  - initramfs (`mkinitcpio` / `dracut` / `update-initramfs`)
- Executes maintenance commands via `pkexec`.

## Requirements

- `love` (LÖVE; used to run the GUI)
- `pkexec` (polkit; for privileged maintenance actions)
- `rsync` (used by the installer script)

## Install

```bash
bash installer.sh
```

Installs into:
- Application bundle: `~/.skrato/`
- Launcher: `~/.local/bin/skrato`
- Desktop entry: `~/.local/share/applications/skrato.desktop`

If `~/.local/bin` is not in your PATH:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

After changing PATH, log out/in or restart your shell.

## Run

From a terminal:

```bash
skrato
```

Or use your desktop launcher.

## Uninstall

```bash
bash installer.sh uninstall
```

Removes:
- `~/.local/bin/skrato`
- `~/.local/share/applications/skrato.desktop`
- `~/.skrato/`

## Notes about what gets installed

The installer copies the project files required to run the LÖVE app into `~/.skrato/`.
It does **not** install compiled Rust build artifacts (`target/`).

