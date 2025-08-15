package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/LumoSolutions/yerd/pkg/extensions"
)

// Error message for invalid extensions
const ErrInvalidExtensions = "invalid extensions provided"

// PrintInvalidExtensionsWithSuggestions prints invalid extensions and their suggestions.
// invalid: List of invalid extension names. Returns error if invalid extensions found.
func PrintInvalidExtensionsWithSuggestions(invalid []string) error {
	if len(invalid) == 0 {
		return nil
	}

	PrintError("Invalid extensions:")
	PrintExtensionsGrid(invalid)
	fmt.Println()

	for _, inv := range invalid {
		suggestions := extensions.SuggestSimilarExtensions(inv)
		if len(suggestions) > 0 {
			fmt.Printf("Did you mean '%s'? Suggestions: %s\n", inv, strings.Join(suggestions, ", "))
		}
	}
	return fmt.Errorf(ErrInvalidExtensions)
}

// PrintExtensionList prints a sorted list of extensions with a title.
// title: Title to display, extensionList: List of extensions to print.
func PrintExtensionList(title string, extensionList []string) {
	if len(extensionList) == 0 {
		return
	}

	PrintWarning("%s:", title)
	sort.Strings(extensionList)
	PrintExtensionsGrid(extensionList)
	fmt.Println()
}

// CreateExtensionMap creates a map from extension slice for efficient lookups.
// extensions: Slice of extension names. Returns map with extension names as keys.
func CreateExtensionMap(extensions []string) map[string]bool {
	extMap := make(map[string]bool)
	for _, ext := range extensions {
		extMap[ext] = true
	}
	return extMap
}

// SliceFromExtensionMap converts extension map back to sorted slice.
// extMap: Map of extension names. Returns sorted slice of extension names.
func SliceFromExtensionMap(extMap map[string]bool) []string {
	var result []string
	for ext := range extMap {
		result = append(result, ext)
	}
	sort.Strings(result)
	return result
}