# Building Claude Statusline

This document covers how to build the Lunar Editor TUI and prepare releases.

## Prerequisites

- **Go 1.21+** - [Download](https://go.dev/dl/)
- **Bash** - For testing the statusline script
- **jq** - Required by statusline.sh at runtime

## Project Structure

```
claude-statusline/
├── statusline.sh          # Main statusline script (no build needed)
├── tui/                   # Go TUI application source
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── config/            # Configuration structs and I/O
│   └── ui/                # Bubble Tea views and model
├── lunar-editor-macos     # Pre-built macOS binary
├── lunar-editor-linux     # Pre-built Linux binary
└── lunar-editor.exe       # Pre-built Windows binary
```

## Building the TUI Editor

### Quick Build (Current Platform)

```bash
cd tui
go build -o ../lunar-editor .
```

### Platform-Specific Builds

#### macOS (Apple Silicon & Intel)

```bash
cd tui

# Apple Silicon (M1/M2/M3)
GOOS=darwin GOARCH=arm64 go build -o ../lunar-editor-macos .

# Intel Mac
GOOS=darwin GOARCH=amd64 go build -o ../lunar-editor-macos-intel .
```

#### Linux

```bash
cd tui

# x86_64
GOOS=linux GOARCH=amd64 go build -o ../lunar-editor-linux .

# ARM64
GOOS=linux GOARCH=arm64 go build -o ../lunar-editor-linux-arm64 .
```

#### Windows

```bash
cd tui
GOOS=windows GOARCH=amd64 go build -o ../lunar-editor.exe .
```

### Build All Platforms

```bash
cd tui

# macOS
GOOS=darwin GOARCH=arm64 go build -o ../lunar-editor-macos .

# Linux
GOOS=linux GOARCH=amd64 go build -o ../lunar-editor-linux .

# Windows
GOOS=windows GOARCH=amd64 go build -o ../lunar-editor.exe .

echo "Build complete!"
ls -lh ../lunar-editor*
```

## Build Options

### Optimized Release Build

Strip debug info and reduce binary size:

```bash
cd tui
go build -ldflags="-s -w" -o ../lunar-editor .
```

### With Version Info

```bash
VERSION="1.0.0"
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

go build -ldflags="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" \
    -o ../lunar-editor .
```

## Development

### Run Without Building

```bash
cd tui
go run .
```

### Run Tests

```bash
cd tui
go test ./...
```

### Update Dependencies

```bash
cd tui
go get -u ./...
go mod tidy
```

### Check for Issues

```bash
cd tui
go vet ./...
```

## Testing the Statusline Script

The shell script doesn't need building, but you can test it:

```bash
# Syntax check
bash -n statusline.sh

# Test with sample input
echo '{"model":{"display_name":"Opus 4.5"},"context_window":{"used_percentage":45}}' | ./statusline.sh
```

## Release Checklist

1. Update version numbers if applicable
2. Run tests: `cd tui && go test ./...`
3. Build all platforms (see above)
4. Test each binary on target platform
5. Verify statusline.sh syntax: `bash -n statusline.sh`
6. Update README.md if features changed
7. Commit binaries and tag release

## Troubleshooting

### "go: command not found"

Install Go from https://go.dev/dl/ and ensure it's in your PATH.

### CGO errors on cross-compilation

This project uses pure Go with no CGO dependencies, so cross-compilation should work without issues. If you encounter CGO errors:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../lunar-editor-linux .
```

### Binary too large

Use ldflags to strip debug info:

```bash
go build -ldflags="-s -w" -o ../lunar-editor .
```

For even smaller binaries, consider [UPX](https://upx.github.io/):

```bash
upx --best lunar-editor
```

### Module errors

```bash
cd tui
go mod tidy
go mod download
```
