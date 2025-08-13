# YERD - PHP Version Manager for Linux

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

## 🎯 What is YERD?

YERD is a Linux PHP version manager that compiles PHP from official source code, giving you complete control over your PHP installations. Perfect for both production servers and development environments.

**Key Benefits:**
- 🚀 **Multiple PHP versions** running simultaneously
- ⚡ **Instant CLI switching** between versions  
- 🛡️ **Safe isolation** - never conflicts with system PHP
- 🧩 **Rich extension support** with automatic dependencies
- 🏗️ **Source-based builds** for maximum reliability
- 🌐 **Multi-distro support** - works on all major Linux distributions

## 🚀 Quick Start

### Installation (One Command)

```bash
curl -sSL https://raw.githubusercontent.com/LumoSolutions/yerd/main/install.sh | bash
```

### Basic Usage

```bash
# Install PHP 8.4
sudo yerd php add 8.4

# Install Composer (optional)
sudo yerd php composer

# Set as default CLI version  
sudo yerd php cli 8.4

# Verify installation
php -v  # PHP 8.4.11 (cli)
composer --version  # Latest Composer

# View all versions
yerd php list
```

## 📋 System Requirements

- **OS**: Any Linux distribution (Ubuntu, Debian, Arch, Fedora, RHEL, openSUSE, etc.)
- **Permissions**: `sudo` access for installation operations
- **Dependencies**: Automatically installed based on your distribution
- **Internet**: Required for downloading PHP source and updates

## 💻 Commands

### Top-Level Commands

| Command | Description | Example |
|---------|-------------|---------|
| `yerd status` | System status overview | Shows conflicts, paths |
| `yerd --help` | Show help information | Display usage guide |
| `yerd --version` | Show YERD version | Display current version |

### PHP Commands

#### Installation & Removal

| Command | Description | Example |
|---------|-------------|---------|
| `sudo yerd php add 8.4` | Install PHP version from source | Builds PHP 8.4 with default extensions |
| `sudo yerd php remove 8.3` | Remove PHP version | Cleans up completely |
| `sudo yerd php composer` | Install/update Composer | Downloads latest stable Composer |

#### Management

| Command | Description | Example |
|---------|-------------|---------|
| `yerd php list` | List available/installed versions | Shows status and updates |
| `sudo yerd php cli 8.4` | Set CLI version | Makes `php` command use 8.4 |

#### Extensions

| Command | Description | Example |
|---------|-------------|---------|
| `yerd php extensions 8.3` | View extensions | Shows installed/available |
| `sudo yerd php extensions add 8.3 mysqli gd` | Add extensions | Rebuilds PHP automatically |
| `sudo yerd php extensions remove 8.3 curl` | Remove extensions | Smart dependency management |

#### Maintenance

| Command | Description | Example |
|---------|-------------|---------|
| `sudo yerd php rebuild 8.3` | Force rebuild | Useful for troubleshooting |
| `sudo yerd php update` | Update versions | Gets latest patches |
| `yerd php doctor` | Diagnostics | Troubleshoot issues |

## 🧩 Extension Management

YERD includes a powerful extension system with 40+ supported extensions:

**Popular Extensions:**
- **Database**: `mysqli`, `pdo-mysql`, `pgsql`, `sqlite3`
- **Graphics**: `gd`, `jpeg`, `freetype`, `exif`  
- **Network**: `curl`, `openssl`, `sockets`
- **Core**: `mbstring`, `opcache`, `zip`, `json`

**Smart Features:**
- 📦 **Auto-dependencies**: Installs required system packages automatically
- 🔄 **Configuration rollback**: Reverts changes if build fails
- 🌐 **Multi-distro support**: Works across all Linux distributions
- ⚡ **Smart rebuilding**: Only rebuilds when extensions change

```bash
# View extensions for PHP 8.3
yerd php extensions 8.3

# Add extensions (rebuilds PHP automatically)
sudo yerd php extensions add 8.3 mysqli gd curl

# Replace all extensions with specific list
sudo yerd php extensions replace 8.3 mbstring opcache curl
```

## 📦 Composer Management

YERD includes integrated Composer management for PHP dependency handling:

**Features:**
- 🚀 **Latest Version**: Always downloads the latest stable Composer
- 🔄 **Easy Updates**: Simple command to update to newest version
- 🌐 **Global Access**: Available system-wide via `/usr/local/bin/composer`
- 🛡️ **YERD Integration**: Stored in YERD directory structure for consistency

```bash
# Install or update Composer
sudo yerd php composer

# Verify installation
composer --version

# Use Composer normally
composer install
composer require vendor/package
composer update
```

**File Locations:**
- **Source**: `/opt/yerd/bin/composer.phar`
- **Global Link**: `/usr/local/bin/composer`

## 🔧 Advanced Features

### Multi-Distribution Support

YERD automatically detects your Linux distribution and installs appropriate dependencies:

| Distribution | Package Manager | Example Extensions |
|--------------|----------------|-------------------|
| Ubuntu/Debian | `apt` | `libgd-dev`, `libcurl4-openssl-dev` |
| Arch/Manjaro | `pacman` | `gd`, `curl` |
| Fedora/RHEL | `dnf`/`yum` | `gd-devel`, `libcurl-devel` |
| openSUSE | `zypper` | `gd-devel`, `libcurl-devel` |

### Safety Features

- **🛡️ System Protection**: Never overwrites existing PHP installations
- **🔐 Privilege Separation**: Build processes run as user, not root
- **📸 Configuration Backup**: Automatic rollback if extension changes fail
- **🔍 Conflict Detection**: Warns about system PHP conflicts before installation

### Performance Optimizations

- **🚀 Multi-core compilation**: Uses all CPU cores for faster builds
- **⚡ Smart caching**: Minimizes API requests with intelligent version caching
- **📦 Parallel downloads**: Efficient source code retrieval
- **🧹 Automatic cleanup**: Removes temporary files after successful builds

## 🗂️ File Locations

```
/opt/yerd/                      # Main directory
├── php/                        # PHP installations
│   ├── php8.3/                # PHP 8.3 installation
│   └── php8.4/                # PHP 8.4 installation
├── bin/                        # YERD-managed binaries
│   ├── php8.3 -> /opt/yerd/php/php8.3/bin/php
│   ├── php8.4 -> /opt/yerd/php/php8.4/bin/php
│   └── composer.phar           # Composer installation
└── etc/                        # Configuration files

/usr/local/bin/                 # Global symlinks
├── php -> /opt/yerd/bin/php8.4 # Current CLI version
├── php8.3 -> /opt/yerd/bin/php8.3
├── php8.4 -> /opt/yerd/bin/php8.4
└── composer -> /opt/yerd/bin/composer.phar

~/.config/yerd/config.json      # User configuration
```

## 🚀 Common Workflows

### Development Environment Setup

```bash
# Install multiple PHP versions for testing
sudo yerd php add 8.1
sudo yerd php add 8.2
sudo yerd php add 8.3
sudo yerd php add 8.4

# Install Composer for dependency management
sudo yerd php composer

# Set 8.3 as default CLI
sudo yerd php cli 8.3

# Test with different versions
php8.1 -v  # PHP 8.1.x
php8.2 -v  # PHP 8.2.x
php8.3 -v  # PHP 8.3.x (also available as 'php')
php8.4 -v  # PHP 8.4.x
composer --version  # Composer globally available
```

### Production Server Management

```bash
# Install specific version for production
sudo yerd php add 8.3

# Install Composer for dependency management
sudo yerd php composer

# Add production extensions
sudo yerd php extensions add 8.3 mysqli pdo-mysql opcache gd curl openssl

# Set as CLI version
sudo yerd php cli 8.3

# Monitor and update
yerd status
sudo yerd php update -y
```

### Extension Development

```bash
# Add development extensions
sudo yerd php extensions add 8.3 mysqli gd curl json xml

# Force rebuild after system updates
sudo yerd php rebuild 8.3

# Troubleshoot issues
yerd php doctor 8.3
```

## 🔍 Troubleshooting

### Common Issues

1. **Build failures**: Run `yerd php doctor <version>` for diagnostics
2. **Permission errors**: Ensure you're using `sudo` for installation commands
3. **System conflicts**: Use `yerd status` to check for existing PHP installations
4. **Extension issues**: Check logs in `~/.config/yerd/` directory

### Getting Help

```bash
# General diagnostics
yerd php doctor

# Version-specific diagnostics  
yerd php doctor 8.3

# System status
yerd status

# Command help
yerd --help
yerd php --help
yerd php extensions --help
```

## 🏗️ Building from Source

If you want to build YERD yourself:

```bash
git clone https://github.com/LumoSolutions/yerd.git
cd yerd
go mod tidy
go build -o yerd .
sudo cp yerd /usr/local/bin/yerd
```

## 🤝 Contributing

We welcome contributions! Please see our [contributing guidelines](CONTRIBUTING.md) for details.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🔗 Links

- **Repository**: https://github.com/LumoSolutions/yerd
- **Issues**: https://github.com/LumoSolutions/yerd/issues
- **Releases**: https://github.com/LumoSolutions/yerd/releases
- **Developer**: [LumoSolutions](https://github.com/LumoSolutions)

---

*Made with ❤️ for the PHP community*