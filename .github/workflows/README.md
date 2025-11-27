# GitHub Actions Workflows

This directory contains the CI/CD workflows for the CloudAMQP CLI.

## Workflows

### CI Workflow (`ci.yml`)

Runs on every push to `main` or `develop` branches, and on pull requests.

**What it does:**
- Tests with multiple Go versions (1.21, 1.22, 1.23)
- Runs format checks, vet, and tests
- Builds the binary with version information
- Uploads test coverage to Codecov
- Verifies the binary works with `--help` and `version` commands

**Version Information:**
The CI build automatically includes:
- Version: From `git describe --tags --always --dirty` (or "dev" as fallback)
- Build Date: Current date in UTC (YYYY-MM-DD)
- Git Commit: Short commit hash

### Release Workflow (`release.yml`)

Triggers when a tag starting with `v` is pushed (e.g., `v1.0.0`), or can be run manually.

**What it does:**
1. **Build Job**: Builds cross-platform binaries for:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64, arm64)

2. **Package Job**: Creates release archives:
   - `.tar.gz` for Linux/macOS binaries
   - `.zip` for Windows binaries
   - Creates GitHub Release with all artifacts

**Version Information:**
Release builds automatically include:
- Version: From the git tag (e.g., `v1.0.0`)
- Build Date: Current date in UTC (YYYY-MM-DD)
- Git Commit: Short commit hash

The version info is embedded using Go's `-ldflags`:
```bash
-ldflags="-w -s -X cloudamqp-cli/cmd.Version=$VERSION -X cloudamqp-cli/cmd.BuildDate=$BUILD_DATE -X cloudamqp-cli/cmd.GitCommit=$GIT_COMMIT"
```

## Creating a Release

To create a new release:

1. **Create and push a tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The release workflow will automatically:
   - Build binaries for all platforms
   - Create archives
   - Create a GitHub release with auto-generated release notes
   - Upload all artifacts

3. The released binaries will show proper version info:
   ```bash
   $ cloudamqp version
   cloudamqp version v1.0.0 (2025-11-26)
   https://github.com/cloudamqp/cli/releases/tag/v1.0.0
   ```

## Manual Workflow Trigger

The release workflow can also be triggered manually from the GitHub Actions UI:
1. Go to Actions tab
2. Select "Build Release" workflow
3. Click "Run workflow"
4. Select branch and click "Run workflow"

## Verification

Both workflows include verification steps:
- **CI**: Tests the built binary with `--help` and `version` commands
- **Release**: Attempts to run `version` command on non-Windows binaries (may skip for cross-compiled binaries)

## Caching

Both workflows use Go module caching to speed up builds:
- Cache key includes Go version and `go.sum` hash
- Cached paths: `~/.cache/go-build` and `~/go/pkg/mod`

## Build Flags

**CI builds:**
- Standard build flags
- Version information included

**Release builds:**
- `-a`: Force rebuilding of packages
- `-installsuffix cgo`: Use suffix for cgo-enabled builds
- `-ldflags="-w -s"`: Strip debug info and symbol tables (smaller binaries)
- Version information via `-X` flags
- `CGO_ENABLED=0`: Disable cgo for static binaries

## Artifacts

**CI Workflow:**
- Uploads test coverage to Codecov

**Release Workflow:**
- Individual binaries (30 days retention)
- Release archives (90 days retention)
- GitHub Release assets (permanent)

## Troubleshooting

If version information is not showing correctly in releases:

1. Check that tags are pushed: `git push origin --tags`
2. Verify tag format: Must start with `v` (e.g., `v1.0.0`)
3. Check workflow logs for version extraction step
4. Ensure no local modifications (avoid `-dirty` in version)

## Local Testing

To test version embedding locally:

```bash
# Using Make (recommended)
make build

# Manual
VERSION=v1.0.0
BUILD_DATE=$(date -u +"%Y-%m-%d")
GIT_COMMIT=$(git rev-parse --short HEAD)

go build -ldflags "\
  -X cloudamqp-cli/cmd.Version=$VERSION \
  -X cloudamqp-cli/cmd.BuildDate=$BUILD_DATE \
  -X cloudamqp-cli/cmd.GitCommit=$GIT_COMMIT" \
  -o cloudamqp .
```
