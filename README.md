# NeetoCal CLI

Command-line interface for [NeetoCal](https://neetocal.com). Manage meetings, bookings, availabilities, and more from the terminal.

## Installation

### macOS

**Homebrew (recommended):**

```bash
brew install neetozone/neetocal/neetocal
```

To update later:

```bash
brew upgrade neetocal
```

**Shell script:**

```bash
curl -fsSL https://raw.githubusercontent.com/neetozone/neeto-cal-cli/main/install.sh | sh
```

This downloads the latest release, extracts it, and installs to `/usr/local/bin` (may prompt for sudo).

### Linux

**Shell script:**

```bash
curl -fsSL https://raw.githubusercontent.com/neetozone/neeto-cal-cli/main/install.sh | sh
```

This downloads the latest release for your architecture (amd64 or arm64), extracts it, and installs to `/usr/local/bin` (may prompt for sudo).

### Windows

**PowerShell (recommended):**

```powershell
irm https://raw.githubusercontent.com/neetozone/neeto-cal-cli/main/install.ps1 | iex
```

**Command Prompt (CMD):**

```cmd
curl -fsSL https://raw.githubusercontent.com/neetozone/neeto-cal-cli/main/install.cmd -o install.cmd && install.cmd
```

Both methods download the latest release, extract it to `%LOCALAPPDATA%\Programs\neetocal`, and add it to your user PATH.

### Manual download

Download the latest release archive for your platform from the [Releases](https://github.com/neetozone/neeto-cal-cli/releases/latest) page, extract it, and place the `neetocal` binary somewhere on your PATH.

| Platform       | Archive                                    |
|----------------|--------------------------------------------|
| macOS (Apple Silicon) | `neeto-cal-cli_*_darwin_arm64.tar.gz` |
| macOS (Intel)  | `neeto-cal-cli_*_darwin_amd64.tar.gz`      |
| Linux (x86_64) | `neeto-cal-cli_*_linux_amd64.tar.gz`       |
| Linux (ARM)    | `neeto-cal-cli_*_linux_arm64.tar.gz`       |
| Windows (x86_64) | `neeto-cal-cli_*_windows_amd64.zip`     |
| Windows (ARM)  | `neeto-cal-cli_*_windows_arm64.zip`        |

### Verify installation

After installing, restart your terminal and run:

```bash
neetocal --help
```

## Prerequisites (development)

- [Go](https://go.dev/dl/) 1.22+
- Access to a NeetoCal organization

## Development

### Setup

```bash
git clone https://github.com/neetozone/neeto-cal-cli.git
cd neeto-cal-cli
go mod download
```

### Build and run

```bash
make build          # Builds ./neetocal binary
./neetocal --help

make install        # Builds and copies to /usr/local/bin
```

### Other make targets

```bash
make test           # Run tests
make lint           # Run golangci-lint
make fmt            # Format code
make vet            # Run go vet (catches bugs the compiler misses, like bad format strings or unreachable code)
make check          # fmt + vet + test
make clean          # Remove built binary
```

### Pointing to a local or staging server

By default the CLI targets `https://{subdomain}.neetocal.com`. Set `NEETOCAL_BASE_URL` to point to a different server:

```bash
# Local development (neeto-cal-web runs on port 8980)
export NEETOCAL_BASE_URL=http://spinkart.lvh.me:8980
neetocal login --subdomain spinkart

# Staging
export NEETOCAL_BASE_URL=https://spinkart.neetocal-staging.neetohq.com
neetocal login --subdomain spinkart
```

### Local development workflow

1. Start neeto-cal-web locally (runs on port 8980 by default).

2. Build and login:
   ```bash
   make build
   export NEETOCAL_BASE_URL=http://spinkart.lvh.me:8980
   ./neetocal login --subdomain spinkart
   ```

3. Verify connectivity:
   ```bash
   ./neetocal doctor
   ```

4. Test commands:
   ```bash
   ./neetocal meetings list --json
   ./neetocal bookings list --type upcoming --json
   ```

## Release

Releases are built with [GoReleaser](https://goreleaser.com/).

### Install GoReleaser

```bash
brew install goreleaser     # macOS
```

### Create a release

1. Tag the version:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. Build release artifacts:
   ```bash
   goreleaser release --clean
   ```

3. For a local test build (no publish):
   ```bash
   goreleaser release --snapshot --clean
   ```

GoReleaser produces archives for Linux, macOS, and Windows (amd64 + arm64). It also publishes to the Homebrew tap at `neetozone/homebrew-neetocal`. Version, commit hash, and build date are injected via ldflags at build time.
