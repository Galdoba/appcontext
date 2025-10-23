package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// Save writes data to a file atomically using a temporary file and atomic rename
func Save(path string, data []byte) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("target file does not exist: %s", path)
	}

	if err := checkPermissions(path); err != nil {
		return fmt.Errorf("permission check failed: %w", err)
	}

	originalInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp.*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()

	if err := os.Chmod(tmpName, originalInfo.Mode()); err != nil {
		tmpFile.Close()
		os.Remove(tmpName)
		return fmt.Errorf("failed to set temp file permissions: %w", err)
	}

	var writeSuccess bool
	defer func() {
		if !writeSuccess {
			if tmpFile != nil {
				tmpFile.Close()
			}
			os.Remove(tmpName)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		return fmt.Errorf("failed to write data to temp file: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	tmpFile = nil

	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("atomic replace failed: %w", err)
	}

	writeSuccess = true
	return nil
}

// checkPermissions verifies read access to the file and write access to its directory
func checkPermissions(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	file.Close()

	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, "write_test.*")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()
	tmpFile.Close()
	os.Remove(tmpName)

	return nil
}
