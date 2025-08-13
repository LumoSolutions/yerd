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

func (b *Builder) validateEnvironment() error {
	// Initialize dependency manager
	depMgr, err := dependencies.NewDependencyManager()
	if err != nil {
		return fmt.Errorf("failed to initialize dependency manager: %v", err)
	}
	
	fmt.Printf("Detected %s with %s package manager\n", depMgr.GetDistro(), depMgr.GetPackageManager())
	
	// Install build dependencies (gcc, make, autoconf, etc.)
	if err := depMgr.InstallBuildDependencies(); err != nil {
		return fmt.Errorf("failed to install build dependencies: %v", err)
	}
	
	// Install extension-specific dependencies
	if err := depMgr.InstallExtensionDependencies(b.Extensions); err != nil {
		return fmt.Errorf("failed to install extension dependencies: %v", err)
	}
	
	// Final validation to ensure all dependencies are available
	missing := depMgr.CheckSystemDependencies(b.Extensions)
	if len(missing) > 0 {
		return fmt.Errorf("dependencies still missing after installation: %s", strings.Join(missing, ", "))
	}
	
	fmt.Println("✓ All dependencies satisfied")
	return nil
}


func (b *Builder) downloadSource() error {
	if _, err := os.Stat(b.SourceDir); err == nil {
		return nil
	}
	
	fullVersion := b.getFullVersion()
	downloadURL := fmt.Sprintf("https://www.php.net/distributions/php-%s.tar.gz", fullVersion)
	
	// Download to temporary directory first (user has permissions there)
	tempDir := os.TempDir()
	tempArchivePath := filepath.Join(tempDir, fmt.Sprintf("php-%s.tar.gz", fullVersion))
	
	
	// Download to temp directory as user
	cmd := exec.Command("curl", "-L", "-o", tempArchivePath, downloadURL)
	if err := b.runCommand(cmd, "Downloading PHP source"); err != nil {
		return err
	}
	
	// Now extract and move as root
	// First ensure the target directory exists
	extractDir := filepath.Dir(b.SourceDir)
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("failed to create extract directory %s: %v", extractDir, err)
	}
	
	// Extract to temp directory first
	tempExtractDir := filepath.Join(tempDir, fmt.Sprintf("php-%s-extract", fullVersion))
	os.RemoveAll(tempExtractDir) // Clean up if exists
	if err := os.MkdirAll(tempExtractDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp extract directory: %v", err)
	}
	
	// Fix permissions for temp extract directory if running as root
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
	
	// Move extracted source to final location (as root)
	tempSourceDir := filepath.Join(tempExtractDir, fmt.Sprintf("php-%s", fullVersion))
	cmd = exec.Command("cp", "-r", tempSourceDir, b.SourceDir)
	cmd.SysProcAttr = nil // Run as root, don't drop privileges
	if err := b.runCommandAsRoot(cmd, "Moving source to final location"); err != nil {
		return err
	}
	
	// Fix ownership of the source directory so user can build in it
	if os.Geteuid() == 0 {
		sudoUser := os.Getenv("SUDO_USER")
		if sudoUser != "" {
			uid, gid, err := utils.GetSudoUserIDs(sudoUser)
			if err == nil {
				// Change ownership of the entire source directory to the user
				chownCmd := exec.Command("chown", "-R", fmt.Sprintf("%d:%d", uid, gid), b.SourceDir)
				if err := b.runCommandAsRoot(chownCmd, "Fixing source directory permissions"); err != nil {
					return fmt.Errorf("failed to fix source directory permissions: %v", err)
				}
			}
		}
	}
	
	// Clean up temp files
	os.Remove(tempArchivePath)
	os.RemoveAll(tempExtractDir)
	
	return nil
}

func (b *Builder) getFullVersion() string {
	// Try to fetch the latest version from PHP.net
	latestVersions, _, err := versions.FetchLatestVersions()
	if err == nil {
		if latestVersion, exists := latestVersions[b.Version]; exists {
			return latestVersion
		}
	}
	
	// Fallback: check if we have version info in config
	cfg, err := config.LoadConfig()
	if err == nil {
		if info, exists := cfg.InstalledPHP[b.Version]; exists && info.Version != "" {
			// If config has more than just major.minor (e.g. "8.4.11"), use it
			if len(strings.Split(info.Version, ".")) > 2 {
				return info.Version
			}
		}
	}
	
	// Last resort: try a reasonable default patch version
	// Most PHP versions start from .0, but let's be more conservative
	return b.Version + ".0"
}

func (b *Builder) configure() error {
	configFlags := b.getConfigureFlags()
	
	args := append([]string{"./configure"}, configFlags...)
	cmd := exec.Command("sh", "-c", strings.Join(args, " "))
	cmd.Dir = b.SourceDir
	
	return b.runCommand(cmd, "Configuring PHP build")
}

func (b *Builder) getConfigureFlags() []string {
	return php.GetConfigureFlagsForVersion(b.Version, b.Extensions)
}

func (b *Builder) compile() error {
	nproc := b.getProcessorCount()
	cmd := exec.Command("make", fmt.Sprintf("-j%d", nproc))
	cmd.Dir = b.SourceDir
	
	return b.runCommand(cmd, "Compiling PHP")
}

func (b *Builder) getProcessorCount() int {
	cmd := exec.Command("nproc")
	output, err := cmd.Output()
	if err != nil {
		return 4
	}
	
	nproc := strings.TrimSpace(string(output))
	if n := parseInt(nproc); n > 0 {
		return n
	}
	
	return 4
}

func parseInt(s string) int {
	var result int
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0
		}
		result = result*10 + int(r-'0')
	}
	return result
}

func (b *Builder) install() error {
	cmd := exec.Command("make", "install")
	cmd.Dir = b.SourceDir
	
	// Install needs root privileges to write to /opt/yerd/
	return b.runCommandAsRoot(cmd, "Installing PHP")
}

func (b *Builder) createSymlinks() error {
	phpBinary := filepath.Join(strings.TrimSuffix(b.InstallDir, "/"), "bin", "php")
	if _, err := os.Stat(phpBinary); err != nil {
		return fmt.Errorf("PHP binary not found at %s", phpBinary)
	}
	
	targetPath := php.GetBinaryPath(b.Version)
	
	if _, err := os.Lstat(targetPath); err == nil {
		if err := os.Remove(targetPath); err != nil {
			return fmt.Errorf("failed to remove existing symlink: %v", err)
		}
	}
	
	if err := os.Symlink(phpBinary, targetPath); err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}
	
	globalSymlink := "/usr/local/bin/php" + b.Version
	if err := os.Remove(globalSymlink); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing global symlink: %v", err)
	}
	
	if err := os.Symlink(targetPath, globalSymlink); err != nil {
		return fmt.Errorf("failed to create global symlink: %v", err)
	}
	
	return nil
}

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
	
	// Don't set user credentials - run as root
	
	if err := cmd.Run(); err != nil {
		logFile.WriteString(fmt.Sprintf("\nCommand failed with error: %v\n", err))
		return fmt.Errorf("%s failed (see %s for details): %v", description, b.LogPath, err)
	}
	
	logFile.WriteString(fmt.Sprintf("\n%s completed successfully\n", description))
	return nil
}

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

func (b *Builder) Cleanup() error {
	if b.SourceDir != "" {
		os.RemoveAll(b.SourceDir)
	}
	
	return nil
}

func (b *Builder) CleanupSuccess() error {
	if b.SourceDir != "" {
		os.RemoveAll(b.SourceDir)
	}
	
	if b.LogPath != "" {
		os.Remove(b.LogPath)
	}
	
	return nil
}