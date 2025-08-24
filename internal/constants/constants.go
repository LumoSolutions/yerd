package constants

import (
	"time"
)

const (
	YerdBaseDir     = "/opt/yerd"
	YerdBinDir      = "/opt/yerd/bin"
	YerdPHPDir      = "/opt/yerd/php"
	YerdEtcDir      = "/opt/yerd/etc"
	YerdWebDir      = "/opt/yerd/web"
	SystemBinDir    = "/usr/local/bin"
	GlobalPhpPath   = SystemBinDir + "/php"
	SpinnerInterval = 200 * time.Millisecond
	LogTimeFormat   = "15:04:05"
	FilePermissions = 0644
	DirPermissions  = 0755

	// Composer
	LocalComposerPath   = YerdBinDir + "/composer.phar"
	GlobalComposerPath  = "/usr/local/bin/composer"
	ComposerPharName    = "composer.phar"
	ComposerDownloadUrl = "https://getcomposer.org/download/latest-stable/composer.phar"

	// FPM Configuration
	FPMSockDir = "/opt/yerd/php/run"
	FPMPidDir  = "/opt/yerd/php/run"
	FPMLogDir  = "/opt/yerd/php/logs"

	// FPM Paths and Names
	FPMPoolDir    = "php-fpm.d"
	FPMPoolConfig = "www.conf"
	SystemdDir    = "/etc/systemd/system"

	// Error Messages
	ErrEmptyPHPVersion = "PHP version cannot be empty"

	// Web
	CertsDir = YerdWebDir + "/certs"

	// Config
	YerdConfigName = "config.json"
)
