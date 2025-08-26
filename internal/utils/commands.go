package utils

import (
	"io"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// CommandExistschecks if a command exists using the 'which' command
func CommandExists(cmd string) (string, bool) {
	output, success := ExecuteCommand("which", cmd)
	return strings.Trim(output, "\n"), success
}

// ExecuteCommand runs a command with full output logging to the specified logger.
// command: Command to run
// args: Additonal arguments
func ExecuteCommandAsUser(command string, args ...string) (string, bool) {
	LogInfo(context, "=== EXECUTING COMMAND AS USER ===")
	LogInfo(context, "Executing: %s", command)
	LogInfo(context, "With Params: %s", strings.Join(args, " "))

	userCtx, err := GetRealUser()
	if err != nil {
		LogError(err, "ExecuteCommand")
		return "", false
	}

	cmd := exec.Command(command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(userCtx.UID),
			Gid: uint32(userCtx.GID),
		},
	}

	return runCommand(cmd)
}

func ExecuteCommandInDirAsUser(directory, command string, args ...string) (string, bool) {
	LogInfo(context, "=== EXECUTING COMMAND AS USER ===")
	LogInfo(context, "Executing: %s", command)
	LogInfo(context, "With Params: %s", strings.Join(args, " "))
	LogInfo(context, "In Directory: %s", directory)

	userCtx, err := GetRealUser()
	if err != nil {
		LogError(err, "ExecuteCommand")
		return "", false
	}

	cmd := exec.Command(command, args...)
	cmd.Dir = directory
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(userCtx.UID),
			Gid: uint32(userCtx.GID),
		},
	}

	return runCommand(cmd)
}

func ExecuteCommand(command string, args ...string) (string, bool) {
	LogInfo(context, "=== EXECUTING COMMAND AS ROOT ===")
	LogInfo(context, "Executing: %s", command)
	LogInfo(context, "With Params: %s", strings.Join(args, " "))

	cmd := exec.Command(command, args...)

	return runCommand(cmd)
}

func ExecuteCommandInDir(directory, command string, args ...string) (string, bool) {
	LogInfo(context, "=== EXECUTING COMMAND AS ROOT ===")
	LogInfo(context, "Executing: %s", command)
	LogInfo(context, "With Params: %s", strings.Join(args, " "))
	LogInfo(context, "In Directory: %s", directory)

	cmd := exec.Command(command, args...)
	cmd.Dir = directory

	return runCommand(cmd)
}

func runCommand(cmd *exec.Cmd) (string, bool) {
	var output strings.Builder
	cmd.Stdout = io.MultiWriter(&output)
	cmd.Stderr = io.MultiWriter(&output)

	err := cmd.Run()
	result := output.String()

	var success bool

	if err != nil {
		success = false
		LogError(err, context)
		LogInfo("output", "%s", result)
	} else {
		success = true
		LogInfo(context, "Command executed successfully")
	}

	return result, success
}

// GetProcessorCount detects the number of CPU cores for parallel processing.
// Returns processor count or defaults to 4 if detection fails.
func GetProcessorCount() int {
	output, success := ExecuteCommand("nproc")
	if !success {
		return 4
	}

	if n, err := strconv.Atoi(strings.TrimSpace(output)); err == nil && n > 0 {
		return n
	}

	return 4
}
