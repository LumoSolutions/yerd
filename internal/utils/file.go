package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

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

func RemoveDirectory(path string) error {
	if !FileExists(path) {
		return nil
	}
	return os.RemoveAll(path)
}

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

func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

func IsBrokenSymlink(path string) bool {
	if !IsSymlink(path) {
		return false
	}
	
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}