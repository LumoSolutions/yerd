# YERD - The Ultimate PHP Development Environment Manager

<div align="center">

![YERD Logo](.meta/yerd_logo.jpg)

**Transform your PHP development workflow with intelligent version management and seamless local environments**

https://github.com/LumoSolutions/yerd

[![Release](https://img.shields.io/github/v/release/LumoSolutions/yerd)](https://github.com/LumoSolutions/yerd/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS-green.svg)](https://github.com/LumoSolutions/yerd)

</div>

---

## 🚀 Why YERD?

**Stop wrestling with PHP versions. Start building.**

YERD revolutionizes PHP development by providing a zero-friction environment manager that just works. Whether you're juggling legacy projects, testing against multiple PHP versions, or deploying production applications, YERD eliminates the complexity and lets you focus on what matters - your code.

### ✨ Key Benefits

- **🎯 Zero Conflicts** - Complete isolation from system PHP. Never break production again.
- **⚡ Instant Switching** - Change PHP versions in milliseconds, not minutes
- **🛠️ Production-Grade** - Built from official PHP source for maximum reliability
- **🔧 30+ Extensions** - Pre-configured with smart dependency management
- **🌐 Local Development** - Integrated nginx with automatic HTTPS for every site
- **🔒 Chrome-Trusted SSL** - Self-signed certificates that browsers actually trust
- **📦 Composer Included** - Latest Composer managed automatically
- **🔄 Self-Updating** - Stay current with one-command updates

### 🏆 Perfect For

- **Development Teams** - Standardize environments across your entire team
- **Agencies** - Manage multiple client projects with different PHP requirements
- **Open Source Maintainers** - Test against multiple PHP versions effortlessly
- **DevOps Engineers** - Deploy consistent, reproducible PHP environments
- **Freelancers** - Switch between client projects without environment conflicts

## ⚡ Quick Start

```bash
curl -sSL https://raw.githubusercontent.com/LumoSolutions/yerd/main/install.sh | sudo bash
```

Then:

```bash
# Install PHP 8.4 with a single command
sudo yerd php 8.4 install

# Set as your default CLI
sudo yerd php 8.4 cli

# You're ready to code!
php -v  # PHP 8.4.x
```

## 📋 System Requirements

- **Operating Systems**: Linux (all distributions), macOS
- **Architecture**: x86_64, ARM64, 32-bit
- **Privileges**: sudo access for installations
- **Dependencies**: Automatically managed per distribution

## 🎯 Core Features

### Multiple PHP Versions
Run PHP 8.1, 8.2, 8.3, and 8.4 simultaneously without conflicts. Each version is completely isolated with its own configuration, extensions, and FPM service.

### Intelligent Extension Management
Choose from 30+ extensions with automatic dependency resolution. YERD handles the complexity of building PHP with your exact requirements.

### Web Development Ready
Integrated nginx support transforms YERD into a complete local development environment. Create sites with custom domains, automatic SSL certificates, and per-site PHP versions. Every site gets a chrome-trusted HTTPS certificate automatically.

### Zero-Configuration Composer
Latest Composer version managed by YERD - always up-to-date, always available globally.

### Enterprise-Grade Reliability
Built from official PHP source code with production-ready configurations. Perfect for both development and production deployments.

## 📚 Complete Command Reference

### PHP Version Management

#### Installation Commands

```bash
# Install a PHP version (8.1, 8.2, 8.3, or 8.4)
sudo yerd php 8.4 install

# Install with fresh source (bypass cache)
sudo yerd php 8.4 install --nocache
```

#### Version Control

```bash
# List all installed PHP versions
yerd php list

# Show detailed PHP status
yerd php status

# Set default CLI version
sudo yerd php 8.4 cli

# Force CLI version update
sudo yerd php 8.4 cli --force
```

#### Extension Management

```bash
# List available extensions for a version
yerd php 8.3 extensions list

# Add extensions (automatically rebuilds PHP)
sudo yerd php 8.3 extensions add gd mysqli opcache

# Remove extensions
sudo yerd php 8.3 extensions remove gd

# Add extensions and rebuild immediately
sudo yerd php 8.3 extensions add gd --rebuild
```

**Available Extensions**: mbstring, bcmath, opcache, curl, openssl, zip, sockets, mysqli, pdo-mysql, gd, jpeg, freetype, xml, json, session, hash, filter, pcre, zlib, bz2, iconv, intl, pgsql, pdo-pgsql, sqlite3, pdo-sqlite, fileinfo, exif, gettext, gmp, ldap, soap, ftp, pcntl

#### Maintenance Operations

```bash
# Update PHP to latest patch version
sudo yerd php 8.4 update

# Rebuild PHP (useful after system updates)
sudo yerd php 8.4 rebuild

# Rebuild with config regeneration
sudo yerd php 8.4 rebuild --config

# Uninstall a PHP version
sudo yerd php 8.2 uninstall

# Skip confirmation prompts
sudo yerd php 8.2 uninstall --yes
```

### Composer Management

```bash
# Install Composer
sudo yerd composer install

# Update to latest version
sudo yerd composer update

# Remove Composer
sudo yerd composer uninstall
```

### Web Services (nginx)

```bash
# Install web components
sudo yerd web install

# Remove web components
sudo yerd web uninstall
```

### Site Management

```bash
# List all sites
yerd sites list

# Add a new site (automatically creates HTTPS certificate)
sudo yerd sites add /path/to/project

# Add with custom domain and PHP version
sudo yerd sites add /var/www/myapp --domain myapp.test --php 8.3

# Specify public directory
sudo yerd sites add /var/www/myapp --folder public

# Remove a site
sudo yerd sites remove /path/to/project

# Update site configuration
sudo yerd sites set php 8.4 myapp.test
```

**🔒 Automatic SSL Certificates**: Every site is served over HTTPS by default with a chrome-trusted SSL certificate, signed by a YERD Certificate Authority generated and managed on your system. No more browser warnings!

### Self-Update

```bash
# Check for and install updates
sudo yerd update

# Auto-confirm updates
sudo yerd update --yes
```

## 🔄 Typical Workflows

### New Project Setup

```bash
# 1. Install required PHP version
sudo yerd php 8.4 install

# 2. Add necessary extensions
sudo yerd php 8.4 extensions add mysqli gd opcache curl

# 3. Install Composer
sudo yerd composer install

# 4. Set as CLI default
sudo yerd php 8.4 cli

# 5. Install web components
sudo yerd web install

# 6. Add your project site (automatic HTTPS)
sudo yerd sites add /var/www/myproject --domain myproject.test --php 8.4

# 7. Start developing with HTTPS!
cd /var/www/myproject
composer install
# Access at: https://myproject.test
```

### Multi-Version Testing

```bash
# Install all PHP versions
sudo yerd php 8.1 install
sudo yerd php 8.2 install
sudo yerd php 8.3 install
sudo yerd php 8.4 install

# Test your code across versions
php8.1 vendor/bin/phpunit
php8.2 vendor/bin/phpunit
php8.3 vendor/bin/phpunit
php8.4 vendor/bin/phpunit

# Or switch CLI versions as needed
sudo yerd php 8.1 cli && phpunit
sudo yerd php 8.4 cli && phpunit
```

### Production Deployment

```bash
# Install specific PHP version
sudo yerd php 8.3 install

# Add production extensions
sudo yerd php 8.3 extensions add opcache mysqli pdo-mysql curl openssl

# Set as system default
sudo yerd php 8.3 cli

# Keep updated
sudo yerd update --yes
sudo yerd php 8.3 update
```

## 🏗️ Architecture

YERD provides a clean, organized structure:

```
/opt/yerd/
├── bin/        # PHP binaries and Composer
├── php/        # PHP installations
├── etc/        # Configuration files
└── web/        # nginx and certificates

/usr/local/bin/
├── php         # Current CLI version
├── php8.1      # Direct version access
├── php8.2
├── php8.3
├── php8.4
└── composer    # Global Composer

~/.config/yerd/config.json  # User configuration
```

## 🛠️ Advanced Features

### 🔒 Automatic SSL Certificate Management
YERD includes a complete SSL certificate infrastructure for local development:
- **YERD Certificate Authority**: A local CA is generated and managed on your system
- **Automatic Certificate Generation**: Every site gets its own SSL certificate automatically
- **Chrome/Browser Trust**: Certificates are signed by the YERD CA, eliminating browser warnings
- **HTTPS by Default**: All sites are served over HTTPS (port 443) with HTTP redirect
- **Zero Configuration**: Just add a site and SSL is handled automatically

### FPM Service Management
Each PHP version includes its own FPM service, managed via systemd:
- Service names: `yerd-php{version}-fpm`
- Sockets: `/opt/yerd/php/run/php{version}-fpm.sock`
- Automatic start on boot
- Graceful reloads during updates

### Smart Dependency Management
YERD automatically detects your Linux distribution and installs appropriate packages:
- **Ubuntu/Debian**: Uses `apt` with development libraries
- **Arch/Manjaro**: Uses `pacman` with build tools
- **Fedora/RHEL**: Uses `dnf`/`yum` with devel packages
- **openSUSE**: Uses `zypper` with development patterns

### Security First
- Complete isolation from system PHP
- Privilege separation during builds
- Automatic configuration backups
- Chrome-trusted SSL certificates for all local sites
- No root processes except installation

## 🚨 Troubleshooting

```bash
# Check system status
yerd php status

# Verify installations
yerd php list

# Force rebuild after system updates
sudo yerd php 8.4 rebuild --config

# Check service status
systemctl status yerd-php84-fpm

# View logs
journalctl -u yerd-php84-fpm
```

## 🤝 Contributing

We welcome contributions! YERD is built with Go and uses the Cobra CLI framework. Check out our [contributing guidelines](CONTRIBUTING.md) to get started.

```bash
# Clone and build from source
git clone https://github.com/LumoSolutions/yerd.git
cd yerd
go build -o yerd main.go
```

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🔗 Resources

- **GitHub**: [LumoSolutions/yerd](https://github.com/LumoSolutions/yerd)
- **Issues**: [Report bugs or request features](https://github.com/LumoSolutions/yerd/issues)
- **Releases**: [Download latest version](https://github.com/LumoSolutions/yerd/releases)

---

<div align="center">

**Built with ❤️ for developers who demand more from their tools**

*Stop managing environments. Start shipping code.*

</div>
