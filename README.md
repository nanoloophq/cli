# @nanoloop/cli

CLI for uploading source maps to [nanoloop](https://nanoloop.app).

## Installation

```bash
npm install -D @nanoloop/cli

# Or install canary for testing
npm install -D @nanoloop/cli@canary
```

## Usage

```bash
# Upload source maps after building
npm run build && npx nanoloop upload --token $NANOLOOP_TOKEN --app $APP_ID --dist ./dist
```

### Options

| Flag | Env Variable | Description |
|------|--------------|-------------|
| `--token` | `NANOLOOP_TOKEN` | API token (required) |
| `--app` | `NANOLOOP_APP_ID` | App ID (required) |
| `--dist` | - | Directory containing source maps (default: `./dist`) |
| `--release` | - | Release version (default: git commit hash) |
| `--url-prefix` | - | URL prefix for source files |
| `--dry-run` | - | List files without uploading |

### CI/CD Example

```yaml
# GitHub Actions
- name: Upload source maps
  run: npx nanoloop upload --dist ./dist
  env:
    NANOLOOP_TOKEN: ${{ secrets.NANOLOOP_TOKEN }}
    NANOLOOP_APP_ID: ${{ vars.NANOLOOP_APP_ID }}
```

## Getting an API Token

1. Go to [nanoloop.app/settings](https://nanoloop.app/settings)
2. Create a new API token
3. Copy the token (it's only shown once)

## Development

```bash
go build -o nanoloop .
./nanoloop upload --help
```

## Publishing

```bash
# Publish canary (for testing)
./scripts/publish-canary.sh

# Build only (no publish)
./scripts/build.sh
```

## License

MIT
