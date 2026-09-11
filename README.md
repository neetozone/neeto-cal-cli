# NeetoCal CLI

Manage meetings, bookings and availabilities from the terminal.

<!-- neeto-cli-commons:installation:start -->
## Installation

### macOS / Linux

**Homebrew (recommended on macOS):**

```bash
brew install neetozone/tap/neetocal
```

**Shell script:**

```bash
curl -fsSL https://neeto-downloads.s3.amazonaws.com/cli/NeetoCal/latest/install.sh | sh
```

This verifies the download's SHA-256 checksum against the published `SHA256SUMS`,
then installs to `/usr/local/bin` (may prompt for sudo). Set `NEETOCAL_INSTALL_DIR`
to a directory you own to install without sudo.

### Windows

**PowerShell:**

```powershell
irm https://neeto-downloads.s3.amazonaws.com/cli/NeetoCal/latest/install.ps1 | iex
```

**Command Prompt (CMD):**

```cmd
curl -fsSL https://neeto-downloads.s3.amazonaws.com/cli/NeetoCal/latest/install.cmd -o install.cmd && install.cmd
```

Both verify the download's SHA-256 checksum before installing to
`%LOCALAPPDATA%\Programs\neetocal` and adding it to your user PATH. Set
`NEETOCAL_INSTALL_DIR` to install somewhere else.
<!-- neeto-cli-commons:installation:end -->

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

<!-- neeto-cli-commons:verify-installation:start -->
### Verify installation

```bash
neetocal --help
```
<!-- neeto-cli-commons:verify-installation:end -->

<!-- neeto-cli-commons:ai-coding-assistants:start -->
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
<!-- neeto-cli-commons:ai-coding-assistants:end -->

<!-- neeto-cli-commons:prerequisites:start -->
## Prerequisites (development)

- [Go](https://go.dev/dl/) 1.26.1+
- Access to a NeetoCal organization
<!-- neeto-cli-commons:prerequisites:end -->

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

<!-- neeto-cli-commons:make-targets:start -->
### Make targets

```bash
make build          # Builds ./neetocal
make test           # Run tests
make lint           # golangci-lint
make fmt            # gofmt -w
make vet            # go vet
make check          # fmt + vet + test
make install        # Installs to /usr/local/bin
make clean          # Remove built binary
```
<!-- neeto-cli-commons:make-targets:end -->

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

<!-- neeto-cli-commons:global-flags:start -->
## Global flags

Every command accepts:

| Flag | Description |
|---|---|
| `--subdomain <name>` | Which logged-in subdomain to use (required when multiple are logged in). |
| `--json` | Force JSON envelope output. |
| `--quiet` | Emit raw data only. Action commands print just the identifier; `delete` prints `success`. |
| `--toon` | TOON (Token-Optimized Output Notation) — compact format for LLMs. |
| `--verbose` | Expand every field of a record instead of a table. |
<!-- neeto-cli-commons:global-flags:end -->

<!-- neeto-cli-commons:release:start -->
## Release

Releases are cut by BigBinary's CI pipeline defined in
`.neetoci/release.yml`. Merging a PR with a `major` / `minor` / `patch`
label to `main` triggers the shared release script published by
`neeto-cli-commons`, which bumps and tags VERSION, runs GoReleaser,
uploads artifacts to `s3://neeto-downloads/cli/NeetoCal/`, updates the
Homebrew tap (`neetozone/tap`), and pushes the version bump commit
straight to `main`.
<!-- neeto-cli-commons:release:end -->
