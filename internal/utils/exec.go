package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// ExecuteCommand runs a system command with arguments and returns combined output.
// command: Command to execute, args: Command arguments. Returns output string and error.
func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// VerifyPHPInstallation checks if a PHP binary exists and is functional.
// phpPath: Path to PHP binary to verify. Returns error if binary missing or non-functional.
func VerifyPHPInstallation(phpPath string) error {
	if _, err := os.Stat(phpPath); os.IsNotExist(err) {
		return fmt.Errorf("PHP binary not found at %s", phpPath)
	}

	cmd := exec.Command(phpPath, "-v")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("PHP binary is not executable: %v", err)
	}

	if !strings.Contains(string(output), "PHP") {
		return fmt.Errorf("invalid PHP binary output")
	}

	return nil
}

// ExecuteCommandAsUser runs a command as the original user with dropped privileges and logging.
// logger: Logger instance, command: Command to run, args: Arguments. Returns output and error.
func ExecuteCommandAsUser(logger *Logger, command string, args ...string) (string, error) {
	userCtx, err := GetRealUser()
	if err != nil {
		if logger != nil {
			logger.WriteLog("Failed to get real user context: %v", err)
		}
		return "", fmt.Errorf("failed to get user context: %v", err)
	}

	if logger != nil {
		logger.WriteLog("Executing as user %s (UID: %d): %s %s", userCtx.Username, userCtx.UID, command, strings.Join(args, " "))
	}

	cmd := exec.Command(command, args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(userCtx.UID),
			Gid: uint32(userCtx.GID),
		},
	}

	cmd.Env = append(os.Environ(), "HOME="+userCtx.HomeDir)

	var output strings.Builder
	if logger != nil {
		cmd.Stdout = io.MultiWriter(&output, &logWriter{logger, "STDOUT"})
		cmd.Stderr = io.MultiWriter(&output, &logWriter{logger, "STDERR"})
	} else {
		cmd.Stdout = &output
		cmd.Stderr = &output
	}

	err = cmd.Run()
	result := output.String()

	if logger != nil {
		if err != nil {
			logger.WriteLog("Command failed with error: %v", err)
		} else {
			logger.WriteLog("Command completed successfully")
		}
		logger.WriteLog("---")
	}

	return result, err
}

// ExecuteCommandWithLogging runs a command with full output logging to the specified logger.
// logger: Logger instance, command: Command to run, args: Arguments. Returns output and error.
func ExecuteCommandWithLogging(logger *Logger, command string, args ...string) (string, error) {
	if logger != nil {
		logger.WriteLog("Executing: %s %s", command, strings.Join(args, " "))
	}

	cmd := exec.Command(command, args...)

	var output strings.Builder
	cmd.Stdout = io.MultiWriter(&output, &logWriter{logger, "STDOUT"})
	cmd.Stderr = io.MultiWriter(&output, &logWriter{logger, "STDERR"})

	err := cmd.Run()
	result := output.String()

	if logger != nil {
		if err != nil {
			logger.WriteLog("Command failed with error: %v", err)
		} else {
			logger.WriteLog("Command completed successfully")
		}
		logger.WriteLog("---")
	}

	return result, err
}

// GetProcessorCount detects the number of CPU cores for parallel processing.
// Returns processor count or defaults to 4 if detection fails.
func GetProcessorCount() int {
	output, err := ExecuteCommand("nproc")
	if err != nil {
		return 4
	}

	if n, err := strconv.Atoi(strings.TrimSpace(output)); err == nil && n > 0 {
		return n
	}

	return 4
}
