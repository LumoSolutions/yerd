package web

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LumoSolutions/yerd/internal/dependencies"
	"github.com/LumoSolutions/yerd/internal/utils"
)

// WebInstaller handles installation of web services
type WebInstaller struct {
	service    string
	config     *ServiceConfig
	depManager *dependencies.DependencyManager
	logger     *utils.Logger
}

// NewWebInstaller creates a new web service installer
func NewWebInstaller(service string) (*WebInstaller, error) {
	config, exists := GetServiceConfig(service)
	if !exists {
		return nil, fmt.Errorf("service '%s' is not supported", service)
	}

	depManager, err := dependencies.NewDependencyManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dependency manager: %v", err)
	}
	
	// Set quiet mode to prevent interference with our spinners
	depManager.SetQuiet(true)

	logger, err := utils.NewLogger(fmt.Sprintf("web-%s", service))
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	return &WebInstaller{
		service:    service,
		config:     config,
		depManager: depManager,
		logger:     logger,
	}, nil
}

// Install performs the complete installation process
func (wi *WebInstaller) Install() error {
	var installSuccess bool
	defer wi.cleanupLogger(&installSuccess)

	fmt.Printf("Installing %s %s...\n", wi.config.Name, wi.config.Version)

	// Check permissions
	if err := utils.CheckInstallPermissions(); err != nil {
		return fmt.Errorf("insufficient permissions: %v", err)
	}

	// Create directories
	if err := wi.createDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	// Install dependencies
	if err := wi.installDependencies(); err != nil {
		return err
	}

	// Download source
	sourceDir, err := wi.downloadSource()
	if err != nil {
		return fmt.Errorf("failed to download source: %v", err)
	}
	defer os.RemoveAll(sourceDir)

	// Compile and install
	if err := wi.compileAndInstall(sourceDir); err != nil {
		return fmt.Errorf("failed to compile and install: %v", err)
	}

	// Create configuration if needed
	if err := wi.createConfiguration(); err != nil {
		return fmt.Errorf("failed to create configuration: %v", err)
	}

	// Mark installation as successful
	installSuccess = true

	fmt.Println()
	utils.PrintSuccess("Successfully installed %s %s", wi.config.Name, wi.config.Version)
	wi.printBinaryLocation()

	return nil
}

// installDependencies handles both build and service dependency installation
func (wi *WebInstaller) installDependencies() error {
	if len(wi.config.Dependencies) == 0 {
		return nil
	}

	spinner := utils.NewLoadingSpinner(fmt.Sprintf("Installing dependencies for %s...", wi.config.Name))
	spinner.Start()
	defer spinner.Stop("✓ Dependencies installed")

	return wi.depManager.InstallExtensionDependencies(wi.config.Dependencies)
}

// executeWithSpinner runs a command with a loading spinner
func (wi *WebInstaller) executeWithSpinner(message string, successMsg string, cmd string, args ...string) error {
	spinner := utils.NewLoadingSpinner(message)
	spinner.Start()
	
	if _, err := utils.ExecuteCommandWithLogging(wi.logger, cmd, args...); err != nil {
		spinner.Stop("✗ " + message + " failed")
		return fmt.Errorf("%s failed: %v", message, err)
	}
	
	spinner.Stop("✓ " + successMsg)
	return nil
}

// printBinaryLocation displays the binary location with consistent formatting
func (wi *WebInstaller) printBinaryLocation() {
	fmt.Printf("Binary location: %s\n\n", GetServiceBinaryPath(wi.service))
}

// InstallWithReplace performs installation with optional replacement of existing service
func (wi *WebInstaller) InstallWithReplace(replaceExisting bool) error {
	if replaceExisting {
		// Build to a temporary location first
		return wi.installWithBackup()
	} else {
		// Normal installation
		return wi.Install()
	}
}

// installWithBackup builds to temp location, then replaces existing on success
func (wi *WebInstaller) installWithBackup() error {
	var installSuccess bool
	defer wi.cleanupLogger(&installSuccess)

	fmt.Printf("Installing %s %s...\n", wi.config.Name, wi.config.Version)

	// Check permissions
	if err := utils.CheckInstallPermissions(); err != nil {
		return fmt.Errorf("insufficient permissions: %v", err)
	}

	// Store original install path
	originalInstallPath := wi.config.InstallPath

	// Install dependencies
	if err := wi.installDependencies(); err != nil {
		return err
	}

	// Download source
	sourceDir, err := wi.downloadSource()
	if err != nil {
		return fmt.Errorf("failed to download source: %v", err)
	}
	defer os.RemoveAll(sourceDir)

	// Remove the old installation first (before building)
	if utils.FileExists(originalInstallPath) {
		if err := utils.RemoveDirectory(originalInstallPath); err != nil {
			return fmt.Errorf("failed to remove existing installation: %v", err)
		}
	}

	// Ensure all required directories exist in final location
	if err := wi.ensureFinalDirectories(); err != nil {
		return fmt.Errorf("failed to create final directories: %v", err)
	}

	// Build and install directly to final location
	if err := wi.compileAndInstall(sourceDir); err != nil {
		return fmt.Errorf("failed to compile and install: %v", err)
	}

	// Create configuration in final location
	// Reset config to use original path for config creation
	wi.config.InstallPath = originalInstallPath
	if err := wi.createConfiguration(); err != nil {
		return fmt.Errorf("failed to create configuration: %v", err)
	}

	fmt.Printf("\n🎉 Build and installation successful!\n")

	// Mark installation as successful
	installSuccess = true

	fmt.Println()
	utils.PrintSuccess("Successfully replaced %s %s", wi.config.Name, wi.config.Version)
	wi.printBinaryLocation()

	return nil
}

// ensureFinalDirectories ensures all required directories exist in the final installation
func (wi *WebInstaller) ensureFinalDirectories() error {
	// Use the original service name to get the final paths
	finalDirs := []string{
		"/opt/yerd/web/" + wi.service + "/conf",
		"/opt/yerd/web/" + wi.service + "/logs", 
		"/opt/yerd/web/" + wi.service + "/run",
	}

	// Add temp path for nginx
	if wi.service == "nginx" {
		finalDirs = append(finalDirs, "/opt/yerd/web/nginx/temp")
	}

	for _, dir := range finalDirs {
		if err := utils.CreateDirectory(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
		wi.logger.WriteLog("Ensured directory exists: %s", dir)
	}

	return nil
}

// createDirectories creates all necessary directories for the service
func (wi *WebInstaller) createDirectories() error {
	dirs := []string{
		wi.config.InstallPath,
		GetServiceConfigPath(wi.service),
		GetServiceLogPath(wi.service),
		GetServiceRunPath(wi.service),
	}

	// Add temp path for nginx
	if wi.service == "nginx" {
		dirs = append(dirs, GetServiceTempPath(wi.service))
	}

	for _, dir := range dirs {
		if err := utils.CreateDirectory(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
		wi.logger.WriteLog("Created directory: %s", dir)
	}

	return nil
}

// downloadSource downloads and extracts the service source code
func (wi *WebInstaller) downloadSource() (string, error) {
	spinner := utils.NewLoadingSpinner("Downloading source code...")
	spinner.Start()

	// Create temporary directory
	tempDir := filepath.Join("/tmp", fmt.Sprintf("yerd-web-%s", wi.service))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		spinner.Stop("✗ Failed to create temp directory")
		return "", err
	}

	// Download file
	filename := filepath.Base(wi.config.DownloadURL)
	downloadPath := filepath.Join(tempDir, filename)

	// Try wget first, then curl
	var downloadCmd string
	var downloadArgs []string

	if _, err := utils.ExecuteCommand("which", "wget"); err == nil {
		downloadCmd = "wget"
		downloadArgs = []string{"-O", downloadPath, wi.config.DownloadURL}
	} else if _, err := utils.ExecuteCommand("which", "curl"); err == nil {
		downloadCmd = "curl"
		downloadArgs = []string{"-L", "-o", downloadPath, wi.config.DownloadURL}
	} else {
		spinner.Stop("✗ Neither wget nor curl found")
		return "", fmt.Errorf("neither wget nor curl is available for downloading")
	}

	wi.logger.WriteLog("Downloading with %s: %s", downloadCmd, wi.config.DownloadURL)
	if _, err := utils.ExecuteCommandWithLogging(wi.logger, downloadCmd, downloadArgs...); err != nil {
		spinner.Stop("✗ Download failed")
		return "", fmt.Errorf("download failed: %v", err)
	}

	// Extract archive
	wi.logger.WriteLog("Extracting archive: %s", downloadPath)
	if _, err := utils.ExecuteCommandWithLogging(wi.logger, "tar", "xzf", downloadPath, "-C", tempDir); err != nil {
		spinner.Stop("✗ Extraction failed")
		return "", fmt.Errorf("extraction failed: %v", err)
	}

	// Find extracted directory
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		spinner.Stop("✗ Failed to read temp directory")
		return "", err
	}

	var sourceDir string
	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(entry.Name(), wi.service) {
			sourceDir = filepath.Join(tempDir, entry.Name())
			break
		}
	}

	if sourceDir == "" {
		spinner.Stop("✗ Source directory not found")
		return "", fmt.Errorf("could not find extracted source directory")
	}

	spinner.Stop("✓ Source downloaded and extracted")
	wi.logger.WriteLog("Source extracted to: %s", sourceDir)
	return sourceDir, nil
}

// compileAndInstall compiles and installs the service
func (wi *WebInstaller) compileAndInstall(sourceDir string) error {
	// Change to source directory
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(sourceDir); err != nil {
		return err
	}

	wi.logger.WriteLog("Changed to source directory: %s", sourceDir)

	// Compile based on service type
	switch wi.service {
	case "nginx":
		return wi.compileNginx()
	case "dnsmasq":
		return wi.compileDnsmasq()
	default:
		return fmt.Errorf("unknown service type: %s", wi.service)
	}
}

// compileNginx compiles nginx with configure script
func (wi *WebInstaller) compileNginx() error {
	// Check if configure script exists
	if !utils.FileExists("./configure") {
		return fmt.Errorf("configure script not found in source directory")
	}

	// Configure
	wi.logger.WriteLog("Configuring nginx with flags: %v", wi.config.BuildFlags)
	if err := wi.executeWithSpinner("Configuring nginx...", "Configuration completed", "./configure", wi.config.BuildFlags...); err != nil {
		return err
	}

	// Build
	wi.logger.WriteLog("Building nginx...")
	if err := wi.executeWithSpinner("Building nginx...", "Build completed", "make", fmt.Sprintf("-j%d", utils.GetProcessorCount())); err != nil {
		return err
	}

	// Install
	wi.logger.WriteLog("Installing nginx...")
	if err := wi.executeWithSpinner("Installing nginx...", "🎉 Nginx compiled and installed", "make", "install"); err != nil {
		return err
	}
	return nil
}

// compileDnsmasq compiles dnsmasq with make
func (wi *WebInstaller) compileDnsmasq() error {
	// Build with custom flags
	makeArgs := append([]string{fmt.Sprintf("-j%d", utils.GetProcessorCount())}, wi.config.BuildFlags...)
	wi.logger.WriteLog("Building dnsmasq with flags: %v", makeArgs)
	if err := wi.executeWithSpinner("Building dnsmasq...", "Build completed", "make", makeArgs...); err != nil {
		return err
	}

	// Install with custom flags
	installArgs := append([]string{"install"}, wi.config.BuildFlags...)
	wi.logger.WriteLog("Installing dnsmasq with flags: %v", installArgs)
	if err := wi.executeWithSpinner("Installing dnsmasq...", "🎉 Dnsmasq compiled and installed", "make", installArgs...); err != nil {
		return err
	}
	return nil
}

// createConfiguration creates basic configuration files
func (wi *WebInstaller) createConfiguration() error {
	switch wi.service {
	case "nginx":
		return wi.createNginxConfig()
	case "dnsmasq":
		return wi.createDnsmasqConfig()
	default:
		return nil // No default config needed
	}
}

// createNginxConfig creates a basic nginx configuration
func (wi *WebInstaller) createNginxConfig() error {
	configDir := "/opt/yerd/web/nginx/conf"
	configPath := filepath.Join(configDir, "nginx.conf")
	
	// Ensure config directory exists
	if err := utils.CreateDirectory(configDir); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	
	if utils.FileExists(configPath) {
		return nil // Don't overwrite existing config
	}

	config := fmt.Sprintf(`user nobody;
worker_processes auto;
pid %s/nginx.pid;

events {
    worker_connections 1024;
}

http {
    # Basic MIME types
    types {
        text/html                             html htm shtml;
        text/css                              css;
        text/xml                              xml;
        image/gif                             gif;
        image/jpeg                            jpeg jpg;
        image/png                             png;
        application/javascript                js;
        application/json                      json;
        text/plain                            txt;
    }
    default_type  application/octet-stream;
    
    sendfile        on;
    tcp_nopush      on;
    tcp_nodelay     on;
    keepalive_timeout  65;
    
    access_log %s/access.log;
    error_log  %s/error.log;
    
    server {
        listen       8080;
        server_name  localhost;
        root         /var/www/html;
        index        index.html index.php;
        
        location / {
            try_files $uri $uri/ =404;
        }
        
        location ~ \.php$ {
            fastcgi_pass   127.0.0.1:9000;
            fastcgi_index  index.php;
            fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
            fastcgi_param  QUERY_STRING     $query_string;
            fastcgi_param  REQUEST_METHOD   $request_method;
            fastcgi_param  CONTENT_TYPE     $content_type;
            fastcgi_param  CONTENT_LENGTH   $content_length;
            fastcgi_param  REQUEST_URI      $request_uri;
            fastcgi_param  DOCUMENT_URI     $document_uri;
            fastcgi_param  DOCUMENT_ROOT    $document_root;
            fastcgi_param  SERVER_PROTOCOL  $server_protocol;
            fastcgi_param  GATEWAY_INTERFACE CGI/1.1;
            fastcgi_param  SERVER_SOFTWARE  nginx/$nginx_version;
        }
    }
}
`, "/opt/yerd/web/nginx/run", "/opt/yerd/web/nginx/logs", "/opt/yerd/web/nginx/logs")

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("failed to write nginx config: %v", err)
	}

	wi.logger.WriteLog("Created nginx configuration: %s", configPath)
	return nil
}

// createDnsmasqConfig creates a basic dnsmasq configuration
func (wi *WebInstaller) createDnsmasqConfig() error {
	configDir := "/opt/yerd/web/dnsmasq/conf"
	configPath := filepath.Join(configDir, "dnsmasq.conf")
	
	// Ensure config directory exists
	if err := utils.CreateDirectory(configDir); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	
	if utils.FileExists(configPath) {
		return nil // Don't overwrite existing config
	}

	config := fmt.Sprintf(`# YERD dnsmasq configuration
port=5353
interface=lo
bind-interfaces

# Process management
pid-file=%s/dnsmasq.pid

# Local domain
local=/dev/
domain=dev

# Cache settings
cache-size=1000

# Log settings
log-queries
log-facility=%s/dnsmasq.log

# Example entries
address=/example.dev/127.0.0.1
address=/test.dev/127.0.0.1
`, "/opt/yerd/web/dnsmasq/run", "/opt/yerd/web/dnsmasq/logs")

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("failed to write dnsmasq config: %v", err)
	}

	wi.logger.WriteLog("Created dnsmasq configuration: %s", configPath)
	return nil
}

// GetInstalledServices returns a list of installed web services
func GetInstalledServices() []string {
	var installed []string
	services := []string{"nginx", "dnsmasq"}

	for _, service := range services {
		binaryPath := GetServiceBinaryPath(service)
		if utils.FileExists(binaryPath) {
			installed = append(installed, service)
		}
	}

	return installed
}

// cleanupLogger handles log file cleanup based on installation success status.
// installSuccess: Pointer to success status.
func (wi *WebInstaller) cleanupLogger(installSuccess *bool) {
	if wi.logger == nil {
		return
	}

	if *installSuccess {
		wi.logger.DeleteLogFile()
		return
	}

	logPath := wi.logger.Close()
	fmt.Printf("\n📝 Check the detailed installation log:\n")
	fmt.Printf("   %s\n", logPath)
	fmt.Printf("\nTo view the log:\n")
	fmt.Printf("   tail -f %s\n", logPath)
}
