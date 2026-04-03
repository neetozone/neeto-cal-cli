#!/bin/sh
set -e

GITHUB_OWNER="neetozone"
GITHUB_REPO="neeto-cal-cli"
BINARY_NAME="neetocal"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    darwin) ;;
    linux) ;;
    *) echo "Error: unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Error: unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version
echo "Fetching latest version..."
VERSION=$(curl -fsSL "https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Error: could not determine latest version."
    exit 1
fi

VERSION_NUM="${VERSION#v}"

# Download
FILENAME="${GITHUB_REPO}_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}/releases/download/${VERSION}/${FILENAME}"

echo "Downloading ${BINARY_NAME} ${VERSION} for ${OS}/${ARCH}..."
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

curl -fsSL "$URL" -o "${TMP_DIR}/${FILENAME}"

# Extract
tar -xzf "${TMP_DIR}/${FILENAME}" -C "${TMP_DIR}"

# Install
if [ -w "$INSTALL_DIR" ]; then
    mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo "Installing to ${INSTALL_DIR} (requires sudo)..."
    sudo mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

echo ""
echo "${BINARY_NAME} ${VERSION} installed successfully."
echo "Run '${BINARY_NAME} --help' to get started."
