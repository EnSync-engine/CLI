# EnSync CLI Tool

A command-line interface for interacting with the EnSync configuration management service.

## Installation

### 1. From Source
```bash
# Clone the repository
git clone https://github.com/EnSync-engine/cli
cd ensync-cli

# Build the binary
go build -o ensync ./main.go

# Verify installation
./ensync version
```

### 2. Download Binary (No Go required)

**macOS & Linux**
1. Download the `tar.gz` archive for your OS/Arch from the [Releases](https://github.com/EnSync-engine/CLI/releases) page.
2. Extract and install:
   ```bash
   tar -xzf CLI_*.tar.gz
   chmod +x ensync
   sudo mv ensync /usr/local/bin/
   ```
3. Verify: `ensync version`

**Windows**
1. Download the `.zip` archive for your OS/Arch from the [Releases](https://github.com/EnSync-engine/CLI/releases) page.
2. Extract the zip file.
3. Open PowerShell as Administrator and move the binary (optional but recommended):
   ```powershell
   mkdir C:\ensync
   move ensync.exe C:\ensync\
   # Add to PATH (current session)
   $env:Path += ";C:\ensync"
   # Add to PATH (permanent)
   [Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\ensync", [EnvironmentVariableTarget]::User)
   ```
4. Verify: `ensync version`

### 3. Using Go Install
```bash
go install github.com/EnSync-engine/cli@latest
```

## Configuration

The CLI can be configured using a config file or environment variables.

### Configuration File
Create a config file at:
*   **macOS / Linux**: `~/.ensync/config.yaml`
*   **Windows**: `%UserProfile%\.ensync\config.yaml`

```yaml
base_url: "https://access.gms.ensync.cloud/api/v1/ensync"
debug: false
```

### Environment Variables

**macOS & Linux**
```bash
export ENSYNC_ACCESS_KEY="your-access-key"
export ENSYNC_BASE_URL="https://access.gms.ensync.cloud/api/v1/ensync"
export ENSYNC_DEBUG=false
```

**Windows (PowerShell)**
```powershell
$env:ENSYNC_ACCESS_KEY="your-access-key"
$env:ENSYNC_BASE_URL="https://access.gms.ensync.cloud/api/v1/ensync"
$env:ENSYNC_DEBUG="false"
```

## Usage

All commands require an access key, either via `--access-key` flag or `ENSYNC_ACCESS_KEY` environment variable.

### Event Management

```bash
# List events (default limit: 10)
ensync event list 
ensync event list --limit 20 --order ASC

# Get event by name
ensync event get "gms/urbanhero/stripe"

# Create event
ensync event create --name "my-event" --payload '{"key":"value"}'

# Update event
ensync event update --id "event-uuid" --name "new-name"
ensync event update --id "event-uuid" --payload '{"new":"data"}'
```

### Access Key Management

```bash
# List access keys
ensync access-key list

# Get access key details
ensync access-key get "key-uuid"

# Create access key
ensync access-key create --name "my-service" --type SERVICE --permissions '{"send":["event1"],"receive":["event2"]}'

# Delete access key
ensync access-key delete "key-uuid"

# Rotate service key pair
ensync access-key rotate "access-key-string"

# Manage Permissions
ensync access-key permissions get "access-key-string"
ensync access-key permissions set "access-key-string" --permissions '{"send":["*"],"receive":["*"]}'
```

### Workspace Management

```bash
# List workspaces (hierarchical tree view)
ensync workspace list

# Create workspace
ensync workspace create --name "my-workspace"
```

### General Options

```bash
# Enable debug output (shows HTTP requests/responses)
ensync --debug event list

# Version info
ensync version --json
```

## Common Flags

- `--limit`: Number of items per page (default: 10)
- `--page`: Page index (default: 0)
- `--order`: Sort order (`ASC` or `DESC`)
- `--order-by`: Field to sort by (e.g., `createdAt`)
- `--access-key`: Authentication key
- `--debug`: Enable verbose logging

## Development

```bash
# Run tests
go test ./...

# Run integration tests
go test ./test/integration/... -v

# Build
go build -o ensync ./main.go
```

## License

MIT License
