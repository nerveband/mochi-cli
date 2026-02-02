#!/bin/bash

set -e

REPO="nerveband/mochi-cli"
BINARY_NAME="mochi"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="x86_64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Installing mochi-cli for $OS/$ARCH..."

# Get latest release
LATEST_URL=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep "browser_download_url.*${OS}_${ARCH}.tar.gz" | cut -d '"' -f 4)

if [ -z "$LATEST_URL" ]; then
    echo "Could not find release for $OS/$ARCH"
    exit 1
fi

# Download and extract
TEMP_DIR=$(mktemp -d)
curl -L -o "$TEMP_DIR/mochi.tar.gz" "$LATEST_URL"
tar -xzf "$TEMP_DIR/mochi.tar.gz" -C "$TEMP_DIR"

# Install
if [ -w "/usr/local/bin" ]; then
    mv "$TEMP_DIR/$BINARY_NAME" "/usr/local/bin/"
else
    echo "Installing to /usr/local/bin requires sudo..."
    sudo mv "$TEMP_DIR/$BINARY_NAME" "/usr/local/bin/"
fi

# Cleanup
rm -rf "$TEMP_DIR"

echo "mochi-cli installed successfully!"
echo "Run 'mochi --help' to get started"
