# Installing the EnSync CLI Tool

This guide outlines how to install and configure the CLI tool for multiple platforms using `curl`.

## Prerequisites:
- **curl**: Make sure `curl` is installed on your system.
- **Internet connection**: You need an active internet connection to download the installation package.

---

## General Installation Steps

1. **Download the CLI Tool**

   Choose the correct download command for your platform.

   - **macOS (Darwin)**:
     ```bash
     curl -L -o CLI_1.0.0_darwin_amd64.tar.gz https://github.com/EnSync-engine/CLI/releases/download/v1.0.0/CLI_1.0.0_darwin_amd64.tar.gz
     ```

   - **Linux**:
     ```bash
     curl -L -o CLI_1.0.0_linux_amd64.tar.gz https://github.com/EnSync-engine/CLI/releases/download/v1.0.0/CLI_1.0.0_linux_amd64.tar.gz
     ```

   - **Windows**:
     ```bash
     curl -LO https://github.com/EnSync-engine/CLI/releases/download/CLI_1.0.0_windows_amd64.zip
     ```

2. **Extract the Downloaded Archive**

   - **macOS/Linux**:
     ```bash
     tar -xvf CLI_1.0.0_<platform>_amd64.tar.gz
     ```

   - **Windows**:
     For `.zip` files, right-click the file and select **Extract All**, or use `7zip`:
     ```bash
     7z x CLI_1.0.0_windows_amd64.zip
     ```

3. **Make the Binary Executable (macOS/Linux)**

   After extracting, make the binary executable (only for macOS/Linux):
   ```bash
   chmod +x ensync-cli
   ```

4. **Move the Binary to a Directory in Your PATH**

   To make the `ensync-cli` command accessible globally, move it to a directory included in your `PATH`.

   - **macOS/Linux**:
     ```bash
     sudo mv ensync-cli /usr/local/bin/
     ```

   - **Windows**:
     Move the `ensync-cli.exe` binary to a folder (e.g., `C:\Program Files\CLI\`) and add that folder to your `PATH` environment variable.

     Steps to update PATH on Windows:
     1. Open the Start Menu, search for "Environment Variables," and select "Edit the system environment variables."
     2. In the System Properties window, click on the "Environment Variables" button.
     3. Under "System variables," find and select the `Path` variable, then click "Edit."
     4. Click "New," and add the folder where you placed `ensync-cli.exe` (e.g., `C:\Program Files\CLI\`).
     5. Click "OK" to save your changes and close all dialogs.

5. **Verify the Installation**

   Confirm the tool was installed successfully by checking the version:
   ```bash
   ensync-cli --version
   ```

   This should output the version of the CLI tool that was installed.

---

## Initial Setup

### macOS/Linux

To simplify the configuration, add the environment variables to your shell configuration file:

1. Open your shell configuration file (e.g., `~/.bashrc`, `~/.zshrc`, or `~/.bash_profile`):
   ```bash
   nano ~/.bashrc
   ```

2. Add the following lines:
   ```bash
   export ENSYNC_BASE_URL="http://localhost:8080/api/v1/ensync"
   export ENSYNC_DEBUG=false
   ```

3. Save the file and apply the changes:
   ```bash
   source ~/.bashrc
   ```

### Windows

To set up environment variables on Windows:

1. Open the Start Menu, search for "Environment Variables," and select "Edit the system environment variables."
2. In the System Properties window, click on the "Environment Variables" button.
3. Under "System variables," click "New," and create two variables:
   - **Variable Name**: `ENSYNC_BASE_URL`
     **Variable Value**: `http://localhost:8080/api/v1/ensync`
   - **Variable Name**: `ENSYNC_DEBUG`
     **Variable Value**: `false`
4. Click "OK" to save the variables, then close all dialogs.
5. Open a new Command Prompt or PowerShell window to apply the changes.

---

### Environment Variables

Alternatively, you can configure the CLI using environment variables:

```bash
export ENSYNC_BASE_URL="http://localhost:8080/api/v1/ensync"
export ENSYNC_DEBUG=false
```

### Command-Line Flags

You can also provide the `accessKey` directly via the `--access-key` flag for each command:

```bash
ensync-cli --access-key "your-access-key" <command>
```
---

## Usage

### Event Management

**List events:**
```bash
# List events with pagination
ensync-cli event list --access-key "your-access-key" --page 0 --limit 10 --order DESC --order-by createdAt

# List events with different ordering
ensync-cli event list --access-key "your-access-key" --order ASC --order-by name
```

**Create an event:**
```bash
ensync-cli event create --access-key "your-access-key" --name "test-event" --payload '{"key":"value","another":"data"}'
```

**Update an event:**
```bash
# Update event name
ensync-cli event update --access-key "your-access-key" --id 1 --name "updated/name/name"

# Update event payload
ensync-cli event update --access-key "your-access-key" --id 1 --payload '{"key":"new-value"}'

# Get event payload by Name
ensync-cli event get --access-key "your-access-key" --name "updated/name/name"
```

### Access Key Management

**List access keys:**
```bash
# List all access keys
ensync-cli access-key list --access-key "your-access-key"
```

**Create an access key:**
```bash
# Create access key with permissions
ensync-cli access-key create --access-key "your-access-key" --permissions '{"send": ["event1"], "receive": ["event2"]}'
```

**Manage permissions:**
```bash
# Get current permissions
ensync-cli access-key permissions get --access-key "your-access-key" --key "IeBTeDncBQmDMzJzKblyKfbctvgEKO8L"

# Update permissions
ensync-cli access-key permissions set --access-key "your-access-key" --key "IeBTeDncBQmDMzJzKblyKfbctvgEKO8L" --permissions '{"send": ["event12344"], "receive": ["event23445"]}'
```

### General Options

**Enable debug mode:**
```bash
# Enable debug output
ensync-cli --debug --access-key "your-access-key" event list
```

**Version information:**
```bash
# Display version
ensync-cli version

# Get version in JSON format
ensync-cli version --json
```

---

## Common Flags

- `--access-key`: Access key for API authentication (required for all commands).
- `--page`: Page number for pagination (default: 0).
- `--limit`: Number of items per page (default: 10).
- `--order`: Sort order (ASC/DESC).
- `--order-by`: Field to sort by (name/createdAt).
- `--debug`: Enable debug mode.
- `--config`: Specify custom config file location.

---

### Troubleshooting

- **Permission Issues**:
  - On **macOS/Linux**, you might need `sudo` to move the binary to a directory like `/usr/local/bin`.
  - On **Windows**, ensure you have the necessary permissions to modify the `PATH` and place the binary in the appropriate folder.

- **Command Not Found**:
  If the `ensync-cli` command doesn't work after installation, ensure that the binary is in your `PATH` and check the system's environment variables.

- **Access Key Errors**:
  If you encounter errors related to the `accessKey`, ensure that:
  - The `accessKey` is provided via the `--access-key` flag.
  - The `accessKey` is valid and has the necessary permissions.

---

### Conclusion

The CLI tool can be installed using `curl` and configured with either a configuration file, environment variables, or command-line flags. The `accessKey` is now required for all API interactions, and this guide provides instructions for installation and configuration across various platforms.
