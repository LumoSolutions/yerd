package utils

import (
	"fmt"
	"strings"
)

// TemplateData represents key-value pairs for template replacement
type TemplateData map[string]string

// Template processes a template string by replacing {{% key %}} patterns with values
func Template(template string, data TemplateData) string {
	result := template

	for key, value := range data {
		placeholder := fmt.Sprintf("{{%% %s %%}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// TemplateWithDefaults processes a template with fallback values for missing keys
func TemplateWithDefaults(template string, data TemplateData, defaults TemplateData) string {
	merged := make(TemplateData)

	for k, v := range defaults {
		merged[k] = v
	}

	for k, v := range data {
		merged[k] = v
	}

	return Template(template, merged)
}

// ExtractTemplateKeys returns all unique template keys found in a template string
func ExtractTemplateKeys(template string) []string {
	uniqueKeys := extractUniqueKeys(template)
	return toSlice(uniqueKeys)
}

// extractUniqueKeys finds all unique keys in the template
func extractUniqueKeys(template string) map[string]bool {
	uniqueKeys := make(map[string]bool)
	remaining := template

	for {
		key, nextPos := findNextKey(remaining)
		if key == "" {
			break
		}
		uniqueKeys[key] = true
		remaining = remaining[nextPos:]
	}

	return uniqueKeys
}

// findNextKey finds the next template key and returns it along with the position after it
func findNextKey(template string) (string, int) {
	const startDelim = "{{%"
	const endDelim = "%}}"

	startIdx := strings.Index(template, startDelim)
	if startIdx == -1 {
		return "", 0
	}

	afterStart := startIdx + len(startDelim)
	endIdx := strings.Index(template[afterStart:], endDelim)
	if endIdx == -1 {
		return "", 0
	}

	key := strings.TrimSpace(template[afterStart : afterStart+endIdx])
	nextPos := afterStart + endIdx + len(endDelim)

	return key, nextPos
}

// toSlice converts a map of keys to a slice
func toSlice(uniqueKeys map[string]bool) []string {
	keys := make([]string, 0, len(uniqueKeys))
	for key := range uniqueKeys {
		if key != "" {
			keys = append(keys, key)
		}
	}
	return keys
}

// ValidateTemplate checks if all template keys have corresponding data values
func ValidateTemplate(template string, data TemplateData) []string {
	keys := ExtractTemplateKeys(template)
	var missing []string

	for _, key := range keys {
		if _, exists := data[key]; !exists {
			missing = append(missing, key)
		}
	}

	return missing
}
