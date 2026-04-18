# @nanoloop/cli

CLI for uploading source maps to [Nanoloop](https://nanoloop.app).

## Installation

```bash
npm install -D @nanoloop/cli
```

## Usage

```bash
npm run build && npx @nanoloop/cli upload --token $NANOLOOP_TOKEN --app $APP_ID --dist ./dist
```

### Options

| Flag | Env Variable | Description |
|------|--------------|-------------|
| `--token` | `NANOLOOP_TOKEN` | API token (required) |
| `--app` | `NANOLOOP_APP_ID` | App ID (required) |
| `--dist` | | Directory containing source maps |
| `--release` | | Release version (default: git commit hash) |
| `--url-prefix` | | URL prefix for source files |
| `--dry-run` | | List files without uploading |

### CI/CD Example

```yaml
# GitHub Actions
- name: Upload source maps
  run: npx @nanoloop/cli upload --dist ./dist
  env:
    NANOLOOP_TOKEN: ${{ secrets.NANOLOOP_TOKEN }}
    NANOLOOP_APP_ID: ${{ vars.NANOLOOP_APP_ID }}
```

## Getting an API Token

1. Go to [nanoloop.app/settings](https://nanoloop.app/settings)
2. Create a new API token
3. Copy the token (it is only shown once)

## Supported Platforms

- macOS (Apple Silicon)
- macOS (Intel)
- Linux (x64)
- Linux (arm64)
- Windows (x64)

## License

MIT
