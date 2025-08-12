package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Logger struct {
	logFile *os.File
	logPath string
}

type logWriter struct {
	logger *Logger
	prefix string
}

func NewLogger(phpVersion string) (*Logger, error) {
	configDir, err := GetUserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %v", err)
	}
	
	if err := os.MkdirAll(configDir, DirPermissions); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}
	
	userCtx, err := GetRealUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get user context: %v", err)
	}
	
	if os.Geteuid() == 0 {
		if err := os.Chown(configDir, userCtx.UID, userCtx.GID); err != nil {
			return nil, fmt.Errorf("failed to set config directory ownership: %v", err)
		}
	}
	
	timestamp := time.Now().Format("20060102_150405")
	logFileName := fmt.Sprintf("install_php%s_%s.log", phpVersion, timestamp)
	logPath := filepath.Join(configDir, logFileName)
	
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, FilePermissions)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %v", err)
	}
	
	if os.Geteuid() == 0 {
		if err := os.Chown(logPath, userCtx.UID, userCtx.GID); err != nil {
			logFile.Close()
			return nil, fmt.Errorf("failed to set log file ownership: %v", err)
		}
	}
	
	logger := &Logger{
		logFile: logFile,
		logPath: logPath,
	}
	
	logger.WriteLog("=== YERD PHP %s Installation Started ===", phpVersion)
	logger.WriteLog("Timestamp: %s", time.Now().Format("2006-01-02 15:04:05"))
	logger.WriteLog("Log file: %s", logPath)
	logger.WriteLog("Running as: UID=%d (effective), User: %s, Home: %s", os.Geteuid(), userCtx.Username, userCtx.HomeDir)
	if os.Getenv("SUDO_USER") != "" {
		logger.WriteLog("Original user (SUDO_USER): %s", os.Getenv("SUDO_USER"))
	}
	logger.WriteLog("")
	
	return logger, nil
}

func (l *Logger) WriteLog(format string, args ...interface{}) {
	if l.logFile == nil {
		return
	}
	
	timestamp := time.Now().Format(LogTimeFormat)
	message := fmt.Sprintf("[%s] %s\n", timestamp, fmt.Sprintf(format, args...))
	l.logFile.WriteString(message)
	l.logFile.Sync()
}

func (l *Logger) WriteLogRaw(message string) {
	if l.logFile == nil {
		return
	}
	
	l.logFile.WriteString(message)
	l.logFile.Sync()
}

func (l *Logger) Close() string {
	logPath := l.logPath
	if l.logFile != nil {
		l.WriteLog("=== Installation log ended ===")
		l.logFile.Close()
		l.logFile = nil
	}
	return logPath
}

func (l *Logger) DeleteLogFile() {
	logPath := l.Close()
	if logPath != "" {
		os.Remove(logPath)
	}
}

func (l *Logger) GetLogPath() string {
	return l.logPath
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	if w.logger != nil && w.logger.logFile != nil {
		lines := strings.Split(string(p), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				w.logger.WriteLogRaw(fmt.Sprintf("[%s] %s: %s\n", 
					time.Now().Format("15:04:05"), w.prefix, line))
			}
		}
	}
	return len(p), nil
}

func SafeLog(logger *Logger, format string, args ...interface{}) {
	if logger != nil {
		logger.WriteLog(format, args...)
	}
}