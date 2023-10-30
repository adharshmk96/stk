#!/bin/bash

# Set variables
REPO_URL="https://github.com/adharshmk96/stk/releases/download"  # replace with your repository URL
VERSION="1.1.0"
BINARY_NAME="stk"

# Detect OS
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS='linux';;
    Darwin*)    OS='darwin';;
    *)          echo "Unknown or unsupported OS"; exit 1;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64)     ARCH='amd64';;
    amd64)      ARCH='amd64';;
    arm64)      ARCH='arm64';;
    aarch64)    ARCH='arm64';;
    *)          echo "Unsupported architecture"; exit 1;;
esac

# Construct file name
FILE="${BINARY_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"

# Download, check for errors, extract, and then set permissions
if curl -sL -f "${REPO_URL}/v${VERSION}/${FILE}" | sudo tar -xz -C /usr/local/bin -f - "${BINARY_NAME}"; then
    sudo chmod +x /usr/local/bin/${BINARY_NAME}
    echo "Installation completed."
else
    echo "Error: Installation failed."
    exit 1
fi
