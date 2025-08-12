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

func NewLoadingSpinner(message string) *LoadingSpinner {
	return &LoadingSpinner{
		message:  message,
		stopChan: make(chan bool, 1),
	}
}

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