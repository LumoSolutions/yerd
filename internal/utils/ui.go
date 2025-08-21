package utils

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/constants"
)

// PrintExtensionsGrid displays extensions in a nicely formatted grid (4 per line).
// extensions: List of extension names to display with proper spacing and alignment.
func PrintExtensionsGrid(extensions []string) {
	for i, ext := range extensions {
		if i%4 == 0 {
			fmt.Print("  ")
		}
		fmt.Printf("%-12s", ext)
		if (i+1)%4 == 0 || i == len(extensions)-1 {
			fmt.Println()
		}
	}
}

// PrintInvalidExtensionsWithSuggestions prints invalid extensions and their suggestions.
// invalid: List of invalid extension names. Returns error if invalid extensions found.
func PrintInvalidExtensionsWithSuggestions(invalid []string) {
	fmt.Println("Invalid extensions:")
	PrintExtensionsGrid(invalid)
	fmt.Println()

	for _, inv := range invalid {
		suggestions := constants.SuggestSimilarExtensions(inv)
		if len(suggestions) > 0 {
			fmt.Printf("Did you mean '%s'? Suggestions: %s\n", inv, strings.Join(suggestions, ", "))
		}
	}
}

// CheckAndPromptForSudo verifies permissions and provides helpful sudo guidance if needed.
// operation: Description of operation, command: Command name, args: Command arguments. Returns true if permissions OK.
func CheckAndPromptForSudo() bool {
	if err := CheckInstallPermissions(); err != nil {
		blue := color.New(color.FgBlue)
		red := color.New(color.FgRed)

		red.Printf("❌ Error: this command requires elevated permissions\n")
		blue.Printf("💡 This is needed to:\n")
		blue.Printf("   • Install or remove installations\n")
		blue.Printf("   • Update system-wide configuration\n\n")
		fmt.Printf("Please rerun this command with sudo\n\n")
		return false
	}
	return true
}
