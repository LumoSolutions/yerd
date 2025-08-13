package version

import (
	"fmt"

	"github.com/fatih/color"
)

const Version = "1.0.5"

// PrintSplash displays the YERD ASCII art logo and version information with colors.
func PrintSplash() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	white := color.New(color.FgWhite)
	gray := color.New(color.FgHiBlack)

	splash := `
██╗   ██╗███████╗██████╗ ██████╗ 
╚██╗ ██╔╝██╔════╝██╔══██╗██╔══██╗
 ╚████╔╝ █████╗  ██████╔╝██║  ██║
  ╚██╔╝  ██╔══╝  ██╔══██╗██║  ██║
   ██║   ███████╗██║  ██║██████╔╝
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═════╝`

	cyan.Print(splash)
	fmt.Println()

	yellow.Printf("                     v%s\n", Version)
	fmt.Println()

	white.Println("A powerful, developer-friendly tool for managing PHP versions")
	white.Println("and local development environments with ease")
	fmt.Println()
	gray.Println("https://github.com/LumoSolutions/yerd")
	fmt.Println()
}

// GetVersion returns the current YERD version string.
func GetVersion() string {
	return Version
}
