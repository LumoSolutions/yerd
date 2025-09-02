package manager

import (
	"fmt"
	"path/filepath"

	"github.com/lumosolutions/yerd/server/internal/core/services"
	"github.com/lumosolutions/yerd/server/internal/utils"
)

type CompileManager struct {
	Log     *services.Logger
	Command *CommandManager
}

func NewCompileManager(logger *services.Logger) *CompileManager {
	return &CompileManager{
		Log:     logger,
		Command: NewCommandManager(logger),
	}
}

func (cm *CompileManager) Configure(path string, flags []string) error {
	configureFile := filepath.Join(path, "configure")
	if !utils.FileExists(configureFile) {
		cm.Log.Info("configure", "file '%s' not found", configureFile)
		return fmt.Errorf("configure file %s does not exist", configureFile)
	}

	if err := utils.Chmod(configureFile, 0755); err != nil {
		cm.Log.Info("configure", "Unable to make configure script executable")
		return fmt.Errorf("unable to make configure executable")
	}

	// TODO:  Handle Elevation, compile as user
	args := append([]string{"/bin/bash", configureFile}, flags...)
	if _, success := cm.Command.ExecuteCommandInDir(path, args[0], args[1:]...); !success {
		cm.Log.Info("configure", "Unable to run configure script")
		return fmt.Errorf("unable to run configure script")
	}

	return nil
}

func (cm *CompileManager) Make(path string) error {
	// TODO: Handle Elevation, make as user
	nproc := cm.Command.GetProcessorCount()
	if _, success := cm.Command.ExecuteCommandInDir(path, "make", fmt.Sprintf("-j%d", nproc)); !success {
		cm.Log.Info("make", "Unable to make PHP")
		return fmt.Errorf("unable to make php")
	}

	return nil
}
