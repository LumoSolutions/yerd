package manager

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/lumosolutions/yerd/server/internal/config"
	"github.com/lumosolutions/yerd/server/internal/constants"
	"github.com/lumosolutions/yerd/server/internal/core/services"
	"github.com/lumosolutions/yerd/server/internal/utils"
)

type SystemdManager struct {
	Service     string
	Version     string
	phpInfo     *config.PhpInfo
	serviceName string
	Log         *services.Logger
	Cm          *CommandManager
}

func (sm *SystemdManager) Reload() error {
	args := []string{"daemon-reload"}
	if !utils.IsRunningElevated() {
		args = []string{"--user", "daemon-reload"}
	}

	if _, success := sm.Cm.ExecuteCommand("systemctl", args...); !success {
		sm.Log.Info("systemd", "Failed to reload daemons")
		return fmt.Errorf("unable to reload systemd daemons")
	}

	return nil
}

func (sm *SystemdManager) Start() error {
	return sm.runCommand("start")
}

func (sm *SystemdManager) Stop() error {
	return sm.runCommand("enable")
}

func (sm *SystemdManager) Enable() error {
	return sm.runCommand("enable")
}

func (sm *SystemdManager) Disable() error {
	return sm.runCommand("disable")
}

func (sm *SystemdManager) runCommand(action string) error {
	args := sm.buildArgs(action)
	if _, success := sm.Cm.ExecuteCommand("systemctl", args...); !success {
		sm.Log.Info("systemd", "Failed to %s service %s", action, sm.serviceName)
		return fmt.Errorf("unable to %s systemd service %s", action, sm.serviceName)
	}

	return nil
}

func (sm *SystemdManager) buildArgs(action string) []string {
	var args []string

	if !utils.IsRunningElevated() {
		args = append(args, "--user")
	}

	args = append(args, action, sm.serviceName)

	return args
}

func NewSystemdManager(service, version string, log *services.Logger) (*SystemdManager, error) {
	manager := &SystemdManager{
		Service: service,
		Version: version,
		Log:     log,
		Cm:      NewCommandManager(log),
	}

	if service == "php" {
		appConfig, _ := config.GetConfig()
		phpInfo, exists := appConfig.Php[version]
		if !exists {
			return nil, fmt.Errorf("PHP %s is not installed", version)
		}

		manager.phpInfo = phpInfo
		manager.serviceName = fmt.Sprintf("yerd-php%s-fpm", phpInfo.Version)
	}

	return manager, nil
}

func CreatePhpSystemdService(info *config.PhpInfo, log *services.Logger) (*SystemdManager, error) {
	fileName := fmt.Sprintf("yerd-php%s-fpm.service", info.Version)
	filePath := filepath.Join(getServiceDir(), fileName)

	data := utils.TemplateData{
		"version":          info.Version,
		"pid_path":         info.FpmPidLocation,
		"fpm_binary_path":  filepath.Join(info.InstallPath, "sbin", "php-fpm"),
		"main_config_path": info.FpmConfig,
	}

	dm := NewDownloadManager(log, time.Second*30)
	content, err := dm.FetchFromGitHub("php", "systemd.conf")
	if err != nil {
		return nil, err
	}

	content = utils.Template(content, data)
	if err := utils.WriteStringToFile(filePath, content, constants.FilePermissions); err != nil {
		log.Error("systemd-php", err)
		return nil, err
	}

	sm := &SystemdManager{
		Service:     "php",
		Version:     info.Version,
		phpInfo:     info,
		serviceName: fmt.Sprintf("yerd-php%s-fpm", info.Version),
		Log:         log,
		Cm:          NewCommandManager(log),
	}

	if err := sm.Reload(); err != nil {
		return nil, err
	}

	return sm, nil
}

func getServiceDir() string {
	if utils.IsRunningElevated() {
		return "/etc/systemd/system"
	}

	user, _ := utils.GetUser()
	return filepath.Join(user.HomeDir, ".config", "systemd", "user")
}
