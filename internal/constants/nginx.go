package constants

import "path/filepath"

// NginxConfig represents configuration for a web Nginx
type NginxConfig struct {
	Name         string
	Version      string
	DownloadURL  string
	BuildFlags   []string
	Dependencies []string
	InstallPath  string
	ConfigPath   string
	LogPath      string
	RunPath      string
	TempPath     string
	BinaryPath   string
	SourcePath   string
}

// GetNginxConfig returns configuration for a specific Nginx
func GetNginxConfig() *NginxConfig {
	return &NginxConfig{
		Name:        "nginx",
		Version:     "1.29.1",
		DownloadURL: "http://nginx.org/download/nginx-1.29.1.tar.gz",
		BuildFlags: []string{
			"--prefix=/opt/yerd/web/nginx",
			"--conf-path=/opt/yerd/web/nginx/conf/nginx.conf",
			"--error-log-path=/opt/yerd/web/nginx/logs/error.log",
			"--pid-path=/opt/yerd/web/nginx/run/nginx.pid",
			"--lock-path=/opt/yerd/web/nginx/run/nginx.lock",
			"--http-client-body-temp-path=/opt/yerd/web/nginx/temp/client_temp",
			"--http-proxy-temp-path=/opt/yerd/web/nginx/temp/proxy_temp",
			"--http-fastcgi-temp-path=/opt/yerd/web/nginx/temp/fastcgi_temp",
			"--http-uwsgi-temp-path=/opt/yerd/web/nginx/temp/uwsgi_temp",
			"--http-scgi-temp-path=/opt/yerd/web/nginx/temp/scgi_temp",
			"--with-http_ssl_module",
			"--with-http_realip_module",
			"--with-stream",
			"--with-stream_realip_module",
			"--with-cc-opt=-I/usr/include",
			"--with-ld-opt=-lsystemd",
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
		InstallPath:  getNginxInstallPath(),
		ConfigPath:   getNginxConfigPath(),
		LogPath:      getNginxLogPath(),
		RunPath:      getNginxRunPath(),
		TempPath:     getNginxTempPath(),
		BinaryPath:   getNginxBinaryPath(),
		SourcePath:   getNginxSrcPath(),
	}
}

// GetNginxInstallPath returns the installation path for a Nginx
func getNginxInstallPath() string {
	return filepath.Join(YerdWebDir, "nginx")
}

// GetNginxConfigPath returns the configuration path for a Nginx
func getNginxConfigPath() string {
	return filepath.Join(getNginxInstallPath(), "conf")
}

// GetNginxLogPath returns the log path for a Nginx
func getNginxLogPath() string {
	return filepath.Join(getNginxInstallPath(), "logs")
}

// GetNginxRunPath returns the runtime path for a Nginx
func getNginxRunPath() string {
	return filepath.Join(getNginxInstallPath(), "run")
}

// GetNginxTempPath returns the temporary files path for a Nginx
func getNginxTempPath() string {
	return filepath.Join(getNginxInstallPath(), "temp")
}

// GetNginxBinaryPath returns the binary path for a Nginx
func getNginxBinaryPath() string {
	return filepath.Join(getNginxInstallPath(), "sbin", "nginx")
}

func getNginxSrcPath() string {
	return filepath.Join(getNginxInstallPath(), "src")
}
