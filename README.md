# YERD - A powerful, developer-friendly tool for managing PHP versions

<div align="center">

![YERD Logo](.meta/yerd_logo.jpg)

**A powerful, developer-friendly tool for managing PHP versions and local development environments with ease**

https://github.com/LumoSolutions/yerd

[![Release](https://img.shields.io/github/v/release/LumoSolutions/yerd)](https://github.com/LumoSolutions/yerd/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Linux-blue.svg)](https://kernel.org)

</div>

---

## 🎯 Purpose

YERD is designed to solve two key challenges:

1. **Production Server Management**: Install and manage multiple PHP versions on Linux servers for different applications and projects
2. **Development Environment Control**: Easily switch between PHP versions on development machines for testing and compatibility

Built by **LumoSolutions**, YERD compiles PHP from official source code, ensuring complete control and reliability without depending on external package repositories.

**🚀 Coming Soon**: NGINX configuration management and site linking/unlinking functionality for complete web server environment control.

**Cross-platform Linux compatibility** with automatic dependency management for all major distributions.

## ✨ Features

### Core Features
- 🚀 **Install and manage multiple PHP versions simultaneously**
- ⚡ **Switch PHP CLI versions instantly with simple commands**
- 🪶 **Lightweight and fast - no unnecessary overhead**
- 👨‍💻 **Developer friendly**

### Installation & Management
- 🏗️ **Build from official PHP source** - Direct from php.net for maximum reliability
- 🔄 **Dynamic version fetching** - Automatically fetches latest PHP versions from php.net
- 📊 **Update checking** - Shows available updates for installed PHP versions
- 🧹 **Clean removal** of PHP versions with complete directory cleanup
- 🔄 **Zero-downtime updates** - Update PHP versions safely with automatic rollback
- 🌐 **Isolated installations** - Each version in separate `/opt/yerd/php/phpX.X/` directories
- 🛠️ **Multi-distro dependency management** - Automatically installs build tools on all major Linux distributions
- 📦 **Rich PHP extensions** - Built-in support for MySQL, GD, cURL, OpenSSL, and more
- 🎯 **Multi-core compilation** - Uses all CPU cores for faster builds
- ⚡ **Intelligent caching** - Hourly refresh of version data to minimize API requests
- 🌐 **Real-time updates** - Always installs the latest patch versions automatically
- 🤖 **Automation support** - `-y` flag for unattended updates in scripts and CI/CD
- 🛡️ **Enterprise-grade reliability** - Install-first, cleanup-after update strategy

### Safety & Security
- 🛡️ **System PHP conflict protection** - Never overwrites existing installations
- 🔒 **Safe removal** - Confirms before removing current CLI versions
- 🔐 **Secure privilege handling** - Build processes run as user, minimal root usage
- 🔍 **Intelligent binary detection** - Finds PHP binaries in installation directories

### User Experience
- 🎨 **Beautiful CLI interface** - Colored output, ASCII art, and loading spinners
- ⏳ **Visual build progress** - Loading indicators for download, configure, compile, and install phases
- 📊 **System status monitoring** - Comprehensive conflict detection and build environment checks
- 🩺 **Advanced diagnostics** - Troubleshoot installation issues with `yerd php doctor`
- 📝 **Comprehensive logging** - Detailed build logs with automatic cleanup on success

## 🚀 Coming Soon

- 🌐 **NGINX Configuration Management** - Automatic NGINX setup and configuration
- 🔗 **Site Linking/Unlinking** - Easy website deployment and management
- 📂 **Virtual Host Management** - Create and manage multiple sites per PHP version
- 🔧 **Development Environment Presets** - Quick setup for common development scenarios

## 📋 Requirements

- **Operating System**: Linux (any distribution)
- **Build Tools**: Automatically detected and installed per distribution
- **Permissions**: `sudo` access for system-wide operations
- **Development Libraries**: Automatically installed during first PHP build
- **For Building YERD**: Go 1.21+ (optional - only for building YERD from source)

## 🐧 Linux Distribution Compatibility

YERD now supports automatic dependency management across all major Linux distributions:

| Distribution | Support Level | Package Manager | Auto-Install |
|--------------|---------------|----------------|--------------|
| **ArchLinux** | ✅ Full | pacman | ✅ Yes |
| **Omarchy** | ✅ Full | pacman | ✅ Yes |
| **Ubuntu** | ✅ Full | apt | ✅ Yes |
| **Debian** | ✅ Full | apt | ✅ Yes |
| **Fedora** | ✅ Full | dnf | ✅ Yes |
| **RHEL/CentOS** | ✅ Full | yum | ✅ Yes |
| **openSUSE** | ✅ Full | zypper | ✅ Yes |

**Automatic dependency installation**: YERD automatically detects your distribution and installs the correct build dependencies using your system's package manager. No manual setup required!

## 🚀 Quick Installation

### Option 1: One-Line Install (Recommended)

The installation script automatically detects your system and handles all setup:

```bash
curl -sSL https://raw.githubusercontent.com/LumoSolutions/yerd/main/install.sh | bash
```

**Alternative with wget:**
```bash
wget -qO- https://raw.githubusercontent.com/LumoSolutions/yerd/main/install.sh | bash
```

**Smart installer features:**
- 🔍 Auto-detects Linux distribution and architecture
- ⚠️ Checks for existing YERD installation (prompts before overwrite)
- 📦 Downloads latest release automatically
- 🔧 Creates necessary directories with proper permissions
- ✅ Verifies installation and provides next steps

### Option 2: Download Pre-built Binary

1. Go to the [Releases](https://github.com/LumoSolutions/yerd/releases) page
2. Download the appropriate binary for your system (usually `yerd_*_linux_amd64.tar.gz`)
3. Extract and install:

```bash
# Download and extract (replace with actual version)
wget https://github.com/LumoSolutions/yerd/releases/download/v1.0.0/yerd_1.0.0_linux_amd64.tar.gz
tar -xzf yerd_1.0.0_linux_amd64.tar.gz

# Install system-wide
sudo mv yerd /usr/local/bin/yerd
sudo chmod +x /usr/local/bin/yerd

# Create YERD directories (important!)
sudo mkdir -p /opt/yerd/{bin,php,etc}

# Verify installation
yerd --help
yerd status
```

### Option 3: Install Script with Options

```bash
# Download install script
wget https://raw.githubusercontent.com/LumoSolutions/yerd/main/install.sh
chmod +x install.sh

# Install latest version
./install.sh

# Install specific version
./install.sh --version 1.0.0
```

### ✅ Post-Installation Verification

After successful installation, verify YERD is working correctly:

```bash
# Check installation
yerd --version

# View system status
yerd status  

# Test with help
yerd --help
yerd php --help

# Check available PHP versions
yerd php list
```

**Expected output:** You should see the YERD splash screen and version information without errors.

**Troubleshooting:**
- If `yerd` command is not found, ensure `/usr/local/bin` is in your PATH
- Add to your shell profile: `export PATH="/usr/local/bin:$PATH"`
- Reload shell: `source ~/.bashrc` or restart terminal

## 🏗️ Building from Source

**Requirements:**
- Go 1.21+ installed
- Git for cloning the repository
- Build tools (gcc, make) - usually pre-installed

```bash
# Clone the repository
git clone https://github.com/LumoSolutions/yerd.git
cd yerd

# Install dependencies and build
go mod tidy
go build -o yerd .

# Install system-wide
sudo cp yerd /usr/local/bin/yerd
sudo chmod +x /usr/local/bin/yerd

# Create YERD directories
sudo mkdir -p /opt/yerd/{bin,php,etc}

# Verify installation
yerd --help
yerd status
```

## 💻 Usage

YERD provides a clean, intuitive command-line interface with beautiful colored output:

```bash
$ yerd
██╗   ██╗███████╗██████╗ ██████╗ 
╚██╗ ██╔╝██╔════╝██╔══██╗██╔══██╗
 ╚████╔╝ █████╗  ██████╔╝██║  ██║
  ╚██╔╝  ██╔══╝  ██╔══██╗██║  ██║
   ██║   ███████╗██║  ██║██████╔╝
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═════╝
                     v1.0.2

A powerful, developer-friendly tool for managing PHP versions
and local development environments with ease

https://github.com/LumoSolutions/yerd
```

### 📝 Commands

**Top-Level Commands:**
| Command | Description | Requires Sudo | Example |
|---------|-------------|---------------|---------|
| `yerd` | Show splash screen and help | No | `yerd` |
| `yerd status` | Show system status and build environment | No | `yerd status` |
| `yerd php` | Manage PHP versions (see PHP commands below) | Varies | `yerd php --help` |

**PHP Management Commands:**
| Command | Description | Requires Sudo | Example |
|---------|-------------|---------------|---------|
| `yerd php list` | List available and installed PHP versions | No | `yerd php list` |
| `yerd php add <version>` | Build and install PHP from source | Yes | `sudo yerd php add 8.4` |
| `yerd php remove <version>` | Remove an installed PHP version | Yes | `sudo yerd php remove 8.3` |
| `yerd php cli <version>` | Set default CLI PHP version | Yes | `sudo yerd php cli 8.4` |
| `yerd php update [version]` | Update PHP versions to latest releases | Yes | `sudo yerd php update 8.4` |
| `yerd php doctor [version]` | Diagnose installations and environment | No | `yerd php doctor 8.3` |

**Version Format Flexibility:**
- ✅ **Short format**: `8.4`, `8.3`, `8.2`
- ✅ **With prefix**: `php8.4`, `php8.3`, `php8.2`
- ✅ **Case insensitive**: `PHP8.4`, `Php8.3`

📝 **Build Logging**: The `add` command creates detailed build logs in `~/.config/yerd/` that are automatically deleted on success or preserved for troubleshooting on failure.

### 🚀 Quick Start Example

```bash
# 1. Check system status
$ yerd status
██╗   ██╗███████╗██████╗ ██████╗ 
╚██╗ ██╔╝██╔════╝██╔══██╗██╔══██╗
 ╚████╔╝ █████╗  ██████╔╝██║  ██║
  ╚██╔╝  ██╔══╝  ██╔══██╗██║  ██║
   ██║   ███████╗██║  ██║██████╔╝
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═════╝
                     v1.0.2

📊 YERD Status
├─ Installed versions: 0
├─ Current CLI: None set
└─ Config: ~/.config/yerd/config.json

🔍 System PHP Check
└─ ✅ No conflicts - ready for YERD management

# 2. Check what's available
$ yerd php list
██╗   ██╗███████╗██████╗ ██████╗ 
╚██╗ ██╔╝██╔════╝██╔══██╗██╔══██╗
 ╚████╔╝ █████╗  ██████╔╝██║  ██║
  ╚██╔╝  ██╔══╝  ██╔══██╗██║  ██║
   ██║   ███████╗██║  ██║██████╔╝
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═════╝
                     v1.0.2

+---------+-----------+-----+----------------+
| VERSION | INSTALLED | CLI |    UPDATES     |
+---------+-----------+-----+----------------+
| PHP 8.2 | No        |     | Latest: 8.2.29 |
| PHP 8.3 | No        |     | Latest: 8.3.14 |
| PHP 8.4 | No        |     | Latest: 8.4.11 |
+---------+-----------+-----+----------------+

# 3. Install PHP 8.4 (note: flexible version format)
$ sudo yerd php add 8.4
██╗   ██╗███████╗██████╗ ██████╗ 
╚██╗ ██╔╝██╔════╝██╔══██╗██╔══██╗
 ╚████╔╝ █████╗  ██████╔╝██║  ██║
  ╚██╔╝  ██╔══╝  ██╔══██╗██║  ██║
   ██║   ███████╗██║  ██║██████╔╝
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═════╝
                     v1.0.2

📦 Installing latest PHP 8.4: 8.4.11
📋 Detected: Arch Linux
⚠️  Installing build dependencies...
✓ Build dependencies installed
Building PHP 8.4 from source...
[/] Downloading php-8.4.11.tar.gz...
✓ Download complete
[|] Extracting source code...
✓ Source extracted
[-] Configuring build...
✓ Configure complete
[\] Building PHP (this may take several minutes)...
✓ Build complete
[/] Installing PHP...
✓ Install complete
✓ PHP 8.4 built and installed successfully
[|] Locating PHP 8.4 binary...
✓ Found PHP binary at: /opt/yerd/php/php8.4/bin/php
[-] Creating symlinks...
✓ Symlinks created successfully
✓ PHP installation verified
✓ Default php.ini created
✓ PHP 8.4 installed successfully

# 4. Set as default CLI version
$ sudo yerd php cli 8.4
██╗   ██╗███████╗██████╗ ██████╗ 
╚██╗ ██╔╝██╔════╝██╔══██╗██╔══██╗
 ╚████╔╝ █████╗  ██████╔╝██║  ██║
  ╚██╔╝  ██╔══╝  ██╔══██╗██║  ██║
   ██║   ███████╗██║  ██║██████╔╝
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═════╝
                     v1.0.2

Setting PHP 8.4 as CLI version...
✓ PHP 8.4 is now the default CLI version
Verify with: php -v

# 5. Verify it's working
$ php -v
PHP 8.4.11 (cli) (built: Aug 12 2025 18:52:45) (NTS)

# 6. Check final status with detailed paths
$ yerd status
📦 Installed PHP Versions
└─ 🎯 PHP 8.4 (Current CLI)
   ├─ Binary: /opt/yerd/bin/php8.4
   ├─ Config: /opt/yerd/etc/php8.4/php.ini
   └─ Install: /opt/yerd/php/php8.4/

# 7. Later, when updates are available
$ sudo yerd php update -y
🔄 Auto-updating all 1 PHP version(s)...
📦 Installing updated PHP 8.4...
✅ PHP 8.4 updated successfully to 8.4.12
```

**Key Features Demonstrated:**
- 🔧 **Automatic dependency management** across Linux distributions
- 🏗️ **Source-based compilation** with build progress indicators
- 📝 **Default php.ini creation** for immediate usability
- 🔄 **Flexible version formats**: `8.4`, `php8.4`, `PHP8.4` all work
- 📊 **Detailed status reporting** with binary and config paths
- ⚡ **Zero-downtime updates** with automatic cleanup

## 📚 Complete Command Reference

### 🏠 `yerd` - Main Command
Shows the beautiful YERD splash screen and help information.

```bash
$ yerd
██╗   ██╗███████╗██████╗ ██████╗ 
╚██╗ ██╔╝██╔════╝██╔══██╗██╔══██╗
 ╚████╔╝ █████╗  ██████╔╝██║  ██║
  ╚██╔╝  ██╔══╝  ██╔══██╗██║  ██║
   ██║   ███████╗██║  ██║██████╔╝
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═════╝
                     v1.0.2

A powerful, developer-friendly tool for managing PHP versions
and local development environments with ease

https://github.com/LumoSolutions/yerd
```

### 📊 `yerd status` - System Status
Shows YERD configuration, system conflicts, and detailed information about installed PHP versions.

**Usage**: `yerd status`  
**Permissions**: None required

```bash
$ yerd status
📦 Installed PHP Versions
├─ 📌 PHP 8.3
│  ├─ Binary: /opt/yerd/bin/php8.3
│  ├─ Config: /opt/yerd/etc/php8.3/php.ini
│  └─ Install: /opt/yerd/php/php8.3/
│
└─ 🎯 PHP 8.4 (Current CLI)
   ├─ Binary: /opt/yerd/bin/php8.4
   ├─ Config: /opt/yerd/etc/php8.4/php.ini
   └─ Install: /opt/yerd/php/php8.4/
```

## 🔧 PHP Management Commands

### 📋 `yerd php list` - List PHP Versions  
Lists all available PHP versions and their installation status.

**Usage**: `yerd php list`  
**Permissions**: None required

```bash
$ yerd php list
+---------+-----------+-----+----------------+
| VERSION | INSTALLED | CLI |    UPDATES     |
+---------+-----------+-----+----------------+
| PHP 8.2 | No        |     | Latest: 8.2.29 |
| PHP 8.3 | Yes       |     | Yes            |
| PHP 8.4 | Yes       | *   | No             |
+---------+-----------+-----+----------------+
```

**Legend:**
- ✅ **Yes**: Version is installed
- ❌ **No**: Version not installed  
- ⭐ **\***: Current CLI version
- 🔄 **Yes/No**: Update available for installed versions
- 📦 **Latest: X.X.X**: Latest available version for uninstalled versions

### ⬇️ `yerd php add <version>` - Build PHP from Source
Downloads and builds a specific PHP version from official source code.

**Usage**: `sudo yerd php add 8.4` or `sudo yerd php add php8.4`  
**Permissions**: Requires `sudo`  
**Supported Versions**: `8.2`, `8.3`, `8.4` (with or without `php` prefix)

```bash
$ sudo yerd php add 8.4
📦 Installing latest PHP 8.4: 8.4.11
📋 Detected: Arch Linux
⚠️  Installing build dependencies...
✓ Build dependencies installed
Building PHP 8.4 from source...
[/] Downloading php-8.4.11.tar.gz...
✓ Download complete
[|] Extracting source code...
✓ Source extracted
[-] Configuring build...
✓ Configure complete
[\] Building PHP (this may take several minutes)...
✓ Build complete
[/] Installing PHP...
✓ Install complete
✓ PHP 8.4 built and installed successfully
✓ Default php.ini created
✓ PHP 8.4 installed successfully
```

### 🗑️ `yerd php remove <version>` - Remove PHP Version
Removes an installed PHP version and cleans up symlinks.

**Usage**: `sudo yerd php remove 8.3` or `sudo yerd php remove php8.3`  
**Permissions**: Requires `sudo`

```bash
$ sudo yerd php remove 8.3
⚠️  Warning: You are about to remove PHP 8.3
Continue? (y/N): y
🗑️  Removing PHP 8.3...
✓ PHP 8.3 removed successfully
```

### 🔄 `yerd php cli <version>` - Set CLI Version
Sets the default PHP version for command line usage.

**Usage**: `sudo yerd php cli 8.4` or `sudo yerd php cli php8.4`  
**Permissions**: Requires `sudo`

```bash
$ sudo yerd php cli 8.4
Setting PHP 8.4 as CLI version...
✓ PHP 8.4 is now the default CLI version
Verify with: php -v
```

### 🔄 `yerd php update [version]` - Update PHP Versions
Updates installed PHP versions to their latest releases.

**Usage**: 
- `sudo yerd php update` (update all versions)
- `sudo yerd php update 8.4` (update specific version)
- `sudo yerd php update -y` (auto-confirm updates)

**Permissions**: Requires `sudo`

```bash
$ sudo yerd php update -y
🔄 Auto-updating all 2 PHP version(s)...
📦 Installing updated PHP 8.3...
✅ PHP 8.3 updated successfully to 8.3.15
📦 Installing updated PHP 8.4...  
✅ PHP 8.4 updated successfully to 8.4.12
```

### 🩺 `yerd php doctor [version]` - Diagnostic Tool
Runs comprehensive diagnostics to troubleshoot installation issues.

**Usage**: 
- `yerd php doctor` (general diagnostics)
- `yerd php doctor 8.4` (version-specific diagnostics)

**Permissions**: None required

```bash
$ yerd php doctor 8.4
🩺 YERD Doctor - System Diagnostics

1️⃣  System Requirements
├─ ✅ Build tool: gcc (Available)
├─ ✅ Build tool: make (Available)
└─ ✅ All requirements satisfied

4️⃣  PHP 8.4 Diagnostics  
├─ ✅ YERD status: Installed
├─ ✅ Binary found: /usr/local/bin/php8.4
└─ ℹ️  Version info: PHP 8.4.11 (cli) (built: Aug 12 2025 18:52:45) (NTS)

✅ Diagnostics complete. No issues found.
```

---

*This documentation reflects the latest YERD CLI structure with PHP subcommands and flexible version format support.*
