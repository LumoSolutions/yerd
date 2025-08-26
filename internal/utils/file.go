package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/version"
)

// FileExists checks if a file or directory exists at the given path.
// path: File system path to check. Returns true if path exists, false otherwise.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDirectory checks if the given path exists and is a directory.
// Returns false if the path doesn't exist, is a file, or if there's an error accessing it.
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// CanWriteToPath checks write permissions by traversing up directory tree to find writable parent.
// path: Target path to check. Returns true if path or parent directory is writable.
func CanWriteToPath(path string) bool {
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

// testWriteAccess attempts to create and remove a test file to verify write permissions.
// dir: Directory to test. Returns true if write access is available.
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

// CreateDirectory creates a directory with proper permissions and error handling.
// path: Directory path to create. Returns error if creation fails.
func CreateDirectory(path string) error {
	if err := os.MkdirAll(path, constants.DirPermissions); err != nil {
		LogError(err, "create-directory")
		return fmt.Errorf("failed to create directory %s: %v", path, err)
	}
	return nil
}

func ReplaceDirectory(path string) error {
	os.RemoveAll(path)
	return CreateDirectory(path)
}

// WriteToFile writes content to a file with specified permissions, overwriting if it exists
func WriteToFile(filename string, content []byte, perm os.FileMode) error {
	err := os.WriteFile(filename, content, perm)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}
	return nil
}

// WriteStringToFile writes string content to a file with specified permissions
func WriteStringToFile(filename string, content string, perm os.FileMode) error {
	return WriteToFile(filename, []byte(content), perm)
}

// Chown changes the ownership of a file or directory
func Chown(path string, uid, gid int) error {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	// Change ownership
	if err := os.Chown(path, uid, gid); err != nil {
		return fmt.Errorf("failed to change ownership of %s: %w", path, err)
	}

	return nil
}

// ChownRecursive changes ownership recursively for directories
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

// ChownSymlink changes ownership of symlink itself (not the target)
func ChownSymlink(path string, uid, gid int) error {
	return os.Lchown(path, uid, gid)
}

// Copy copies a file from src to dst using the cp command
func Copy(src, dst string) error {
	output, success := ExecuteCommand("cp", src, dst)
	if !success {
		return fmt.Errorf("failed to copy %s to %s: %s", src, dst, output)
	}
	return nil
}

// CopyRecursive copies a directory recursively from src to dst using cp -r
func CopyRecursive(src, dst string) error {
	output, success := ExecuteCommand("cp", "-rT", src, dst)
	if !success {
		return fmt.Errorf("failed to copy recursively %s to %s: %s", src, dst, output)
	}
	return nil
}

func Chmod(filePath string, mode os.FileMode) error {
	if err := os.Chmod(filePath, mode); err != nil {
		LogError(err, "chmod")
		return fmt.Errorf("failed to chmod %s: %v", filePath, err)
	}
	return nil
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

// ReadSymlink returns the target path that a symbolic link points to.
// path: File system path of the symlink. Returns the target path and nil error if successful,
// or empty string and error if path is not a symlink or cannot be read.
func ReadSymlink(path string) (string, error) {
	if !IsSymlink(path) {
		return "", fmt.Errorf("path %s is not a symbolic link", path)
	}

	target, err := os.Readlink(path)
	if err != nil {
		return "", fmt.Errorf("failed to read symlink %s: %w", path, err)
	}

	return target, nil
}

// CreateSymlink creates a symbolic link, removing existing link if present.
// target: Path to link target, link: Path where symlink should be created. Returns error if creation fails.
func CreateSymlink(target, link string) error {
	if FileExists(link) {
		if err := os.Remove(link); err != nil {
			LogError(err, "create-symlink")
			return fmt.Errorf("failed to remove existing symlink %s: %v", link, err)
		}
	}

	linkDir := filepath.Dir(link)
	if err := os.MkdirAll(linkDir, constants.DirPermissions); err != nil {
		LogError(err, "create-symlink")
		return fmt.Errorf("failed to create symlink directory %s: %v", linkDir, err)
	}

	LogInfo("create-symlink", "Creating link from %s to %s", link, target)
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

// ExtractArchive extracts a tar achieve to a folder
func ExtractArchive(archivePath, toFolder string, userCtx *UserContext) error {
	ReplaceDirectory(toFolder)
	Chown(toFolder, userCtx.UID, userCtx.GID)
	if _, success := ExecuteCommand("tar", "-xzf", archivePath, "-C", toFolder); !success {
		return fmt.Errorf("tar command failed")
	}

	return nil
}

// FetchFromGitHub downloads a file from github and returns it as
// a string value
func FetchFromGitHub(folder, file string) (string, error) {
	filePath := filepath.Join(".config", folder, file)

	url := fmt.Sprintf(
		"https://raw.githubusercontent.com/%s/%s/%s",
		version.GetRepo(),
		version.GetBranch(),
		filePath,
	)

	LogInfo("github", "Attempting to download %s", url)

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		LogError(err, "github")
		LogInfo("github", "Failed to create request")
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "YERD/1.0")
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Expires", "0")

	resp, err := client.Do(req)

	if err != nil {
		LogError(err, "github")
		LogInfo("github", "Failed to fetch data from github")
		return "", fmt.Errorf("failed to fetch from GitHub: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		LogInfo("github", "Status code was %d which was not expected", resp.StatusCode)
		return "", fmt.Errorf("GitHub returned status %d: %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		LogInfo("github", "Failed to read body")
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(content), nil
}

// RemoveFile removes a single file
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

// RemoveFolder removes a directory and all its contents
func RemoveFolder(folderPath string) error {
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

func GetWorkingDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		LogError(err, "pwd")
		return "", err
	}

	return dir, nil
}
