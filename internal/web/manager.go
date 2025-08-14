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

	opts := utils.DefaultDownloadOptions()
	_, err := utils.DownloadAndExtractTarGz(config.DownloadURL, buildDir, opts)
	if err != nil {
		return fmt.Errorf("failed to download and extract source: %w", err)
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

// StartService starts a web service (nginx or dnsmasq)
func (wm *WebManager) StartService(serviceName string) error {
	_, exists := GetServiceConfig(serviceName)
	if !exists {
		return fmt.Errorf("unsupported service: %s", serviceName)
	}

	if !IsServiceInstalled(serviceName) {
		return fmt.Errorf("service %s is not installed. Run 'yerd web install' first", serviceName)
	}

	err := wm.checkPermissions()
	if err != nil {
		return fmt.Errorf("permission check failed: %w", err)
	}

	binaryPath := GetServiceBinaryPath(serviceName)
	
	switch serviceName {
	case "nginx":
		return wm.startNginx(binaryPath)
	case "dnsmasq":
		return wm.startDnsmasq(binaryPath)
	default:
		return fmt.Errorf("starting %s is not supported yet", serviceName)
	}
}

// StopService stops a web service (nginx or dnsmasq)
func (wm *WebManager) StopService(serviceName string) error {
	_, exists := GetServiceConfig(serviceName)
	if !exists {
		return fmt.Errorf("unsupported service: %s", serviceName)
	}

	if !IsServiceInstalled(serviceName) {
		return fmt.Errorf("service %s is not installed", serviceName)
	}

	err := wm.checkPermissions()
	if err != nil {
		return fmt.Errorf("permission check failed: %w", err)
	}

	switch serviceName {
	case "nginx":
		return wm.stopNginx()
	case "dnsmasq":
		return wm.stopDnsmasq()
	default:
		return fmt.Errorf("stopping %s is not supported yet", serviceName)
	}
}

// StartAllServices starts nginx and dnsmasq
func (wm *WebManager) StartAllServices() error {
	services := []string{"nginx", "dnsmasq"}
	var failed []string

	for _, service := range services {
		if err := wm.StartService(service); err != nil {
			failed = append(failed, fmt.Sprintf("%s: %v", service, err))
		}
	}

	if len(failed) > 0 {
		return fmt.Errorf("failed to start some services: %v", failed)
	}

	return nil
}

// StopAllServices stops nginx and dnsmasq
func (wm *WebManager) StopAllServices() error {
	services := []string{"nginx", "dnsmasq"}
	var failed []string

	for _, service := range services {
		if err := wm.StopService(service); err != nil {
			failed = append(failed, fmt.Sprintf("%s: %v", service, err))
		}
	}

	if len(failed) > 0 {
		return fmt.Errorf("failed to stop some services: %v", failed)
	}

	return nil
}

// startNginx starts the nginx service
func (wm *WebManager) startNginx(binaryPath string) error {
	configPath := filepath.Join(GetServiceConfigPath("nginx"), "nginx.conf")
	
	if !utils.FileExists(configPath) {
		return fmt.Errorf("nginx configuration not found at %s", configPath)
	}

	// Test configuration first
	_, err := utils.ExecuteCommand(binaryPath, "-t", "-c", configPath)
	if err != nil {
		return fmt.Errorf("nginx configuration test failed: %w", err)
	}

	spinner := utils.NewLoadingSpinner("Starting nginx...")
	spinner.Start()
	defer spinner.Stop("✓ nginx started")

	_, err = utils.ExecuteCommand(binaryPath, "-c", configPath)
	if err != nil {
		return fmt.Errorf("failed to start nginx: %w", err)
	}

	return nil
}

// stopNginx stops the nginx service
func (wm *WebManager) stopNginx() error {
	pidPath := filepath.Join(GetServiceRunPath("nginx"), "nginx.pid")
	
	if !utils.FileExists(pidPath) {
		return fmt.Errorf("nginx is not running (PID file not found)")
	}

	spinner := utils.NewLoadingSpinner("Stopping nginx...")
	spinner.Start()
	defer spinner.Stop("✓ nginx stopped")

	binaryPath := GetServiceBinaryPath("nginx")
	configPath := filepath.Join(GetServiceConfigPath("nginx"), "nginx.conf")
	
	_, err := utils.ExecuteCommand(binaryPath, "-c", configPath, "-s", "quit")
	if err != nil {
		return fmt.Errorf("failed to stop nginx: %w", err)
	}

	return nil
}

// startDnsmasq starts the dnsmasq service  
func (wm *WebManager) startDnsmasq(binaryPath string) error {
	configPath := filepath.Join(GetServiceConfigPath("dnsmasq"), "dnsmasq.conf")
	
	if !utils.FileExists(configPath) {
		return fmt.Errorf("dnsmasq configuration not found at %s. Run 'sudo yerd web install -f' to recreate configurations", configPath)
	}

	// Test configuration first
	_, err := utils.ExecuteCommand(binaryPath, "--test", "-C", configPath)
	if err != nil {
		return fmt.Errorf("dnsmasq configuration test failed: %w", err)
	}

	spinner := utils.NewLoadingSpinner("Starting dnsmasq...")
	spinner.Start()
	defer spinner.Stop("✓ dnsmasq started")

	_, err = utils.ExecuteCommand(binaryPath, "-C", configPath)
	if err != nil {
		return fmt.Errorf("failed to start dnsmasq: %w", err)
	}

	return nil
}

// stopDnsmasq stops the dnsmasq service
func (wm *WebManager) stopDnsmasq() error {
	pidPath := filepath.Join(GetServiceRunPath("dnsmasq"), "dnsmasq.pid")
	
	if !utils.FileExists(pidPath) {
		return fmt.Errorf("dnsmasq is not running (PID file not found)")
	}

	spinner := utils.NewLoadingSpinner("Stopping dnsmasq...")
	spinner.Start()
	defer spinner.Stop("✓ dnsmasq stopped")

	_, err := utils.ExecuteCommand("pkill", "-f", "dnsmasq")
	if err != nil {
		return fmt.Errorf("failed to stop dnsmasq: %w", err)
	}

	return nil
}
