package constants

import (
	"os/user"
	"path/filepath"
	"strings"
)

const (
	Repo   = "LumoSolutions/yerd"
	Branch = "main"

	YerdPort              = 9000
	SiteHttpPort          = 8080
	SiteElevatedHttpPort  = 80
	SiteHttpsPort         = 8081
	SiteElevatedHttpsPort = 443

	ConfigPath         = "~/.local/yerd"
	ConfigElevatedPath = "/opt/yerd"
	ConfigName         = "server_config"

	WorkingPath         = "~/.local/yerd"
	WorkingElevatedPath = "/opt/yerd"
	LogPathRelative     = "logs"
	SrcPathRelative     = "src"
	EtcPathRelative     = "etc"
	PhpPathRelative     = "php"
	NginxPathRelative   = "nginx"
	PhpRunPathRelative  = "run"
	PhpLogsPathRelative = "logs"
	BinPathRelative     = "bin"

	BinPath         = "~/.local/bin"
	BinElevatedPath = "/usr/local/bin"

	FilePermissions = 0644
	DirPermissions  = 0755
)

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		return filepath.Join(usr.HomeDir, path[2:])
	}
	return path
}
