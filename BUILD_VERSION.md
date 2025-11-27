# Building with Version Information

The `cloudamqp` CLI supports version information that can be set at build time using Go's `-ldflags`.

## Version Variables

The following variables are available in `cmd/version.go`:
- `Version` - The version number (e.g., "1.0.0")
- `BuildDate` - The build date (e.g., "2025-11-25")
- `GitCommit` - The git commit hash (optional)

## Build Examples

### Using Make (Recommended)

The Makefile automatically extracts version information from git:

```bash
# Build with automatic git version
make build
```

Output:
```
$ ./cloudamqp version
cloudamqp version v0.1.0-8-g285bd2b (2025-11-26)
https://github.com/cloudamqp/cli/releases/tag/v0.1.0-8-g285bd2b
```

Show version information before building:
```bash
make version-info
```

Override version manually:
```bash
make build VERSION=1.0.0 BUILD_DATE=2025-11-25
```

### Manual Build (without Make)

#### Development Build
```bash
go build -o cloudamqp-cli
```
Output:
```
$ cloudamqp-cli version
cloudamqp version dev (development build)
```

#### Release Build
```bash
go build -ldflags "\
  -X cloudamqp-cli/cmd.Version=1.0.0 \
  -X cloudamqp-cli/cmd.BuildDate=2025-11-25 \
  -X cloudamqp-cli/cmd.GitCommit=abc123" \
  -o cloudamqp-cli
```
Output:
```
$ cloudamqp-cli version
cloudamqp version 1.0.0 (2025-11-25)
https://github.com/cloudamqp/cli/releases/tag/v1.0.0
```

### Using Git Information
```bash
VERSION=$(git describe --tags --always --dirty)
BUILD_DATE=$(date -u +"%Y-%m-%d")
GIT_COMMIT=$(git rev-parse --short HEAD)

go build -ldflags "\
  -X cloudamqp-cli/cmd.Version=${VERSION} \
  -X cloudamqp-cli/cmd.BuildDate=${BUILD_DATE} \
  -X cloudamqp-cli/cmd.GitCommit=${GIT_COMMIT}" \
  -o cloudamqp-cli
```

## Usage

The version can be displayed in two ways:

```bash
# Using the version command
cloudamqp version

# Using the --version or -v flag
cloudamqp --version
cloudamqp -v
```

Both produce identical output, similar to GitHub CLI (`gh`).

## CI/CD Integration

### GitHub Actions Example
```yaml
- name: Build
  run: |
    VERSION=${GITHUB_REF#refs/tags/}
    BUILD_DATE=$(date -u +"%Y-%m-%d")
    GIT_COMMIT=${GITHUB_SHA::7}

    go build -ldflags "\
      -X cloudamqp-cli/cmd.Version=${VERSION} \
      -X cloudamqp-cli/cmd.BuildDate=${BUILD_DATE} \
      -X cloudamqp-cli/cmd.GitCommit=${GIT_COMMIT}" \
      -o cloudamqp-cli
```

### GoReleaser Example
```yaml
builds:
  - ldflags:
      - -X cloudamqp-cli/cmd.Version={{.Version}}
      - -X cloudamqp-cli/cmd.BuildDate={{.Date}}
      - -X cloudamqp-cli/cmd.GitCommit={{.ShortCommit}}
```
