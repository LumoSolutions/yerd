package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

type UserContext struct {
	User     *user.User
	Username string
	HomeDir  string
	UID      int
	GID      int
}

func CheckInstallPermissions() error {
	if os.Geteuid() == 0 {
		return nil
	}
	
	requiredPaths := map[string]string{
		YerdBaseDir:  "YERD installation directory",
		SystemBinDir: "system binary directory",
	}
	
	var failedPaths []string
	for path, desc := range requiredPaths {
		if !canWriteToPath(path) {
			failedPaths = append(failedPaths, fmt.Sprintf("%s (%s)", path, desc))
		}
	}
	
	if len(failedPaths) > 0 {
		return fmt.Errorf("cannot write to required directories: %s", strings.Join(failedPaths, ", "))
	}
	
	return nil
}

func testWriteAccess(dir string) bool {
	testFile := filepath.Join(dir, ".yerd_test_write")
	file, err := os.Create(testFile)
	if err != nil {
		return false
	}
	file.Close()
	os.Remove(testFile)
	return true
}

func canWriteToPath(path string) bool {
	for current := path; current != "/" && current != "."; current = filepath.Dir(current) {
		if info, err := os.Stat(current); err == nil {
			if !info.IsDir() {
				return false
			}
			return testWriteAccess(current)
		}
		if current == filepath.Dir(current) {
			break
		}
	}
	return false
}

func GetRealUser() (*UserContext, error) {
	var realUser *user.User
	var err error
	
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {
		realUser, err = user.Lookup(sudoUser)
		if err != nil {
			return nil, fmt.Errorf("failed to lookup sudo user %s: %v", sudoUser, err)
		}
	} else {
		realUser, err = user.Current()
		if err != nil {
			return nil, fmt.Errorf("failed to get current user: %v", err)
		}
	}
	
	uid, err := strconv.Atoi(realUser.Uid)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UID: %v", err)
	}
	
	gid, err := strconv.Atoi(realUser.Gid)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GID: %v", err)
	}
	
	return &UserContext{
		User:     realUser,
		Username: realUser.Username,
		HomeDir:  realUser.HomeDir,
		UID:      uid,
		GID:      gid,
	}, nil
}

func GetUserConfigDir() (string, error) {
	userCtx, err := GetRealUser()
	if err != nil {
		return "", err
	}
	
	return filepath.Join(userCtx.HomeDir, ".config", "yerd"), nil
}

func GetSudoUserIDs(sudoUser string) (int, int, error) {
	realUser, err := user.Lookup(sudoUser)
	if err != nil {
		return 0, 0, err
	}
	
	uid, err := strconv.Atoi(realUser.Uid)
	if err != nil {
		return 0, 0, err
	}
	
	gid, err := strconv.Atoi(realUser.Gid)
	if err != nil {
		return 0, 0, err
	}
	
	return uid, gid, nil
}

func CheckAndPromptForSudo(operation, command string, args ...string) bool {
	if err := CheckInstallPermissions(); err != nil {
		fmt.Printf("❌ Error: %s requires elevated permissions\n", operation)
		fmt.Printf("💡 This is needed to:\n")
		
		// Show operation-specific reasons
		switch operation {
		case "PHP installation":
			fmt.Printf("   • Create directories in /opt/yerd/\n")
			fmt.Printf("   • Install PHP binaries and libraries\n")
			fmt.Printf("   • Create symlinks in /usr/local/bin/\n")
		case "PHP removal":
			fmt.Printf("   • Remove symlinks from /usr/local/bin/\n")
			fmt.Printf("   • Clean up installation directories in /opt/yerd/\n")
		case "Setting CLI version":
			fmt.Printf("   • Create symlinks in /usr/local/bin/\n")
			fmt.Printf("   • Remove existing PHP symlinks\n")
		case "PHP update":
			fmt.Printf("   • Update PHP installations in /opt/yerd/\n")
			fmt.Printf("   • Manage symlinks in /usr/local/bin/\n")
		case "YERD update":
			fmt.Printf("   • Update YERD binary in /usr/local/bin/\n")
			fmt.Printf("   • Replace existing installation\n")
		default:
			fmt.Printf("   • Manage files in /opt/yerd/ and /usr/local/bin/\n")
		}
		
		fmt.Printf("   • Update system-wide configuration\n\n")
		fmt.Printf("Please run with sudo:\n")
		if command == "update" && operation == "YERD update" {
			fmt.Printf("   sudo yerd %s", command)
		} else {
			fmt.Printf("   sudo yerd php %s", command)
		}
		for _, arg := range args {
			fmt.Printf(" %s", arg)
		}
		fmt.Println()
		return false
	}
	return true
}