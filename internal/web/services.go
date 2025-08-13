package web

import (
	"path/filepath"
	"strings"

	"github.com/LumoSolutions/yerd/internal/utils"
)

// ServiceConfig represents configuration for a web service
type ServiceConfig struct {
	Name         string
	Version      string
	DownloadURL  string
	BuildFlags   []string
	Dependencies []string
	InstallPath  string
}

// GetServiceConfig returns configuration for a specific service
func GetServiceConfig(serviceName string) (*ServiceConfig, bool) {
	configs := map[string]*ServiceConfig{
		"nginx": {
			Name:        "nginx",
			Version:     "1.24.0",
			DownloadURL: "http://nginx.org/download/nginx-1.24.0.tar.gz",
			BuildFlags: []string{
				"--prefix=" + GetServiceInstallPath("nginx"),
				"--conf-path=" + GetServiceConfigPath("nginx") + "/nginx.conf",
				"--error-log-path=" + GetServiceLogPath("nginx") + "/error.log",
				"--access-log-path=" + GetServiceLogPath("nginx") + "/access.log",
				"--pid-path=" + GetServiceRunPath("nginx") + "/nginx.pid",
				"--lock-path=" + GetServiceRunPath("nginx") + "/nginx.lock",
				"--http-client-body-temp-path=" + GetServiceTempPath("nginx") + "/client_temp",
				"--http-proxy-temp-path=" + GetServiceTempPath("nginx") + "/proxy_temp",
				"--http-fastcgi-temp-path=" + GetServiceTempPath("nginx") + "/fastcgi_temp",
				"--http-uwsgi-temp-path=" + GetServiceTempPath("nginx") + "/uwsgi_temp",
				"--http-scgi-temp-path=" + GetServiceTempPath("nginx") + "/scgi_temp",
				"--with-http_ssl_module",
				"--with-http_realip_module",
				"--with-http_addition_module",
				"--with-http_sub_module",
				"--with-http_dav_module",
				"--with-http_flv_module",
				"--with-http_mp4_module",
				"--with-http_gunzip_module",
				"--with-http_gzip_static_module",
				"--with-http_auth_request_module",
				"--with-http_random_index_module",
				"--with-http_secure_link_module",
				"--with-http_degradation_module",
				"--with-http_slice_module",
				"--with-http_stub_status_module",
				"--with-http_v2_module",
				"--with-file-aio",
				"--with-threads",
			},
			Dependencies: []string{"pcre2", "zlib", "openssl"},
			InstallPath:  GetServiceInstallPath("nginx"),
		},
	}

	config, exists := configs[strings.ToLower(serviceName)]
	return config, exists
}

// GetSupportedServices returns a list of all supported web services
func GetSupportedServices() []string {
	return []string{"nginx"}
}

// IsValidService checks if a service name is supported
func IsValidService(serviceName string) bool {
	_, exists := GetServiceConfig(serviceName)
	return exists
}

// GetServiceInstallPath returns the installation path for a service
func GetServiceInstallPath(serviceName string) string {
	return filepath.Join(utils.YerdWebDir, serviceName)
}

// GetServiceConfigPath returns the configuration path for a service
func GetServiceConfigPath(serviceName string) string {
	return filepath.Join(GetServiceInstallPath(serviceName), "conf")
}

// GetServiceLogPath returns the log path for a service
func GetServiceLogPath(serviceName string) string {
	return filepath.Join(GetServiceInstallPath(serviceName), "logs")
}

// GetServiceRunPath returns the runtime path for a service
func GetServiceRunPath(serviceName string) string {
	return filepath.Join(GetServiceInstallPath(serviceName), "run")
}

// GetServiceTempPath returns the temporary files path for a service
func GetServiceTempPath(serviceName string) string {
	return filepath.Join(GetServiceInstallPath(serviceName), "temp")
}

// GetServiceBinaryPath returns the binary path for a service
func GetServiceBinaryPath(serviceName string) string {
	return filepath.Join(GetServiceInstallPath(serviceName), "sbin", serviceName)
}
