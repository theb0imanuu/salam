#!/bin/bash

set -e

REPO="theb0imanuu/salam"
BINARY_NAME="salam"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "🕊️  Installing Salam (Server Health Monitor)..."

# Check for existing installation
if command -v $BINARY_NAME &> /dev/null; then
    echo "⚠️  Salam is already installed. Updating..."
fi

# Download latest release
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
fi

echo "📦 Downloading from ${DOWNLOAD_URL}..."

TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

if command -v curl &> /dev/null; then
    curl -fsSL -o "$TMP_DIR/$BINARY_NAME" "$DOWNLOAD_URL"
elif command -v wget &> /dev/null; then
    wget -q -O "$TMP_DIR/$BINARY_NAME" "$DOWNLOAD_URL"
else
    echo "❌ curl or wget is required"
    exit 1
fi

chmod +x "$TMP_DIR/$BINARY_NAME"

# Install
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    echo "✓ Installed to $INSTALL_DIR/$BINARY_NAME"
else
    echo "🔐 Sudo required for installation to $INSTALL_DIR"
    sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    echo "✓ Installed to $INSTALL_DIR/$BINARY_NAME"
fi

# Verify installation
if command -v $BINARY_NAME &> /dev/null; then
    echo ""
    echo "🎉 Installation complete!"
    echo "Run 'salam --help' to get started"
    echo ""
    echo "Quick start:"
    echo "  salam check          # Run one-time health check"
    echo "  salam watch          # Continuous monitoring"
    echo "  salam config         # Generate config file"
else
    echo "❌ Installation failed. Please check your PATH."
    exit 1
fi