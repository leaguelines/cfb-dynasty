# cfb-dynasty-gui

Desktop explorer for EA Sports College Football dynasty saves. Built with [Wails v2](https://wails.io).

## Run

From the repository root:

```bash
go run -tags production ./cmd/cfb-dynasty-gui --schema-dir /path/to/schemas
```

Flags:

| Flag | Description |
|------|-------------|
| `--schema-dir` | Folder containing `C27_*.gz` schema bundles |
| `--save` | Optional save file to open immediately |

If `--schema-dir` is omitted, the app checks `SCHEMA_DIR`, then the last folder you picked (saved under the OS user config dir), then `./data/schemas` / `./schemas` next to the working directory or binary.

## Build

```bash
# Single platform
go build -tags production -o cfb-dynasty-gui ./cmd/cfb-dynasty-gui

# Wails packaging (icons, platform bundles)
wails build
```

CI builds Windows, macOS, and Linux packages on every `main` push (nightly) and on `v*` tags (versioned releases). See [`.github/workflows/release.yml`](../../.github/workflows/release.yml).

## Notes

- Schema bundles are game-derived and must not be redistributed with the app.
- CSV/JSON export uses the same collection flattening as cfb-dynasty-web.
- Optional team/conference logo PNGs can be placed under `frontend/dist/static/teams` and `.../conferences` (gitignored).
