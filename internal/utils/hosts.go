package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	HostsFilePath   = "/etc/hosts"
	YerdStartMarker = "# YERD START - DO NOT MODIFY THIS LINE"
	YerdEndMarker   = "# YERD END - DO NOT MODIFY THIS LINE"
	DefaultIP       = "127.0.0.1"
)

type HostsManager struct {
	hostsPath string
}

func NewHostsManager() *HostsManager {
	return &HostsManager{hostsPath: HostsFilePath}
}

func NewHostsManagerWithPath(path string) *HostsManager {
	return &HostsManager{hostsPath: path}
}

// Install creates YERD comment markers in the hosts file if they don't exist
func (hm *HostsManager) Install() error {
	content, err := hm.loadHostsFile()
	if err != nil {
		return err
	}

	if hm.findYerdBounds(content) != nil {
		return nil
	}

	content = append(content, "", YerdStartMarker, YerdEndMarker)
	return hm.writeHostsFile(content)
}

// Add inserts a hostname entry within the YERD section
func (hm *HostsManager) Add(hostname string) error {
	hostname = strings.TrimSpace(hostname)
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}
	if strings.Contains(hostname, " ") {
		return fmt.Errorf("hostname cannot contain spaces")
	}

	content, err := hm.loadHostsFile()
	if err != nil {
		return err
	}

	bounds := hm.findYerdBounds(content)
	if bounds == nil {
		return fmt.Errorf("YERD section not found in hosts file")
	}

	if hm.hostExists(content, bounds, hostname) {
		return fmt.Errorf("host '%s' already exists in YERD section", hostname)
	}

	newEntry := fmt.Sprintf("%s\t%s", DefaultIP, hostname)
	newContent := make([]string, 0, len(content)+1)
	newContent = append(newContent, content[:bounds[1]]...)
	newContent = append(newContent, newEntry)
	newContent = append(newContent, content[bounds[1]:]...)

	return hm.writeHostsFile(newContent)
}

// Remove deletes a hostname entry from the YERD section
func (hm *HostsManager) Remove(hostname string) error {
	hostname = strings.TrimSpace(hostname)
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}

	content, err := hm.loadHostsFile()
	if err != nil {
		return err
	}

	bounds := hm.findYerdBounds(content)
	if bounds == nil {
		return fmt.Errorf("YERD section not found in hosts file")
	}

	found := false
	newContent := make([]string, 0, len(content))

	for i, line := range content {
		if i > bounds[0] && i < bounds[1] && hm.isHostEntry(line, hostname) {
			found = true
			continue
		}
		newContent = append(newContent, line)
	}

	if !found {
		return fmt.Errorf("host '%s' not found in YERD section", hostname)
	}

	return hm.writeHostsFile(newContent)
}

// Uninstall removes the entire YERD section and all its contents
func (hm *HostsManager) Uninstall() error {
	content, err := hm.loadHostsFile()
	if err != nil {
		return err
	}

	bounds := hm.findYerdBounds(content)
	if bounds == nil {
		return nil
	}

	newContent := append(content[:bounds[0]], content[bounds[1]+1:]...)
	
	for len(newContent) > 0 && strings.TrimSpace(newContent[len(newContent)-1]) == "" {
		newContent = newContent[:len(newContent)-1]
	}

	return hm.writeHostsFile(newContent)
}

// ListYerdHosts returns all hostnames managed by YERD
func (hm *HostsManager) ListYerdHosts() ([]string, error) {
	content, err := hm.loadHostsFile()
	if err != nil {
		return nil, err
	}

	bounds := hm.findYerdBounds(content)
	if bounds == nil {
		return []string{}, nil
	}

	var hosts []string
	for i := bounds[0] + 1; i < bounds[1]; i++ {
		line := strings.TrimSpace(content[i])
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[0] == DefaultIP {
			hosts = append(hosts, parts[1])
		}
	}

	return hosts, nil
}

// loadHostsFile reads and returns the hosts file content
func (hm *HostsManager) loadHostsFile() ([]string, error) {
	file, err := os.Open(hm.hostsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read hosts file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// writeHostsFile atomically writes content to the hosts file
func (hm *HostsManager) writeHostsFile(content []string) error {
	tempPath := hm.hostsPath + ".yerd.tmp"
	
	file, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range content {
		if _, err := file.WriteString(line + "\n"); err != nil {
			os.Remove(tempPath)
			return err
		}
	}

	return os.Rename(tempPath, hm.hostsPath)
}

// findYerdBounds returns [startIdx, endIdx] of YERD section or nil if not found
func (hm *HostsManager) findYerdBounds(content []string) []int {
	startIdx, endIdx := -1, -1
	
	for i, line := range content {
		trimmed := strings.TrimSpace(line)
		if trimmed == YerdStartMarker {
			startIdx = i
		}
		if trimmed == YerdEndMarker {
			endIdx = i
		}
	}
	
	if startIdx == -1 || endIdx == -1 {
		return nil
	}
	return []int{startIdx, endIdx}
}

// hostExists checks if hostname already exists in the YERD section
func (hm *HostsManager) hostExists(content []string, bounds []int, hostname string) bool {
	for i := bounds[0] + 1; i < bounds[1]; i++ {
		if hm.isHostEntry(content[i], hostname) {
			return true
		}
	}
	return false
}

// isHostEntry checks if a line is a host entry for the given hostname
func (hm *HostsManager) isHostEntry(line, hostname string) bool {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return false
	}

	parts := strings.Fields(line)
	return len(parts) >= 2 && parts[0] == DefaultIP && parts[1] == hostname
}