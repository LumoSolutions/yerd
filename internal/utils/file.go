package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileExists checks if a file or directory exists at the given path.
// path: File system path to check. Returns true if path exists, false otherwise.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// EnsureDirectories creates all required YERD directories if they don't exist.
// Returns error if any directory creation fails.
func EnsureDirectories() error {
	dirs := []string{
		YerdBaseDir,
		YerdBinDir,
		YerdPHPDir,
		YerdEtcDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, DirPermissions); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	return nil
}

// RemoveDirectory removes a directory and all its contents, ignoring if it doesn't exist.
// path: Directory path to remove. Returns error if removal fails.
func RemoveDirectory(path string) error {
	if !FileExists(path) {
		return nil
	}
	return os.RemoveAll(path)
}

// CreateSymlink creates a symbolic link, removing existing link if present.
// target: Path to link target, link: Path where symlink should be created. Returns error if creation fails.
func CreateSymlink(target, link string) error {
	if FileExists(link) {
		if err := os.Remove(link); err != nil {
			return fmt.Errorf("failed to remove existing symlink %s: %v", link, err)
		}
	}

	linkDir := filepath.Dir(link)
	if err := os.MkdirAll(linkDir, DirPermissions); err != nil {
		return fmt.Errorf("failed to create symlink directory %s: %v", linkDir, err)
	}

	return os.Symlink(target, link)
}

// RemoveSymlink safely removes a symbolic link with validation to prevent removing regular files.
// link: Path to symlink to remove. Returns error if not a symlink or removal fails.
func RemoveSymlink(link string) error {
	if !FileExists(link) {
		return nil
	}

	info, err := os.Lstat(link)
	if err != nil {
		return fmt.Errorf("failed to stat symlink %s: %v", link, err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("%s is not a symlink, refusing to remove", link)
	}

	return os.Remove(link)
}

// IsSymlink checks if the given path is a symbolic link.
// path: File system path to check. Returns true if path is a symlink, false otherwise.
func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

// IsBrokenSymlink checks if the path is a symlink that points to a non-existent target.
// path: File system path to check. Returns true if symlink exists but target doesn't.
func IsBrokenSymlink(path string) bool {
	if !IsSymlink(path) {
		return false
	}

	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// CreateDirectory creates a directory with proper permissions and error handling.
// path: Directory path to create. Returns error if creation fails.
func CreateDirectory(path string) error {
	if err := os.MkdirAll(path, DirPermissions); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", path, err)
	}
	return nil
}
