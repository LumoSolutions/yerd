package utils

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Spinner represents an animated CLI spinner with status messages
type Spinner struct {
	mu         sync.Mutex
	writer     io.Writer
	active     bool
	phrase     string
	frames     []string
	frameIndex int
	stopChan   chan struct{}
	doneChan   chan struct{}
	hideCursor bool
	delay      time.Duration
}

// New creates a new spinner instance
func NewSpinner(phrase string) *Spinner {
	return &Spinner{
		writer:     os.Stdout,
		phrase:     phrase,
		frames:     []string{"-", "\\", "|", "/"},
		stopChan:   make(chan struct{}),
		doneChan:   make(chan struct{}),
		hideCursor: true,
		delay:      100 * time.Millisecond,
	}
}

// SetWriter sets the output writer (useful for testing)
func (s *Spinner) SetWriter(w io.Writer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writer = w
}

// SetDelay sets the animation delay between frames
func (s *Spinner) SetDelay(delay time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.delay = delay * time.Millisecond
}

// SetFrames sets custom spinner frames
func (s *Spinner) SetFrames(frames []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.frames = frames
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	if s.hideCursor {
		fmt.Fprint(s.writer, "\033[?25l") // Hide cursor
	}

	// Print initial spinner
	s.mu.Lock()
	fmt.Fprintf(s.writer, "[%s]  %s", s.frames[s.frameIndex], s.phrase)
	s.mu.Unlock()

	go s.animate()
}

// animate handles the spinner animation loop
func (s *Spinner) animate() {
	ticker := time.NewTicker(s.delay)
	defer ticker.Stop()
	defer close(s.doneChan)

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.updateSpinner()
		}
	}
}

// updateSpinner updates only the spinner line
func (s *Spinner) updateSpinner() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.frameIndex = (s.frameIndex + 1) % len(s.frames)

	// Clear current line and redraw spinner
	fmt.Fprint(s.writer, "\r\033[2K") // Move to start of line and clear it
	fmt.Fprintf(s.writer, "[%s]  %s", s.frames[s.frameIndex], s.phrase)
}

// UpdatePhrase updates the spinner phrase
func (s *Spinner) UpdatePhrase(phrase string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.phrase = phrase

	if s.active {
		// Clear current line and redraw with new phrase
		fmt.Fprint(s.writer, "\r\033[2K")
		fmt.Fprintf(s.writer, "[%s]  %s", s.frames[s.frameIndex], s.phrase)
	}
}

// AddSuccessStatus adds a status message above the spinner
func (s *Spinner) AddSuccessStatus(status string, args ...any) {
	message := fmt.Sprintf("✓ "+status, args...)
	s.addStatus(message, color.New(color.FgGreen))
}

// AddErrorStatus adds a status message above the spinner
func (s *Spinner) AddErrorStatus(status string, args ...any) {
	message := fmt.Sprintf("✗ "+status, args...)
	s.addStatus(message, color.New(color.FgRed))
}

// AddInfoStatus add a status message above the spinner
func (s *Spinner) AddInfoStatus(status string, args ...any) {
	message := fmt.Sprintf("- "+status, args...)
	s.addStatus(message, color.New(color.FgBlue))
}

// AddWarningStatus adds a status message above the spinner
func (s *Spinner) AddWarningStatus(status string, args ...any) {
	message := fmt.Sprintf("- "+status, args...)
	s.addStatus(message, color.New(color.FgYellow))
}

func (s *Spinner) addStatus(status string, outputColor *color.Color) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		// If spinner isn't active, just print the status
		outputColor.Fprintf(s.writer, "%s\n", status)
		return
	}

	// Clear the current spinner line
	fmt.Fprint(s.writer, "\r\033[2K")

	// Print the status message (this becomes a permanent line)
	outputColor.Fprintf(s.writer, "%s\n", status)

	// Redraw the spinner on the new line
	fmt.Fprintf(s.writer, "[%s]  %s", s.frames[s.frameIndex], s.phrase)
}

// Stop stops the spinner with a final message
func (s *Spinner) Stop(finalMessage string, outputColor *color.Color) {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()

	close(s.stopChan)
	<-s.doneChan

	s.mu.Lock()
	// Clear the spinner line
	fmt.Fprint(s.writer, "\r\033[2K")

	// Print the final message if provided
	if finalMessage != "" {
		outputColor.Fprintf(s.writer, "%s\n", finalMessage)
	}
	s.mu.Unlock()

	if s.hideCursor {
		fmt.Fprint(s.writer, "\033[?25h") // Show cursor
	}
}

// StopWithSuccess stops the spinner with a success message (convenience method)
func (s *Spinner) StopWithSuccess(message string, args ...any) {
	message = fmt.Sprintf(message, args...)
	s.Stop(fmt.Sprintf("🎉 %s", message), color.New(color.FgGreen))
}

// StopWithError stops the spinner with an error message (convenience method)
func (s *Spinner) StopWithError(message string, args ...any) {
	message = fmt.Sprintf(message, args...)
	s.Stop(fmt.Sprintf("✗ %s", message), color.New(color.FgRed))
}
