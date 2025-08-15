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

// ExtractTemplateKeys returns all template keys found in a template string
func ExtractTemplateKeys(template string) []string {
	var keys []string
	start := 0
	
	for {
		startIdx := strings.Index(template[start:], "{{%")
		if startIdx == -1 {
			break
		}
		startIdx += start
		
		endIdx := strings.Index(template[startIdx:], "%}}")
		if endIdx == -1 {
			break
		}
		endIdx += startIdx
		
		keyPart := template[startIdx+3 : endIdx]
		key := strings.TrimSpace(keyPart)
		
		if key != "" {
			found := false
			for _, existing := range keys {
				if existing == key {
					found = true
					break
				}
			}
			if !found {
				keys = append(keys, key)
			}
		}
		
		start = endIdx + 3
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