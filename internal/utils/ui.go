package utils

import (
	"fmt"
	"sync"
	"time"
)

type LoadingSpinner struct {
	message   string
	isRunning bool
	mutex     sync.Mutex
	stopChan  chan bool
}

// NewLoadingSpinner creates a new loading spinner with the specified message.
// message: Text to display alongside the spinner. Returns LoadingSpinner instance.
func NewLoadingSpinner(message string) *LoadingSpinner {
	return &LoadingSpinner{
		message:  message,
		stopChan: make(chan bool, 1),
	}
}

// Start begins the spinner animation in a separate goroutine.
// Does nothing if spinner is already running.
func (s *LoadingSpinner) Start() {
	s.mutex.Lock()
	if s.isRunning {
		s.mutex.Unlock()
		return
	}
	s.isRunning = true
	s.mutex.Unlock()

	go s.spin()
}

// Stop halts the spinner animation and optionally displays a completion message.
// completionMessage: Optional message to display after stopping (empty string for none).
func (s *LoadingSpinner) Stop(completionMessage string) {
	s.mutex.Lock()
	if !s.isRunning {
		s.mutex.Unlock()
		return
	}
	s.isRunning = false
	s.mutex.Unlock()

	s.stopChan <- true

	fmt.Print("\r\033[K")
	if completionMessage != "" {
		fmt.Printf("%s\n", completionMessage)
	}
}

// spin handles the actual spinner animation loop with rotating characters.
// Runs in a separate goroutine until stopped.
func (s *LoadingSpinner) spin() {
	spinChars := []string{"|", "/", "-", "\\"}
	i := 0

	for {
		select {
		case <-s.stopChan:
			return
		default:
			s.mutex.Lock()
			if !s.isRunning {
				s.mutex.Unlock()
				return
			}
			s.mutex.Unlock()

			fmt.Printf("\r[%s] %s", spinChars[i], s.message)
			i = (i + 1) % len(spinChars)
			time.Sleep(SpinnerInterval)
		}
	}
}

// PrintSuccess displays a success message with green checkmark.
// format: Printf format string, args: Format arguments.
func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf("✓ "+format+"\n", args...)
}

// PrintError displays an error message with red X mark.
// format: Printf format string, args: Format arguments.
func PrintError(format string, args ...interface{}) {
	fmt.Printf("✗ "+format+"\n", args...)
}

// PrintWarning displays a warning message with warning symbol.
// format: Printf format string, args: Format arguments.
func PrintWarning(format string, args ...interface{}) {
	fmt.Printf("⚠️  "+format+"\n", args...)
}

// PrintInfo displays an info message with cyan formatting.
// format: Printf format string, args: Format arguments.
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
