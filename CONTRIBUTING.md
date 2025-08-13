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
- [Making Changes](#making-changes)
- [Code Standards](#code-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Pull Request Process](#pull-request-process)
- [Code Review](#code-review)
- [Community Guidelines](#community-guidelines)

## 🏁 Getting Started

### Prerequisites

- **Go 1.21+** installed
- **Git** for version control
- **Linux environment** (Ubuntu, Debian, Arch, etc.)
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
   ./yerd php --help
   ```

4. **Run tests:**
   ```bash
   go test ./...
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

## ✏️ Making Changes

### Branch Naming

Use descriptive branch names that indicate the type of change:

- **Features:** `feature/add-nginx-support`
- **Bug fixes:** `fix/extension-detection-error`
- **Documentation:** `docs/update-readme-examples`
- **Refactoring:** `refactor/simplify-config-loading`

### Commit Messages

Write clear, descriptive commit messages:

```
Add support for custom PHP configure flags

- Allow users to specify custom configure flags via config file
- Add validation for configure flag syntax
- Update documentation with examples
- Add tests for flag validation

Closes #123
```

**Format:**
- Use imperative mood ("Add", "Fix", "Update")
- First line: concise summary (50 chars or less)
- Blank line, then detailed explanation if needed
- Reference issues with "Closes #123" or "Fixes #456"

### File Organization

YERD follows a structured organization:

```
yerd/
├── cmd/                    # CLI commands
│   ├── php/               # PHP-specific commands
│   └── root.go            # Root command
├── internal/              # Internal packages
│   ├── builder/           # PHP building logic
│   ├── config/            # Configuration management
│   ├── installer/         # Installation logic
│   ├── utils/             # Utility functions
│   └── versions/          # Version management
├── pkg/                   # Public packages
│   ├── constants/         # Shared constants
│   ├── extensions/        # PHP extensions
│   └── php/               # PHP version definitions
└── scripts/               # Build and release scripts
```

## 📏 Code Standards

### Go Style

- Follow **Go conventions** and **gofmt** formatting
- Use **meaningful variable names**
- Add **comments for exported functions**
- Keep functions **focused and small**
- Use **early returns** to reduce nesting

### Documentation

- All **exported functions** must have comments
- Use **Go doc comment format:**
  ```go
  // FunctionName does something specific and returns an error if it fails.
  // parameter: Description of what this parameter does.
  func FunctionName(parameter string) error {
      // implementation
  }
  ```

### Error Handling

- Always **handle errors appropriately**
- Provide **context in error messages**
- Use **fmt.Errorf** for error wrapping
- Don't ignore errors with `_`

### Example:
```go
// InstallPHP installs a PHP version with specified extensions.
// version: PHP version to install, extensions: List of extensions to include.
func InstallPHP(version string, extensions []string) error {
    if version == "" {
        return fmt.Errorf("version cannot be empty")
    }
    
    cfg, err := config.LoadConfig()
    if err != nil {
        return fmt.Errorf("failed to load config: %v", err)
    }
    
    // ... rest of implementation
}
```

## 🧪 Testing

### Manual Testing

Before submitting a PR, test your changes:

1. **Build successfully:**
   ```bash
   go build -o yerd .
   ```

2. **Test basic functionality:**
   ```bash
   ./yerd --help
   ./yerd php list
   ./yerd status
   ```

3. **Test your specific changes:**
   - If you added a new command, test all its options
   - If you fixed a bug, verify the fix works
   - If you added a feature, test edge cases

### Automated Tests

Run the test suite:
```bash
go test ./...
go vet ./...
```

If tests fail, fix them before submitting your PR.

## 📤 Submitting Changes

### Before Submitting

1. **Rebase** on the latest upstream main:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Squash commits** if you have multiple small commits:
   ```bash
   git rebase -i HEAD~3  # Interactive rebase for last 3 commits
   ```

3. **Test one final time:**
   ```bash
   go build -o yerd . && ./yerd --help
   ```

### Create Pull Request

1. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create PR on GitHub:**
   - Go to your fork on GitHub
   - Click "Compare & pull request"
   - Fill out the PR template

## 🔄 Pull Request Process

### PR Template

When creating a PR, include:

```markdown
## Description
Brief description of changes made.

## Type of Change
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)  
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## How Has This Been Tested?
- [ ] Manual testing performed
- [ ] Unit tests added/updated
- [ ] Integration tests pass

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex logic
- [ ] Documentation updated if needed
- [ ] No new warnings introduced
```

### PR Requirements

- ✅ **Clear description** of changes
- ✅ **Tests passing** (manual and automated)
- ✅ **Documentation updated** if needed
- ✅ **No merge conflicts** with main branch
- ✅ **Follows code standards**

## 👀 Code Review

### Review Process

1. **Automated checks** run first (builds, tests)
2. **Maintainer review** for code quality and design
3. **Feedback addressed** through discussion
4. **Approval and merge** when ready

### Addressing Feedback

- **Be responsive** to reviewer comments
- **Ask questions** if feedback is unclear
- **Make requested changes** in new commits
- **Explain your reasoning** if you disagree

## 🤝 Community Guidelines

### Be Respectful

- **Constructive feedback** only
- **Assume positive intent** in discussions
- **Help others learn** and improve

### Communication

- **Use GitHub issues** for bug reports and feature requests
- **Use discussions** for general questions
- **Be patient** - maintainers are volunteers

### Contribution Types

We welcome various types of contributions:

- 🐛 **Bug fixes**
- ✨ **New features**
- 📚 **Documentation improvements**
- 🧪 **Test coverage**
- 🎨 **Code refactoring**
- 🌍 **Translations**
- 💡 **Ideas and suggestions**

## 🆘 Getting Help

- **Documentation:** Check [README.md](README.md) and [CLAUDE.md](CLAUDE.md)
- **Issues:** Browse [existing issues](https://github.com/LumoSolutions/yerd/issues)
- **Discussions:** Use [GitHub Discussions](https://github.com/LumoSolutions/yerd/discussions)
- **Contact:** Reach out to the maintainers

## 📄 License

By contributing to YERD, you agree that your contributions will be licensed under the [MIT License](LICENSE).

---

**Thank you for contributing to YERD!** 🎉

Every contribution, no matter how small, helps make YERD better for the entire PHP community.