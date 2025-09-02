package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lumosolutions/yerd/server/internal/constants"
)

func CreateDirectory(path string) error {
	if err := os.MkdirAll(path, constants.DirPermissions); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", path, err)
	}
	return nil
}

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func WriteToFile(filename string, content []byte, perm os.FileMode) error {
	dir := filepath.Dir(filename)
	if err := CreateDirectory(dir); err != nil {
		return fmt.Errorf("failed to create directory for file %s: %w", filename, err)
	}

	err := os.WriteFile(filename, content, perm)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}
	return nil
}

func WriteStringToFile(filename string, content string, perm os.FileMode) error {
	return WriteToFile(filename, []byte(content), perm)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ReplaceDirectory(path string) error {
	os.RemoveAll(path)
	return CreateDirectory(path)
}

func ChownRecursive(path string, uid, gid int) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if err := os.Chown(name, uid, gid); err != nil {
			return fmt.Errorf("failed to chown %s: %w", name, err)
		}
		return nil
	})
}

func Chmod(filePath string, mode os.FileMode) error {
	if err := os.Chmod(filePath, mode); err != nil {
		return fmt.Errorf("failed to chmod %s: %v", filePath, err)
	}
	return nil
}

func RemoveDirectory(folderPath string) error {
	info, err := os.Stat(folderPath)
	if err != nil {
		return nil
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", folderPath)
	}

	if err := os.RemoveAll(folderPath); err != nil {
		return fmt.Errorf("failed to remove folder: %w", err)
	}

	return nil
}

func RemoveFile(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", filePath)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	return nil
}

func EnsureFoldersCreated(folders []string) error {
	for _, folder := range folders {
		if !IsDirectory(folder) {
			if err := CreateDirectory(folder); err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateSymlink(target, link string) error {
	if FileExists(link) {
		if err := os.Remove(link); err != nil {
			return fmt.Errorf("failed to remove existing symlink %s: %v", link, err)
		}
	}

	linkDir := filepath.Dir(link)
	if err := os.MkdirAll(linkDir, constants.DirPermissions); err != nil {
		return fmt.Errorf("failed to create symlink directory %s: %v", linkDir, err)
	}
	
	return os.Symlink(target, link)
}
