package manager

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/lumosolutions/yerd/server/internal/core/services"
	"github.com/lumosolutions/yerd/server/internal/utils"
)

var context = "command"

type CommandManager struct {
	Log *services.Logger
}

func NewCommandManager(logger *services.Logger) *CommandManager {
	return &CommandManager{
		Log: logger,
	}
}

func (cm *CommandManager) ExecuteCommand(command string, args ...string) (string, bool) {
	cm.Log.Info(context, "=== EXECUTING COMMAND ===")
	cm.Log.Info(context, "Executing: %s", command)
	cm.Log.Info(context, "With Params: %s", strings.Join(args, " "))

	cmd := exec.Command(command, args...)

	return cm.runCommand(cmd)
}

func (cm *CommandManager) ExecuteCommandInDir(directory, command string, args ...string) (string, bool) {
	cm.Log.Info(context, "=== EXECUTING COMMAND ===")
	cm.Log.Info(context, "Executing: %s", command)
	cm.Log.Info(context, "With Params: %s", strings.Join(args, " "))
	cm.Log.Info(context, "In Directory: %s", directory)

	cmd := exec.Command(command, args...)
	cmd.Dir = directory

	return cm.runCommand(cmd)
}

func (cm *CommandManager) CommandExists(cmd string) (string, bool) {
	output, success := cm.ExecuteCommand("which", cmd)
	return strings.Trim(output, "\n"), success
}

func (cm *CommandManager) runCommand(cmd *exec.Cmd) (string, bool) {
	var output strings.Builder
	cmd.Stdout = io.MultiWriter(&output)
	cmd.Stderr = io.MultiWriter(&output)

	err := cmd.Run()
	result := output.String()

	var success bool

	if err != nil {
		success = false
		cm.Log.Error(context, err)
		cm.Log.Info("output", "%s", result)
	} else {
		success = true
		cm.Log.Info(context, "Command executed successfully")
	}

	return result, success
}

func (cm *CommandManager) GetProcessorCount() int {
	output, success := cm.ExecuteCommand("nproc")
	if !success {
		return 4
	}

	if n, err := strconv.Atoi(strings.TrimSpace(output)); err == nil && n > 0 {
		return n
	}

	return 4
}

func (cm *CommandManager) CopyRecursive(src, dst string) error {
	output, success := cm.ExecuteCommand("cp", "-rT", src, dst)
	if !success {
		return fmt.Errorf("failed to copy recursively %s to %s: %s", src, dst, output)
	}
	return nil
}

func (cm *CommandManager) ExtractArchive(archivePath, toFolder string) error {
	utils.ReplaceDirectory(toFolder)
	if _, success := cm.ExecuteCommand("tar", "-xzf", archivePath, "-C", toFolder); !success {
		return fmt.Errorf("tar command failed")
	}

	return nil
}
