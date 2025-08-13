package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/dependencies"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/versions"
	"github.com/LumoSolutions/yerd/pkg/php"
)

type Builder struct {
	Version    string
	Extensions []string
	SourceDir  string
	InstallDir string
	LogPath    string
}

// NewBuilder creates a new PHP builder instance with source and install directories.
// version: PHP version to build, extensions: List of extensions to include. Returns configured Builder.
func NewBuilder(version string, extensions []string) *Builder {
	sourceDir := filepath.Join(utils.YerdPHPDir, "src", "php-"+version)
	installDir := php.GetInstallPath(version)
	logDir := filepath.Join(os.TempDir(), "yerd-build")
	os.MkdirAll(logDir, 0755)
	logPath := filepath.Join(logDir, fmt.Sprintf("php-%s-build.log", version))

	return &Builder{
		Version:    version,
		Extensions: extensions,
		SourceDir:  sourceDir,
		InstallDir: installDir,
		LogPath:    logPath,
	}
}

// RebuildPHP performs complete PHP compilation from source with dependencies and configuration.
// Returns error if any build step fails.
func (b *Builder) RebuildPHP() error {
	if err := b.validateEnvironment(); err != nil {
		return fmt.Errorf("environment validation failed: %v", err)
	}

	if err := b.downloadSource(); err != nil {
		return fmt.Errorf("source download failed: %v", err)
	}

	if err := b.configure(); err != nil {
		return fmt.Errorf("configure failed: %v", err)
	}

	if err := b.compile(); err != nil {
		return fmt.Errorf("compilation failed: %v", err)
	}

	if err := b.install(); err != nil {
		return fmt.Errorf("installation failed: %v", err)
	}

	if err := b.createSymlinks(); err != nil {
		return fmt.Errorf("symlink creation failed: %v", err)
	}

	return nil
}

// validateEnvironment checks and installs required build dependencies and extensions.
// Returns error if dependencies cannot be satisfied.
func (b *Builder) validateEnvironment() error {
	depMgr, err := dependencies.NewDependencyManager()
	if err != nil {
		return fmt.Errorf("failed to initialize dependency manager: %v", err)
	}

	utils.PrintInfo("Detected %s with %s package manager", depMgr.GetDistro(), depMgr.GetPackageManager())

	if err := depMgr.InstallBuildDependencies(); err != nil {
		return fmt.Errorf("failed to install build dependencies: %v", err)
	}

	if err := depMgr.InstallExtensionDependencies(b.Extensions); err != nil {
		return fmt.Errorf("failed to install extension dependencies: %v", err)
	}

	missing := depMgr.CheckSystemDependencies(b.Extensions)
	if len(missing) > 0 {
		return fmt.Errorf("dependencies still missing after installation: %s", strings.Join(missing, ", "))
	}

	utils.PrintSuccess("All dependencies satisfied")
	return nil
}

// downloadSource downloads and extracts PHP source code from official distribution.
// Returns error if download or extraction fails.
func (b *Builder) downloadSource() error {
	if _, err := os.Stat(b.SourceDir); err == nil {
		return nil
	}

	fullVersion := b.getFullVersion()
	downloadURL := fmt.Sprintf("https://www.php.net/distributions/php-%s.tar.gz", fullVersion)

	tempDir := os.TempDir()
	tempArchivePath := filepath.Join(tempDir, fmt.Sprintf("php-%s.tar.gz", fullVersion))

	cmd := exec.Command("curl", "-L", "-o", tempArchivePath, downloadURL)
	if err := b.runCommand(cmd, "Downloading PHP source"); err != nil {
		return err
	}

	extractDir := filepath.Dir(b.SourceDir)
	if err := utils.CreateDirectory(extractDir); err != nil {
		return err
	}

	tempExtractDir := filepath.Join(tempDir, fmt.Sprintf("php-%s-extract", fullVersion))
	os.RemoveAll(tempExtractDir)
	if err := utils.CreateDirectory(tempExtractDir); err != nil {
		return err
	}

	if os.Geteuid() == 0 {
		sudoUser := os.Getenv("SUDO_USER")
		if sudoUser != "" {
			uid, gid, err := utils.GetSudoUserIDs(sudoUser)
			if err == nil {
				os.Chown(tempExtractDir, uid, gid)
			}
		}
	}

	cmd = exec.Command("tar", "-xzf", tempArchivePath, "-C", tempExtractDir)
	if err := b.runCommand(cmd, "Extracting PHP source"); err != nil {
		return err
	}

	tempSourceDir := filepath.Join(tempExtractDir, fmt.Sprintf("php-%s", fullVersion))
	cmd = exec.Command("cp", "-r", tempSourceDir, b.SourceDir)
	cmd.SysProcAttr = nil
	if err := b.runCommandAsRoot(cmd, "Moving source to final location"); err != nil {
		return err
	}

	if os.Geteuid() == 0 {
		sudoUser := os.Getenv("SUDO_USER")
		if sudoUser != "" {
			uid, gid, err := utils.GetSudoUserIDs(sudoUser)
			if err == nil {
				chownCmd := exec.Command("chown", "-R", fmt.Sprintf("%d:%d", uid, gid), b.SourceDir)
				if err := b.runCommandAsRoot(chownCmd, "Fixing source directory permissions"); err != nil {
					return fmt.Errorf("failed to fix source directory permissions: %v", err)
				}
			}
		}
	}

	os.Remove(tempArchivePath)
	os.RemoveAll(tempExtractDir)

	return nil
}

// getFullVersion resolves short version string to full version number from remote or config.
// Returns full version string with patch number.
func (b *Builder) getFullVersion() string {
	latestVersions, _, err := versions.FetchLatestVersions()
	if err == nil {
		if latestVersion, exists := latestVersions[b.Version]; exists {
			return latestVersion
		}
	}

	cfg, err := config.LoadConfig()
	if err == nil {
		if info, exists := cfg.InstalledPHP[b.Version]; exists && info.Version != "" {
			if len(strings.Split(info.Version, ".")) > 2 {
				return info.Version
			}
		}
	}

	return b.Version + ".0"
}

// configure runs PHP's configure script with appropriate flags for version and extensions.
// Returns error if configuration fails.
func (b *Builder) configure() error {
	configFlags := b.getConfigureFlags()

	args := append([]string{"./configure"}, configFlags...)
	cmd := exec.Command("sh", "-c", strings.Join(args, " "))
	cmd.Dir = b.SourceDir

	return b.runCommand(cmd, "Configuring PHP build")
}

// getConfigureFlags returns configure script flags based on PHP version and extensions.
// Returns slice of configure flag strings.
func (b *Builder) getConfigureFlags() []string {
	return php.GetConfigureFlagsForVersion(b.Version, b.Extensions)
}

// compile runs make with optimal parallel job count to build PHP from source.
// Returns error if compilation fails.
func (b *Builder) compile() error {
	nproc := b.getProcessorCount()
	cmd := exec.Command("make", fmt.Sprintf("-j%d", nproc))
	cmd.Dir = b.SourceDir

	return b.runCommand(cmd, "Compiling PHP")
}

// getProcessorCount detects system CPU count for parallel compilation jobs.
// Returns processor count or defaults to 4 if detection fails.
func (b *Builder) getProcessorCount() int {
	return utils.GetProcessorCount()
}

// install runs make install to install compiled PHP to target directory with root privileges.
// Returns error if installation fails.
func (b *Builder) install() error {
	cmd := exec.Command("make", "install")
	cmd.Dir = b.SourceDir

	return b.runCommandAsRoot(cmd, "Installing PHP")
}

// createSymlinks creates version-specific and global symlinks for the PHP binary.
// Returns error if symlink creation fails.
func (b *Builder) createSymlinks() error {
	phpBinary := filepath.Join(strings.TrimSuffix(b.InstallDir, "/"), "bin", "php")
	if _, err := os.Stat(phpBinary); err != nil {
		return fmt.Errorf("PHP binary not found at %s", phpBinary)
	}

	targetPath := php.GetBinaryPath(b.Version)
	if err := utils.CreateSymlink(phpBinary, targetPath); err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}

	globalSymlink := "/usr/local/bin/php" + b.Version
	if err := utils.CreateSymlink(targetPath, globalSymlink); err != nil {
		return fmt.Errorf("failed to create global symlink: %v", err)
	}

	return nil
}

// runCommand executes a command with logging and user privilege handling.
// cmd: Command to execute, description: Operation description. Returns error if command fails.
func (b *Builder) runCommand(cmd *exec.Cmd, description string) error {
	logFile, err := os.OpenFile(b.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	logFile.WriteString(fmt.Sprintf("\n=== %s ===\n", description))
	logFile.WriteString(fmt.Sprintf("Command: %s\n", strings.Join(cmd.Args, " ")))
	logFile.WriteString(fmt.Sprintf("Directory: %s\n\n", cmd.Dir))

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if os.Geteuid() == 0 {
		if err := b.setUserCredentials(cmd); err != nil {
			logFile.WriteString(fmt.Sprintf("Warning: failed to set user credentials: %v\n", err))
		}
	}

	if err := cmd.Run(); err != nil {
		logFile.WriteString(fmt.Sprintf("\nCommand failed with error: %v\n", err))
		return fmt.Errorf("%s failed (see %s for details): %v", description, b.LogPath, err)
	}

	logFile.WriteString(fmt.Sprintf("\n%s completed successfully\n", description))
	return nil
}

// runCommandAsRoot executes a command with root privileges and full logging.
// cmd: Command to execute, description: Operation description. Returns error if command fails.
func (b *Builder) runCommandAsRoot(cmd *exec.Cmd, description string) error {
	logFile, err := os.OpenFile(b.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	logFile.WriteString(fmt.Sprintf("\n=== %s (as root) ===\n", description))
	logFile.WriteString(fmt.Sprintf("Command: %s\n", strings.Join(cmd.Args, " ")))
	logFile.WriteString(fmt.Sprintf("Directory: %s\n\n", cmd.Dir))

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Run(); err != nil {
		logFile.WriteString(fmt.Sprintf("\nCommand failed with error: %v\n", err))
		return fmt.Errorf("%s failed (see %s for details): %v", description, b.LogPath, err)
	}

	logFile.WriteString(fmt.Sprintf("\n%s completed successfully\n", description))
	return nil
}

// setUserCredentials configures command to run as original SUDO user instead of root.
// cmd: Command to configure. Returns error if user context cannot be determined.
func (b *Builder) setUserCredentials(cmd *exec.Cmd) error {
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser == "" {
		return nil
	}

	uid, gid, err := utils.GetSudoUserIDs(sudoUser)
	if err != nil {
		return err
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}

	return nil
}

// Cleanup removes source directory after failed build to free disk space.
// Returns error if cleanup fails.
func (b *Builder) Cleanup() error {
	if b.SourceDir != "" {
		os.RemoveAll(b.SourceDir)
	}

	return nil
}

// CleanupSuccess removes source directory and logs after successful build.
// Returns error if cleanup fails.
func (b *Builder) CleanupSuccess() error {
	if b.SourceDir != "" {
		os.RemoveAll(b.SourceDir)
	}

	if b.LogPath != "" {
		os.Remove(b.LogPath)
	}

	return nil
}
