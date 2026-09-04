# NeetoCal CLI

Command-line interface for [NeetoCal](https://neetocal.com). Manage meetings, bookings, availabilities, and more from the terminal.

## Installation

### macOS

**Homebrew (recommended):**

```bash
brew trust neetozone/tap
brew install neetozone/tap/neetocal
```

To update later:

```bash
brew upgrade neetocal
```

**Shell script:**

```bash
curl -fsSL https://neetocal.com/cli/install.sh | sh
```

This downloads the latest release, verifies its SHA-256 checksum against the published `SHA256SUMS`, extracts it, and installs to `/usr/local/bin` (may prompt for sudo). Set `NEETOCAL_INSTALL_DIR` to a directory you own to install without sudo.

### Linux

**Shell script:**

```bash
curl -fsSL https://neetocal.com/cli/install.sh | sh
```

This downloads the latest release for your architecture (amd64 or arm64), verifies its SHA-256 checksum against the published `SHA256SUMS`, extracts it, and installs to `/usr/local/bin` (may prompt for sudo). Set `NEETOCAL_INSTALL_DIR` to a directory you own to install without sudo.

### Windows

**PowerShell (recommended):**

```powershell
irm https://neetocal.com/cli/install.ps1 | iex
```

**Command Prompt (CMD):**

```cmd
curl -fsSL https://neetocal.com/cli/install.cmd -o install.cmd && install.cmd
```

Both methods download the latest release, verify its SHA-256 checksum against the published `SHA256SUMS`, extract it to `%LOCALAPPDATA%\Programs\neetocal`, and add it to your user PATH. Set `NEETOCAL_INSTALL_DIR` to install somewhere else.

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

## AI coding assistants

```bash
neetocal setup claude      # Register plugin with Claude Code
neetocal setup cursor      # Write .cursor/rules/neetocal.mdc
neetocal setup windsurf    # Write .windsurf/rules/neetocal.md
neetocal setup copilot     # Add a NeetoCal section to .github/copilot-instructions.md
neetocal setup gemini      # Add a NeetoCal section to GEMINI.md
neetocal setup codex       # Add a NeetoCal section to AGENTS.md
```

Every command except `setup claude` writes into the current project directory, so
run these commands from the root of the project the assistant works in. Re-run
them after every upgrade: `setup cursor` and `setup windsurf` overwrite their rule
file, while `setup copilot`, `setup gemini` and `setup codex` keep the existing
content of their file and replace only the NeetoCal section instead of adding a
duplicate.

## Prerequisites (development)

- [Go](https://go.dev/dl/) 1.22+
- Access to a NeetoCal organization

## Development

### Setup

```bash
git clone https://github.com/neetozone/neeto-cal-cli.git
cd neeto-cal-cli
bin/setup
```

This installs Go dependencies, golangci-lint, configures git hooks, and builds the binary.

### Build and run

```bash
# Local installation
make build          # Builds ./neetocal binary
./neetocal --help

# Global installation
make install        # Builds and copies to /usr/local/bin
neetocal help
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
export NEETOCAL_BASE_URL=https://spinkart.neetocal.net
neetocal login --subdomain spinkart
```

### Local development workflow

1. Start neeto-cal-web locally (runs on port 8980 by default).

2. Build and install globally:
   ```bash
   make install    # Builds and copies to /usr/local/bin (may need sudo)
   ```

3. Login:
   ```bash
   export NEETOCAL_BASE_URL=http://spinkart.lvh.me:8980
   neetocal login --subdomain spinkart
   ```

4. Verify connectivity:
   ```bash
   neetocal doctor
   ```

5. Test commands:
   ```bash
   neetocal meetings list
   neetocal bookings list --type upcoming --json
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

GoReleaser produces archives for Linux, macOS, and Windows (amd64 + arm64). It also publishes to the Homebrew tap at `neetozone/homebrew-tap`. Version, commit hash, and build date are injected via ldflags at build time.
