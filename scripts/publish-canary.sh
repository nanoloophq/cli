#!/bin/bash
set -e

echo "Platform packages are published via GitHub Actions."
echo "Push to main to trigger a release."
echo ""
echo "For local testing, build the binary:"
echo "  go build -o nanoloop ."
echo "  ./nanoloop upload --help"
