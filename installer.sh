#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  bash installer.sh [install|uninstall]

Default: install
EOF
}

ACTION="${1:-install}"
case "$ACTION" in
  install|uninstall) ;;
  -h|--help) usage; exit 0 ;;
  *) echo "Unknown action: $ACTION"; usage; exit 1 ;;
esac

HOME_DIR="${HOME}"
TARGET_DIR="$HOME_DIR/.skrato"
LAUNCHER_DIR="$HOME_DIR/.local/bin"
LAUNCHER="$LAUNCHER_DIR/skrato"
DESKTOP_DIR="$HOME_DIR/.local/share/applications"
DESKTOP_FILE="$DESKTOP_DIR/skrato.desktop"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing dependency: $1" >&2
    exit 1
  fi
}

# ============================================================
# Uninstall
# ============================================================
if [[ "$ACTION" == "uninstall" ]]; then
  rm -f "$LAUNCHER" "$DESKTOP_FILE" || true
  rm -rf "$TARGET_DIR" || true
  exit 0
fi

# ============================================================
# Install
# ============================================================
require_cmd love
# pkexec is required for actions inside the app; fail early with a clear message.
require_cmd pkexec

mkdir -p "$TARGET_DIR" "$LAUNCHER_DIR" "$DESKTOP_DIR"

# Copy only what LÖVE needs to run (avoid copying Rust build/targets).
# This keeps the install small and reliable.
rsync -a --delete \
  --exclude '/target/' \
  --exclude '/.git/' \
  --exclude '/Cargo.toml' \
  --exclude '/Cargo.lock' \
  --exclude '/src/' \
  --exclude '/installer.sh' \
  ./* "$TARGET_DIR"/

# Launcher
cat >"$LAUNCHER" <<EOF
#!/usr/bin/env bash
exec love "$TARGET_DIR" "\$@"
EOF
chmod 0755 "$LAUNCHER"

# Desktop entry (Icon optional)
ICON_PATH="$TARGET_DIR/icon.png"
if [[ -f "$ICON_PATH" ]]; then
  ICON_LINE="Icon=$ICON_PATH"
else
  ICON_LINE="# Icon missing (optional)"
fi

cat >"$DESKTOP_FILE" <<EOF
[Desktop Entry]
Type=Application
Name=skrato
Comment=Refresh GRUB and Initramfs
Exec=$LAUNCHER
$ICON_LINE
Terminal=false
Categories=System;
EOF
chmod 0644 "$DESKTOP_FILE"

cat <<EOF
Installed skrato to:
- $TARGET_DIR
- $LAUNCHER
- $DESKTOP_FILE

Run: $LAUNCHER
EOF



