package version

import (
	"fmt"

	"github.com/fatih/color"
)

const Version = "1.1.0"
const Branch = "feat/nginx"
const Repo = "LumoSolutions/yerd"

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

	white.Println("A powerful, developer-friendly tool for managing")
	white.Println("multiple PHP versions and local development ")
	white.Println("environments via nginx with ease")
	fmt.Println()
	gray.Println("Consider contributing today")
	gray.Println("https://github.com/LumoSolutions/yerd")
	fmt.Println()
}

// GetVersion returns the current YERD version string.
func GetVersion() string {
	return Version
}

// GetBranch returns the current branch that files should
// be downloaded from, useful during development, should
// be reset to main on build for release
func GetBranch() string {
	return Branch
}

// GetRepo returns the current repo used to store files
// and will be where the .config files are located
func GetRepo() string {
	return Repo
}
