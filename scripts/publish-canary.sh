#!/bin/bash
set -e

CANARY_NUM=${1:-$(git rev-parse --short HEAD 2>/dev/null || date +%s)}
VERSION="0.0.0-canary.$CANARY_NUM"

echo "Publishing canary version: $VERSION"

./scripts/build.sh "$VERSION"

echo ""
echo "Creating GitHub release v$VERSION..."
gh release create "v$VERSION" dist/* --title "v$VERSION" --notes "Canary release" --prerelease

PUBLISH_DIR="publish"
rm -rf "$PUBLISH_DIR"
mkdir -p "$PUBLISH_DIR/bin" "$PUBLISH_DIR/scripts"

cp npm/scripts/postinstall.js "$PUBLISH_DIR/scripts/"

cat > "$PUBLISH_DIR/package.json" << EOF
{
  "name": "@nanoloop/cli",
  "version": "$VERSION",
  "description": "Nanoloop CLI for uploading source maps",
  "bin": {
    "nanoloop": "bin/nanoloop"
  },
  "scripts": {
    "postinstall": "node scripts/postinstall.js"
  },
  "repository": {
    "type": "git",
    "url": "https://github.com/nanoloop/cli"
  },
  "license": "MIT"
}
EOF

echo "Publishing @nanoloop/cli@$VERSION to npm..."
(cd "$PUBLISH_DIR" && npm publish --tag canary --access public)

echo ""
echo "Done! Install with:"
echo "  npm install @nanoloop/cli@canary"
echo "  npx @nanoloop/cli@canary upload --help"
