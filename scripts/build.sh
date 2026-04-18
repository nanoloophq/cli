#!/bin/bash
set -e

VERSION=${1:-"0.0.0-canary.$(git rev-parse --short HEAD 2>/dev/null || date +%s)"}
OUT_DIR="dist"

rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

platforms=(
    "darwin/arm64"
    "darwin/amd64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    GOOS="${platform%/*}"
    GOARCH="${platform#*/}"

    output="$OUT_DIR/nanoloop-$GOOS-$GOARCH"
    if [ "$GOOS" = "windows" ]; then
        output="$output.exe"
    fi

    echo "Building $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w -X main.version=$VERSION" -o "$output" .
done

echo ""
echo "Version: $VERSION"
echo "Binaries in $OUT_DIR/"
ls -la "$OUT_DIR"
