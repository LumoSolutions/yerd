package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

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