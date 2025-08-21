package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lumosolutions/yerd/internal/constants"
)

type Logger struct {
	logFile *os.File
	logPath string
	mu      sync.Mutex // Mutex for thread-safe operations
}

type logWriter struct {
	logger *Logger
	prefix string
}

var (
	instance *Logger
	once     sync.Once
)

// GetLogger returns the singleton Logger instance, creating it if necessary.
// Uses sync.Once to ensure thread-safe initialization.
func GetLogger() (*Logger, error) {
	var err error
	once.Do(func() {
		instance, err = initLogger()
	})

	if err != nil {
		return nil, err
	}

	return instance, nil
}

// initLogger creates the singleton logger instance with proper file setup.
// This is called only once through sync.Once.
func initLogger() (*Logger, error) {
	configDir, err := GetUserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %v", err)
	}

	if err := os.MkdirAll(configDir, constants.DirPermissions); err != nil {
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

	// Create log filename with format: yerd_YYYYMMDD_timestamp.log
	now := time.Now()
	date := now.Format("20060102")
	timestamp := now.Format("150405")
	nanos := now.Nanosecond()
	logFileName := fmt.Sprintf("yerd_%s_%s%d.log", date, timestamp, nanos/1000)
	logPath := filepath.Join(configDir, logFileName)

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, constants.FilePermissions)
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

	logger.WriteLog("=== YERD Logger Initialized ===")
	logger.WriteLog("Timestamp: %s", time.Now().Format("2006-01-02 15:04:05"))
	logger.WriteLog("Log file: %s", logPath)
	logger.WriteLog("Running as: UID=%d (effective), User: %s, Home: %s", os.Geteuid(), userCtx.Username, userCtx.HomeDir)
	if os.Getenv("SUDO_USER") != "" {
		logger.WriteLog("Original user (SUDO_USER): %s", os.Getenv("SUDO_USER"))
	}
	logger.WriteLog("")

	return logger, nil
}

// ResetLogger forces creation of a new logger instance (useful for testing or log rotation).
// This should be used sparingly as it breaks the singleton pattern temporarily.
func ResetLogger() {
	if instance != nil {
		instance.Close()
	}
	instance = nil
	once = sync.Once{}
}

// WriteLog writes a timestamped message to the log file with formatting.
// Thread-safe through mutex locking.
func (l *Logger) WriteLog(format string, args ...interface{}) {
	if l == nil || l.logFile == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format(constants.LogTimeFormat)
	message := fmt.Sprintf("[%s] %s\n", timestamp, fmt.Sprintf(format, args...))
	l.logFile.WriteString(message)
	l.logFile.Sync()
}

// WriteLogRaw writes a message to the log file without timestamp formatting.
// Thread-safe through mutex locking.
func (l *Logger) WriteLogRaw(message string) {
	if l == nil || l.logFile == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.logFile.WriteString(message)
	l.logFile.Sync()
}

// Close finalizes and closes the log file, returning the log file path.
// Thread-safe through mutex locking.
func (l *Logger) Close() string {
	if l == nil {
		return ""
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	logPath := l.logPath
	if l.logFile != nil {
		l.WriteLogRaw(fmt.Sprintf("[%s] === Logger Closed ===\n",
			time.Now().Format(constants.LogTimeFormat)))
		l.logFile.Close()
		l.logFile = nil
	}
	return logPath
}

// DeleteLogFile closes and removes the log file from disk.
// Used for cleanup when logs are no longer needed.
func (l *Logger) DeleteLogFile() {
	logPath := l.Close()
	if logPath != "" {
		os.Remove(logPath)
	}
}

// GetLogPath returns the file system path to the current log file.
func (l *Logger) GetLogPath() string {
	if l == nil {
		return ""
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.logPath
}

// CreateLogWriter creates an io.Writer that writes to the logger with a prefix.
// Useful for capturing command output.
func (l *Logger) CreateLogWriter(prefix string) *logWriter {
	return &logWriter{
		logger: l,
		prefix: prefix,
	}
}

// Write implements io.Writer interface for capturing command output to log.
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

// SafeLog writes to the singleton logger instance if available.
// Attempts to get the logger if not provided, handles errors gracefully.
func SafeLog(format string, args ...interface{}) {
	logger, err := GetLogger()
	if err != nil || logger == nil {
		// Fallback to stderr if logger unavailable
		fmt.Fprintf(os.Stderr, "[LOG ERROR] "+format+"\n", args...)
		return
	}
	logger.WriteLog(format, args...)
}

// LogError is a convenience function for logging errors.
func LogError(err error, context string) {
	if err == nil {
		return
	}
	SafeLog("ERROR [%s]: %v", context, err)
}

// LogInfo is a convenience function for logging informational messages.
func LogInfo(context string, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	SafeLog("INFO [%s]: %s", context, message)
}

// LogDebug is a convenience function for logging debug messages.
func LogDebug(context string, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	SafeLog("DEBUG [%s]: %s", context, message)
}

// LogWarning is a convenience function for logging warning messages.
func LogWarning(context string, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	SafeLog("WARNING [%s]: %s", context, message)
}
