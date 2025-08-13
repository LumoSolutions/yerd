package web

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LumoSolutions/yerd/internal/dependencies"
	"github.com/LumoSolutions/yerd/internal/utils"
)

// WebManager handles web service installation and management
type WebManager struct {
	depManager *dependencies.DependencyManager
}

// NewWebManager creates a new web service manager
func NewWebManager() (*WebManager, error) {
	depManager, err := dependencies.NewDependencyManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dependency manager: %w", err)
	}

	return &WebManager{
		depManager: depManager,
	}, nil
}

// InstallService installs a web service
func (wm *WebManager) InstallService(serviceName string) error {
	config, exists := GetServiceConfig(serviceName)
	if !exists {
		return fmt.Errorf("unsupported service: %s", serviceName)
	}

	err := wm.checkPermissions()
	if err != nil {
		return fmt.Errorf("permission check failed: %w", err)
	}

	err = wm.ensureDirectories(config)
	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	err = wm.installDependencies(config)
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	err = wm.downloadAndBuildService(config)
	if err != nil {
		return fmt.Errorf("failed to download and build service: %w", err)
	}

	return nil
}

// checkPermissions verifies installation permissions
func (wm *WebManager) checkPermissions() error {
	return utils.CheckInstallPermissions()
}

// ensureDirectories creates necessary directories for the service
func (wm *WebManager) ensureDirectories(config *ServiceConfig) error {
	spinner := utils.NewLoadingSpinner(fmt.Sprintf("Creating directories for %s...", config.Name))
	spinner.Start()
	defer spinner.Stop("✓ Directories created")

	dirs := []string{
		config.InstallPath,
		GetServiceConfigPath(config.Name),
		GetServiceLogPath(config.Name),
		GetServiceRunPath(config.Name),
		GetServiceTempPath(config.Name),
		filepath.Join(config.InstallPath, "sbin"),
	}

	for _, dir := range dirs {
		err := utils.CreateDirectory(dir)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// installDependencies installs system dependencies for the service
func (wm *WebManager) installDependencies(config *ServiceConfig) error {
	if len(config.Dependencies) == 0 {
		return nil
	}

	spinner := utils.NewLoadingSpinner(fmt.Sprintf("Installing dependencies for %s...", config.Name))
	spinner.Start()
	defer spinner.Stop("✓ Dependencies installed")

	return wm.depManager.InstallExtensionDependencies(config.Dependencies)
}

// downloadAndBuildService downloads source and builds the service
func (wm *WebManager) downloadAndBuildService(config *ServiceConfig) error {
	buildDir := filepath.Join("/tmp", fmt.Sprintf("yerd-build-%s", config.Name))

	err := os.RemoveAll(buildDir)
	if err != nil {
		return fmt.Errorf("failed to clean build directory: %w", err)
	}

	err = utils.CreateDirectory(buildDir)
	if err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}
	defer os.RemoveAll(buildDir)

	err = wm.downloadSource(config, buildDir)
	if err != nil {
		return fmt.Errorf("failed to download source: %w", err)
	}

	err = wm.buildService(config, buildDir)
	if err != nil {
		return fmt.Errorf("failed to build service: %w", err)
	}

	return nil
}

// downloadSource downloads and extracts the service source
func (wm *WebManager) downloadSource(config *ServiceConfig, buildDir string) error {
	spinner := utils.NewLoadingSpinner(fmt.Sprintf("Downloading %s source...", config.Name))
	spinner.Start()
	defer spinner.Stop("✓ Source downloaded")

	tarballPath := filepath.Join(buildDir, fmt.Sprintf("%s.tar.gz", config.Name))

	_, err := utils.ExecuteCommand("wget", "-O", tarballPath, config.DownloadURL)
	if err != nil {
		_, err = utils.ExecuteCommand("curl", "-L", "-o", tarballPath, config.DownloadURL)
		if err != nil {
			return fmt.Errorf("failed to download with both wget and curl: %w", err)
		}
	}

	_, err = utils.ExecuteCommand("tar", "-xzf", tarballPath, "-C", buildDir, "--strip-components=1")
	if err != nil {
		return fmt.Errorf("failed to extract tarball: %w", err)
	}

	return nil
}

// buildService configures, compiles and installs the service
func (wm *WebManager) buildService(config *ServiceConfig, buildDir string) error {
	configureSpinner := utils.NewLoadingSpinner(fmt.Sprintf("Configuring %s build...", config.Name))
	configureSpinner.Start()

	_, err := utils.ExecuteCommand("sh", append([]string{"-c", "cd " + buildDir + " && ./configure"}, config.BuildFlags...)...)
	configureSpinner.Stop("✓ Build configured")
	if err != nil {
		return fmt.Errorf("configure failed: %w", err)
	}

	buildSpinner := utils.NewLoadingSpinner(fmt.Sprintf("Building %s...", config.Name))
	buildSpinner.Start()

	processorCount := utils.GetProcessorCount()
	_, err = utils.ExecuteCommand("sh", "-c", fmt.Sprintf("cd %s && make -j%d", buildDir, processorCount))
	buildSpinner.Stop("✓ Build completed")
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	installSpinner := utils.NewLoadingSpinner(fmt.Sprintf("Installing %s...", config.Name))
	installSpinner.Start()

	_, err = utils.ExecuteCommand("sh", "-c", fmt.Sprintf("cd %s && make install", buildDir))
	installSpinner.Stop("✓ Installation completed")
	if err != nil {
		return fmt.Errorf("install failed: %w", err)
	}

	return nil
}
