package utils

import "time"

const (
	YerdBaseDir     = "/opt/yerd"
	YerdBinDir      = "/opt/yerd/bin"
	YerdPHPDir      = "/opt/yerd/php"
	YerdEtcDir      = "/opt/yerd/etc"
	YerdWebDir      = "/opt/yerd/web"
	SystemBinDir    = "/usr/local/bin"
	SpinnerInterval = 200 * time.Millisecond
	LogTimeFormat   = "15:04:05"
	FilePermissions = 0644
	DirPermissions  = 0755

	// FPM Configuration
	FPMUser    = "nobody"
	FPMGroup   = "nobody"
	FPMSockDir = "/opt/yerd/php/run"
	FPMPidDir  = "/opt/yerd/php/run"
	FPMLogDir  = "/opt/yerd/php/logs"

	// FPM Paths and Names
	FPMPoolDir     = "php-fpm.d"
	FPMPoolConfig  = "www.conf"
	SystemdDir     = "/etc/systemd/system"
	
	// Error Messages
	ErrEmptyPHPVersion = "PHP version cannot be empty"
)
