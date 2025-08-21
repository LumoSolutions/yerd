package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lumosolutions/yerd/internal/constants"
)

type UserContext struct {
	User     *user.User
	Username string
	HomeDir  string
	UID      int
	GID      int
}

// GetUserConfigDir returns the YERD configuration directory path for the real user.
// Returns config directory path or error if user context cannot be determined.
func GetUserConfigDir() (string, error) {
	userCtx, err := GetRealUser()
	if err != nil {
		return "", err
	}

	return filepath.Join(userCtx.HomeDir, ".config", "yerd"), nil
}

// GetRealUser determines the actual user context, handling sudo scenarios correctly.
// Returns UserContext with user details including UID/GID or error if user lookup fails.
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

// CheckInstallPermissions verifies write access to directories required for YERD operations.
// Returns error if insufficient permissions, nil if all required paths are writable.
func CheckInstallPermissions() error {
	if os.Geteuid() == 0 {
		return nil
	}

	requiredPaths := map[string]string{
		constants.YerdBaseDir:  "YERD installation directory",
		constants.SystemBinDir: "system binary directory",
	}

	var failedPaths []string
	for path, desc := range requiredPaths {
		if !CanWriteToPath(path) {
			failedPaths = append(failedPaths, fmt.Sprintf("%s (%s)", path, desc))
		}
	}

	if len(failedPaths) > 0 {
		return fmt.Errorf("cannot write to required directories: %s", strings.Join(failedPaths, ", "))
	}

	return nil
}
