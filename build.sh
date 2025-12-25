#!/bin/bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

OUT_DIR="$ROOT_DIR/dist"
mkdir -p "$OUT_DIR"

APP_NAME="pars"

have_garble() {
  command -v garble >/dev/null 2>&1
}

echo "==> Output dir: $OUT_DIR"
if have_garble; then
  echo "==> garble: yes"
else
  echo "==> garble: no"
fi

echo "==> Building Linux (amd64)..."
if have_garble; then
  CGO_ENABLED=1 garble -literals -tiny -seed=random build -trimpath -ldflags="-s -w" -o "$OUT_DIR/${APP_NAME}-linux-amd64"
else
  CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o "$OUT_DIR/${APP_NAME}-linux-amd64"
fi

echo "==> Building Windows (amd64)..."
if ! command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1; then
  echo "ERROR: mingw-w64 not installed (need x86_64-w64-mingw32-gcc)"
  echo "Install:"
  echo "  sudo apt-get install mingw-w64"
  exit 1
fi

export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++

if have_garble; then
  garble -literals -tiny -seed=random build -trimpath -ldflags="-s -w" -o "$OUT_DIR/${APP_NAME}-windows-amd64.exe"
else
  go build -trimpath -ldflags="-s -w" -o "$OUT_DIR/${APP_NAME}-windows-amd64.exe"
fi

unset GOOS GOARCH CGO_ENABLED CC CXX

echo "✅ Done:"
echo " - $OUT_DIR/${APP_NAME}-linux-amd64"
echo " - $OUT_DIR/${APP_NAME}-windows-amd64.exe"

