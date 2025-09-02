package services

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/lumosolutions/yerd/server/internal/utils"
)

type Logger struct {
	identifier string
	logFile    *os.File
	logPath    string
	mu         sync.Mutex
}

func NewLogger(identifier string) (*Logger, error) {
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s_%d.log", identifier, timestamp, time.Now().UnixNano())

	logDir := utils.GetYerdLogPath()
	fullPath := filepath.Join(logDir, fileName)

	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %v", err)
	}

	return &Logger{
		identifier: identifier,
		logFile:    file,
		logPath:    fullPath,
	}, nil
}

func (log *Logger) Info(context, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	log.writeLog("INFO", context, formattedMessage)
}

func (log *Logger) Error(context string, err error) {
	log.writeLog("ERROR", context, err.Error())
}

func (log *Logger) Warning(context, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	log.writeLog("WARNING", context, formattedMessage)
}

func (log *Logger) GetLogFile() string {
	return log.logPath
}

func (log *Logger) writeLog(logType, context, message string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	logEntry := fmt.Sprintf("%s [%s]: %s\n", logType, context, message)

	if _, err := log.logFile.WriteString(logEntry); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", err)
		return
	}

	if err := log.logFile.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to sync log file: %v\n", err)
	}
}

func (log *Logger) Close() {
	log.mu.Lock()
	defer log.mu.Unlock()

	if log.logFile != nil {
		log.logFile.Close()
	}
}
