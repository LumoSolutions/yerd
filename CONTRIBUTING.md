# Contributing to YERD

Thank you for your interest in contributing to YERD! We welcome contributions from the community and appreciate your help in making YERD better.

## 🚀 Quick Start

1. **Fork** the repository
2. **Clone** your fork locally
3. **Create** a feature branch
4. **Make** your changes
5. **Test** thoroughly
6. **Submit** a pull request

## 📋 Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Code Standards](#code-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Pull Request Process](#pull-request-process)
- [Code Review](#code-review)
- [Community Guidelines](#community-guidelines)

## 🏁 Getting Started

### Prerequisites

- **Go 1.24+** installed
- **Git** for version control
- **Linux or macOS environment** (Ubuntu, Debian, Arch, Fedora, macOS, etc.)
- **sudo access** for testing installation features

### Fork and Clone

1. **Fork the repository** on GitHub:
   - Navigate to https://github.com/LumoSolutions/yerd
   - Click the "Fork" button in the top-right corner
   - Select your GitHub account as the destination

2. **Clone your fork locally:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/yerd.git
   cd yerd
   ```

3. **Add the upstream remote:**
   ```bash
   git remote add upstream https://github.com/LumoSolutions/yerd.git
   ```

4. **Verify remotes:**
   ```bash
   git remote -v
   # origin    https://github.com/YOUR_USERNAME/yerd.git (fetch)
   # origin    https://github.com/YOUR_USERNAME/yerd.git (push)
   # upstream  https://github.com/LumoSolutions/yerd.git (fetch)
   # upstream  https://github.com/LumoSolutions/yerd.git (push)
   ```

## 🛠️ Development Setup

### Build and Test

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Build YERD:**
   ```bash
   go build -o yerd .
   ```

3. **Test basic functionality:**
   ```bash
   ./yerd --help
   ./yerd --version
   ./yerd php --help
   ```

### Development Workflow

1. **Keep your fork updated:**
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   git push origin main
   ```

2. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-description
   ```

3. **Make your changes** (see [Making Changes](#making-changes))

4. **Commit and push:**
   ```bash
   git add .
   git commit -m "Add descriptive commit message"
   git push origin feature/your-feature-name
   ```

## 📁 Project Structure

YERD follows a clean, modular architecture using the Cobra CLI framework:

```
yerd/
├── cmd/                         # CLI command definitions
│   ├── root.go                 # Root command and initialization
│   ├── update.go               # Self-update command
│   ├── php.go                  # PHP parent command
│   ├── composer.go             # Composer parent command
│   ├── web.go                  # Web services parent command
│   ├── sites.go                # Sites parent command
│   ├── php/                    # PHP subcommands
│   │   ├── cli.go              # Set CLI version
│   │   ├── extensions.go       # Manage extensions
│   │   ├── install.go          # Install PHP version
│   │   ├── list.go             # List installed versions
│   │   ├── php_version.go      # Version-specific commands
│   │   ├── rebuild.go          # Rebuild PHP
│   │   ├── status.go           # Show PHP status
│   │   ├── uninstall.go        # Uninstall PHP version
│   │   └── update.go           # Update PHP version
│   ├── composer/               # Composer subcommands
│   │   ├── install.go          # Install Composer
│   │   ├── uninstall.go        # Uninstall Composer
│   │   └── update.go           # Update Composer
│   ├── web/                    # Web service subcommands
│   │   ├── install.go          # Install nginx
│   │   └── uninstall.go        # Uninstall nginx
│   └── sites/                  # Site management subcommands
│       ├── add.go              # Add new site
│       ├── list.go             # List sites
│       ├── remove.go           # Remove site
│       └── set.go              # Set site configuration
│
├── internal/                    # Internal packages (not exported)
│   ├── config/                 # Configuration management
│   │   ├── config.go           # Config file operations
│   │   ├── php.go              # PHP-specific config
│   │   └── web.go              # Web services config
│   ├── constants/              # Application constants
│   │   ├── constants.go        # General constants
│   │   ├── dependencies.go     # System dependencies
│   │   ├── nginx.go            # nginx-specific constants
│   │   └── php.go              # PHP versions and extensions
│   ├── installers/             # Installation logic
│   │   ├── composer/           # Composer installer
│   │   │   └── common.go       # Composer operations
│   │   ├── nginx/              # nginx installer
│   │   │   └── installer.go    # nginx build and install
│   │   └── php/                # PHP installer
│   │       ├── cli.go          # CLI version management
│   │       ├── extensions.go   # Extension management
│   │       ├── general.go      # Common PHP operations
│   │       ├── install.go      # PHP installation
│   │       ├── uninstall.go    # PHP removal
│   │       └── versions.go     # Version checking
│   ├── manager/                # Site and service management
│   │   ├── certificate.go      # SSL certificate generation
│   │   ├── manager.go          # Site manager
│   │   └── site.go             # Site operations
│   ├── utils/                  # Utility functions
│   │   ├── arrays.go           # Array helpers
│   │   ├── commands.go         # Command execution
│   │   ├── common.go           # Common utilities
│   │   ├── download.go         # File downloads
│   │   ├── file.go             # File operations
│   │   ├── hosts.go           # /etc/hosts management
│   │   ├── logger.go           # Logging utilities
│   │   ├── spinner.go          # CLI spinner
│   │   ├── systemd.go          # Systemd operations
│   │   ├── template.go         # Template rendering
│   │   ├── ui.go               # UI helpers
│   │   └── user.go             # User operations
│   └── version/                # Version information
│       └── version.go          # Version constants and splash
│
├── scripts/                    # Build and release scripts
│   └── build-releases.sh       # Multi-platform build script
│
├── main.go                     # Application entry point
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── README.md                   # Project documentation
├── CONTRIBUTING.md             # Contribution guidelines
├── LICENSE                     # MIT License
└── install.sh                  # Installation script
```

### Key Packages

- **cmd/**: Contains all CLI commands using Cobra framework
- **internal/config/**: Manages YERD configuration files
- **internal/constants/**: Defines PHP versions, extensions, and paths
- **internal/installers/**: Handles PHP, Composer, and nginx installation
- **internal/manager/**: Manages sites and SSL certificates
- **internal/utils/**: Common utilities for file ops, systemd, UI, etc.

## ✏️ Making Changes

### Branch Naming

Use descriptive branch names that indicate the type of change:

- **Features:** `feature/add-redis-extension`
- **Bug fixes:** `fix/fpm-socket-permission`
- **Documentation:** `docs/update-ssl-docs`
- **Refactoring:** `refactor/simplify-installer`

### Commit Messages

Write clear, descriptive commit messages:

```
Add support for Redis PHP extension

- Add Redis to available extensions list
- Include hiredis dependency for Redis
- Update extension validation logic
- Add Redis to documentation

Closes #123
```

**Format:**
- Use imperative mood ("Add", "Fix", "Update")
- First line: concise summary (50 chars or less)
- Blank line, then detailed explanation if needed
- Reference issues with "Closes #123" or "Fixes #456"

### Adding New Features

#### Adding a New PHP Extension

1. Update `internal/constants/php.go`:
   ```go
   // Add to availableExtensions map
   "redis": {
       Name:         "redis",
       ConfigFlag:   "--enable-redis",
       Dependencies: []string{"hiredis"},
   },
   ```

2. Update documentation in README.md

#### Adding a New Command

1. Create command file in appropriate `cmd/` subdirectory
2. Use Cobra command structure:
   ```go
   func BuildYourCommand() *cobra.Command {
       return &cobra.Command{
           Use:   "command",
           Short: "Brief description",
           Long:  `Detailed description`,
           Run: func(cmd *cobra.Command, args []string) {
               // Implementation
           },
       }
   }
   ```

3. Register command in parent command's init

## 📏 Code Standards

### Go Style

- Follow **standard Go conventions**
- Run **gofmt** before committing
- Use **meaningful variable and function names**
- Keep functions **focused and concise**
- Use **early returns** to reduce nesting

### Documentation

- All **exported functions** must have comments
- Use **Go doc comment format:**
  ```go
  // InstallPhp installs the specified PHP version with given extensions.
  // Returns error if installation fails or version is invalid.
  func InstallPhp(version string, extensions []string) error {
      // implementation
  }
  ```

### Error Handling

- Always **check and handle errors**
- Provide **context in error messages**
- Use **fmt.Errorf** for error wrapping:
  ```go
  if err := someFunction(); err != nil {
      return fmt.Errorf("failed to do something: %w", err)
  }
  ```

### Logging and Output

- Use **color package** for colored output
- Follow existing patterns for user feedback:
  ```go
  green := color.New(color.FgGreen)
  red := color.New(color.FgRed)
  
  green.Println("✓ Operation successful")
  red.Printf("❌ Error: %v\n", err)
  ```

## 🧪 Testing

### Manual Testing

Before submitting a PR, test your changes:

1. **Build successfully:**
   ```bash
   go build -o yerd .
   ```

2. **Run format checks:**
   ```bash
   go fmt ./...
   go vet ./...
   ```

3. **Test basic functionality:**
   ```bash
   ./yerd --version
   ./yerd php list
   ./yerd php status
   ```

4. **Test your specific changes:**
   - New command: test all flags and arguments
   - Bug fix: verify the issue is resolved
   - New feature: test edge cases and error handling

### Testing Checklist

- [ ] Code compiles without warnings
- [ ] Basic commands work (`--help`, `--version`)
- [ ] New features work as expected
- [ ] Error cases handled gracefully
- [ ] No regression in existing functionality

## 📤 Submitting Changes

### Before Submitting

1. **Update from upstream:**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run final checks:**
   ```bash
   go fmt ./...
   go vet ./...
   go build -o yerd .
   ```

3. **Update documentation** if needed

### Create Pull Request

1. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create PR on GitHub** with clear description

## 🔄 Pull Request Process

### PR Template

```markdown
## Description
Brief description of changes made and why.

## Type of Change
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring

## Testing
- [ ] Manual testing performed
- [ ] Tested on Linux
- [ ] Tested on macOS (if applicable)

## Checklist
- [ ] Code follows Go conventions
- [ ] Self-review completed
- [ ] Comments added for complex logic
- [ ] Documentation updated if needed
- [ ] No new warnings introduced
```

### Review Process

1. **Automated checks** run on PR creation
2. **Maintainer review** for code quality
3. **Discussion and feedback**
4. **Approval and merge**

## 👀 Code Review

### What We Look For

- **Code quality** and Go best practices
- **Clear logic** and readability
- **Error handling** completeness
- **Documentation** accuracy
- **Backwards compatibility**

### Addressing Feedback

- **Respond promptly** to reviewer comments
- **Ask questions** if unclear
- **Make requested changes** in new commits
- **Mark conversations** as resolved

## 🤝 Community Guidelines

### Communication

- **Be respectful** and constructive
- **Assume positive intent**
- **Help others** learn and improve
- **Use GitHub issues** for bugs and features
- **Use discussions** for questions

### Types of Contributions

We welcome:
- 🐛 **Bug fixes**
- ✨ **New features**
- 📚 **Documentation improvements**
- 🧪 **Test coverage**
- 🎨 **Code refactoring**
- 💡 **Ideas and suggestions**
- 🔧 **PHP extension additions**
- 🌍 **Multi-platform support**

## 🆘 Getting Help

- **Documentation:** [README.md](README.md) and [CLAUDE.md](CLAUDE.md)
- **Issues:** [GitHub Issues](https://github.com/LumoSolutions/yerd/issues)
- **Discussions:** [GitHub Discussions](https://github.com/LumoSolutions/yerd/discussions)

## 📄 License

By contributing to YERD, you agree that your contributions will be licensed under the [MIT License](LICENSE).

---

**Thank you for contributing to YERD!** 🎉

Every contribution helps make PHP development better for the entire community.