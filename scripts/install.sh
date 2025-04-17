#!/bin/bash

set -e

VERSION="v0.1.0"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
fi

BINARY="kindctl-${OS}-${ARCH}"
URL="https://github.com/<your-username>/kindctl/releases/download/${VERSION}/${BINARY}"

echo "Installing kindctl ${VERSION} for ${OS}/${ARCH}..."

curl -L -o kindctl "${URL}"
chmod +x kindctl
sudo mv kindctl /usr/local/bin/

echo "kindctl installed successfully!"