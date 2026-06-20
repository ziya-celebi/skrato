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

The runtime installer/app needs:

- `love` (LÖVE; used to run the GUI)
- `pkexec` (polkit; for privileged maintenance actions)
- `rsync` (used by `installer.sh`)

## Install (most common)

This installs the LÖVE app bundle + a launcher script + a desktop entry.

### 1) Clone the repo

```bash
git clone https://github.com/<YOUR_USERNAME>/skrato.git
cd skrato
```

> If you don’t use GitHub, replace the clone URL with your source location.

### 2) Install skrato

From inside the repo directory (the one containing `installer.sh`):

```bash
bash installer.sh
```

### 3) Where it installs

After installation, skrato is placed in:

- Application bundle: `~/.skrato/`
- Launcher: `~/.local/bin/skrato`
- Desktop entry: `~/.local/share/applications/skrato.desktop`

If `~/.local/bin` is not in your PATH, add it:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Then **log out/in** or restart your shell.

## Run

### Option A: from a terminal

```bash
skrato
```

### Option B: from the desktop launcher

Use the app entry created at: `~/.local/share/applications/skrato.desktop`.

## Uninstall

From anywhere (the uninstall command is a self-contained action):

```bash
bash installer.sh uninstall
```

Removes:

- `~/.local/bin/skrato`
- `~/.local/share/applications/skrato.desktop`
- `~/.skrato/`

## Notes about what gets installed

`installer.sh` copies only the project files required to run the LÖVE app into `~/.skrato/`.

It does **not** install compiled Rust build artifacts (the `target/` directory).

