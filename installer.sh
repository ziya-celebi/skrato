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
  *) echo "Unknown action: $ACTION" >&2; usage; exit 1 ;;
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
  rm -f "$LAUNCHER" "$DESKTOP_FILE" 2>/dev/null || true
  rm -rf "$TARGET_DIR" 2>/dev/null || true
  echo "Uninstalled skrato (removed $TARGET_DIR, $LAUNCHER, $DESKTOP_FILE)."
  exit 0
fi

# ============================================================
# Install
# ============================================================
require_cmd pkexec

# We install the prebuilt Rust binary (release mode) into ~/.skrato.
# installer.sh can build it if missing, so rsync/love are not required at runtime.

# Optional build step: if the release binary isn't present, build it.


mkdir -p "$TARGET_DIR" "$LAUNCHER_DIR" "$DESKTOP_DIR"

BIN_NAME="skrato"
RELEASE_BIN="$(pwd)/target/release/$BIN_NAME"

# Build release binary if missing
if [[ ! -f "$RELEASE_BIN" ]]; then
  echo "Release binary not found at $RELEASE_BIN. Building with cargo..."
  cargo build --release
fi

if [[ ! -f "$RELEASE_BIN" ]]; then
  echo "Expected release binary not found: $RELEASE_BIN" >&2
  exit 1
fi

# Install binary
cp -f "$RELEASE_BIN" "$TARGET_DIR/$BIN_NAME"
chmod 0755 "$TARGET_DIR/$BIN_NAME"

# Launcher: run the installed Rust binary.
cat >"$LAUNCHER" <<EOF
#!/usr/bin/env bash
exec "$TARGET_DIR/$BIN_NAME" "\$@"
EOF
chmod 0755 "$LAUNCHER"


# Desktop entry (icon optional)
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
Comment=Maintain bootloader and initramfs
Exec=$LAUNCHER
$ICON_LINE
Terminal=false
Categories=System;
EOF

chmod 0644 "$DESKTOP_FILE"

echo "Installed skrato to:"
echo "- $TARGET_DIR"
echo "- $LAUNCHER"
echo "- $DESKTOP_FILE"
echo "Run: $LAUNCHER"

