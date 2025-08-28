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

// HostsManager manages YERD-specific entries in the hosts file
type HostsManager struct {
	hostsPath string
}

// NewHostsManager creates a new HostsManager with the default hosts file path
func NewHostsManager() *HostsManager {
	return &HostsManager{hostsPath: HostsFilePath}
}

// NewHostsManagerWithPath creates a new HostsManager with a custom hosts file path
func NewHostsManagerWithPath(path string) *HostsManager {
	return &HostsManager{hostsPath: path}
}

// Install creates YERD comment markers in the hosts file if they don't exist
func (hm *HostsManager) Install() error {
	content, err := hm.loadHostsFile()
	if err != nil {
		LogError(err, "hosts")
		return fmt.Errorf("failed to load hosts file: %w", err)
	}

	if hm.findYerdBounds(content) != nil {
		LogInfo("hosts", "YERD section already exists")
		return nil
	}

	if len(content) > 0 && strings.TrimSpace(content[len(content)-1]) != "" {
		content = append(content, "")
	}
	content = append(content, YerdStartMarker, YerdEndMarker)

	if err := hm.writeHostsFile(content); err != nil {
		LogError(err, "hosts")
		return fmt.Errorf("failed to install YERD section: %w", err)
	}

	LogInfo("hosts", "YERD section installed successfully")
	return nil
}

// Add inserts a hostname entry within the YERD section
func (hm *HostsManager) Add(hostname string) error {
	if err := hm.validateHostname(hostname); err != nil {
		return err
	}
	hostname = strings.TrimSpace(hostname)

	if err := hm.ensureInstalled(); err != nil {
		return fmt.Errorf("failed to ensure YERD section: %w", err)
	}

	content, err := hm.loadHostsFile()
	if err != nil {
		return fmt.Errorf("failed to load hosts file: %w", err)
	}

	bounds := hm.findYerdBounds(content)
	if bounds == nil {
		return fmt.Errorf("YERD section not found after installation")
	}

	if hm.hostExists(content, bounds, hostname) {
		return fmt.Errorf("host '%s' already exists in YERD section", hostname)
	}

	newEntry := fmt.Sprintf("%s\t%s", DefaultIP, hostname)
	newContent := hm.insertLine(content, bounds.endLine, newEntry)

	if err := hm.writeHostsFile(newContent); err != nil {
		return fmt.Errorf("failed to add host: %w", err)
	}

	LogInfo("hosts", "Added host: %s", hostname)
	return nil
}

// Remove deletes a hostname entry from the YERD section
func (hm *HostsManager) Remove(hostname string) error {
	if err := hm.validateHostname(hostname); err != nil {
		return err
	}
	hostname = strings.TrimSpace(hostname)

	content, err := hm.loadHostsFile()
	if err != nil {
		return fmt.Errorf("failed to load hosts file: %w", err)
	}

	bounds := hm.findYerdBounds(content)
	if bounds == nil {
		return fmt.Errorf("YERD section not found in hosts file")
	}

	removed := false
	newContent := make([]string, 0, len(content))

	for i, line := range content {
		if hm.shouldRemoveLine(i, line, bounds, hostname) {
			removed = true
			LogInfo("hosts", "Removing host: %s", hostname)
			continue
		}
		newContent = append(newContent, line)
	}

	if !removed {
		return fmt.Errorf("host '%s' not found in YERD section", hostname)
	}

	if err := hm.writeHostsFile(newContent); err != nil {
		return fmt.Errorf("failed to remove host: %w", err)
	}

	return nil
}

// Uninstall removes the entire YERD section and all its contents
func (hm *HostsManager) Uninstall() error {
	content, err := hm.loadHostsFile()
	if err != nil {
		return fmt.Errorf("failed to load hosts file: %w", err)
	}

	bounds := hm.findYerdBounds(content)
	if bounds == nil {
		LogInfo("hosts", "YERD section not found, nothing to uninstall")
		return nil
	}

	newContent := hm.removeSection(content, bounds)

	if err := hm.writeHostsFile(newContent); err != nil {
		return fmt.Errorf("failed to uninstall YERD section: %w", err)
	}

	LogInfo("hosts", "YERD section uninstalled successfully")
	return nil
}

// ListYerdHosts returns all hostnames managed by YERD
func (hm *HostsManager) ListYerdHosts() ([]string, error) {
	content, err := hm.loadHostsFile()
	if err != nil {
		return nil, fmt.Errorf("failed to load hosts file: %w", err)
	}

	bounds := hm.findYerdBounds(content)
	if bounds == nil {
		return []string{}, nil
	}

	hosts := make([]string, 0)
	for i := bounds.startLine + 1; i < bounds.endLine; i++ {
		hostname := hm.extractHostname(content[i])
		if hostname != "" {
			hosts = append(hosts, hostname)
		}
	}

	return hosts, nil
}

// yerdBounds represents the boundaries of the YERD section
type yerdBounds struct {
	startLine int
	endLine   int
}

// ensureInstalled makes sure the YERD section exists in the hosts file
func (hm *HostsManager) ensureInstalled() error {
	content, err := hm.loadHostsFile()
	if err != nil {
		return err
	}

	if hm.findYerdBounds(content) == nil {
		return hm.Install()
	}
	return nil
}

// validateHostname checks if a hostname is valid
func (hm *HostsManager) validateHostname(hostname string) error {
	hostname = strings.TrimSpace(hostname)
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}
	if strings.Contains(hostname, " ") || strings.Contains(hostname, "\t") {
		return fmt.Errorf("hostname cannot contain whitespace")
	}
	if strings.Contains(hostname, "#") {
		return fmt.Errorf("hostname cannot contain '#' character")
	}
	return nil
}

// loadHostsFile reads and returns the hosts file content as lines
func (hm *HostsManager) loadHostsFile() ([]string, error) {
	file, err := os.Open(hm.hostsPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// writeHostsFile atomically writes content to the hosts file
func (hm *HostsManager) writeHostsFile(content []string) error {
	tempPath := hm.hostsPath + ".yerd.tmp"

	file, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		if _, err := os.Stat(tempPath); err == nil {
			os.Remove(tempPath)
		}
	}()

	writer := bufio.NewWriter(file)
	for i, line := range content {
		if i > 0 {
			writer.WriteString("\n")
		}
		writer.WriteString(line)
	}

	if len(content) > 0 {
		writer.WriteString("\n")
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return os.Rename(tempPath, hm.hostsPath)
}

// findYerdBounds returns the boundaries of the YERD section or nil if not found
func (hm *HostsManager) findYerdBounds(content []string) *yerdBounds {
	startIdx, endIdx := -1, -1

	for i, line := range content {
		trimmed := strings.TrimSpace(line)
		if trimmed == YerdStartMarker {
			startIdx = i
		} else if trimmed == YerdEndMarker {
			endIdx = i
			break
		}
	}

	LogInfo("hosts", "bounds check - start: %d, end: %d", startIdx, endIdx)

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return nil
	}

	return &yerdBounds{
		startLine: startIdx,
		endLine:   endIdx,
	}
}

// hostExists checks if hostname already exists in the YERD section
func (hm *HostsManager) hostExists(content []string, bounds *yerdBounds, hostname string) bool {
	for i := bounds.startLine + 1; i < bounds.endLine; i++ {
		if hm.isHostEntry(content[i], hostname) {
			return true
		}
	}
	return false
}

// isHostEntry checks if a line is a host entry for the given hostname
func (hm *HostsManager) isHostEntry(line, hostname string) bool {
	extracted := hm.extractHostname(line)
	return extracted != "" && extracted == hostname
}

// extractHostname extracts the hostname from a hosts file entry line
func (hm *HostsManager) extractHostname(line string) string {
	line = strings.TrimSpace(line)

	if line == "" || strings.HasPrefix(line, "#") {
		return ""
	}

	parts := strings.Fields(line)
	if len(parts) >= 2 && parts[0] == DefaultIP {
		return parts[1]
	}

	return ""
}

// insertLine inserts a new line at the specified position
func (hm *HostsManager) insertLine(content []string, position int, line string) []string {
	newContent := make([]string, 0, len(content)+1)
	newContent = append(newContent, content[:position]...)
	newContent = append(newContent, line)
	newContent = append(newContent, content[position:]...)
	return newContent
}

// removeSection removes the YERD section from content
func (hm *HostsManager) removeSection(content []string, bounds *yerdBounds) []string {
	newContent := make([]string, 0, len(content))
	newContent = append(newContent, content[:bounds.startLine]...)

	if bounds.endLine+1 < len(content) {
		newContent = append(newContent, content[bounds.endLine+1:]...)
	}

	for len(newContent) > 0 && strings.TrimSpace(newContent[len(newContent)-1]) == "" {
		newContent = newContent[:len(newContent)-1]
	}

	return newContent
}

// shouldRemoveLine determines if a line should be removed during Remove operation
func (hm *HostsManager) shouldRemoveLine(index int, line string, bounds *yerdBounds, hostname string) bool {
	if index <= bounds.startLine || index >= bounds.endLine {
		return false
	}
	return hm.isHostEntry(line, hostname)
}
