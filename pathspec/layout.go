package pathspec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NewLayout creates a new Layout with the given app name and paths
// [ai generated commentary]
func NewLayout(appname string, paths []Path) (*Layout, error) {
	layout := &Layout{
		AppName:      appname,
		ConfigPaths:  []Path{},
		DataPaths:    []Path{},
		CachePaths:   []Path{},
		RuntimePaths: []Path{},
		TempPaths:    []Path{},
	}

	for i, path := range paths {
		if path.AppName == "" {
			path.AppName = appname
			paths[i] = path
		}

		if err := validate(path); err != nil {
			return nil, fmt.Errorf("validation failed for path %s: %w", path.Name, err)
		}

		switch path.BaseDir {
		case Config:
			layout.ConfigPaths = append(layout.ConfigPaths, path)
		case Data:
			layout.DataPaths = append(layout.DataPaths, path)
		case Cache:
			layout.CachePaths = append(layout.CachePaths, path)
		case Runtime:
			layout.RuntimePaths = append(layout.RuntimePaths, path)
		case Temp:
			layout.TempPaths = append(layout.TempPaths, path)
		}
	}

	return layout, nil
}

// Import loads a Layout from a JSON file and validates all paths
// [ai generated commentary]
func Import(filePath string) (*Layout, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var layout Layout
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&layout); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Validate all paths in the layout
	allPaths := layout.GetAllPaths()
	for _, path := range allPaths {
		// Ensure AppName is set from layout if missing in individual path
		if path.AppName == "" {
			path.AppName = layout.AppName
		}

		if err := validate(path); err != nil {
			return nil, fmt.Errorf("validation failed for path %s: %w", path.Name, err)
		}
	}

	return &layout, nil
}

// GetAllPaths returns all paths combined from all categories
// [ai generated commentary]
func (l *Layout) GetAllPaths() []Path {
	var allPaths []Path
	allPaths = append(allPaths, l.ConfigPaths...)
	allPaths = append(allPaths, l.DataPaths...)
	allPaths = append(allPaths, l.CachePaths...)
	allPaths = append(allPaths, l.RuntimePaths...)
	allPaths = append(allPaths, l.TempPaths...)
	return allPaths
}

// Generate creates all directories and empty files in the layout
// [ai generated commentary]
func (l *Layout) Generate() error {
	var errors []string
	allPaths := l.GetAllPaths()

	for _, path := range allPaths {
		fullPath := path.String()

		if err := validate(path); err != nil {
			return fmt.Errorf("generation imposible: invalid path: %v: %v", path, err)
		}

		switch path.PathType {
		case DirectoryType:
			if err := os.MkdirAll(fullPath, os.FileMode(path.DefaultPerm)); err != nil {
				errors = append(errors, fmt.Sprintf("directory %s: %v", fullPath, err))
			}
		case FileType:
			parentDir := filepath.Dir(fullPath)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				errors = append(errors, fmt.Sprintf("parent directory for %s: %v", fullPath, err))
				continue
			}

			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, os.FileMode(path.DefaultPerm))
				if err != nil {
					errors = append(errors, fmt.Sprintf("file %s: %v", fullPath, err))
				} else {
					file.Close()
				}
			}
		case SymlinkType:
			errors = append(errors, fmt.Sprintf("symlink %s: symlink creation not supported", fullPath))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("generation errors:\n%s", strings.Join(errors, "\n"))
	}
	return nil
}

// Assess evaluates real files against layout parameters and returns assessment results
// [ai generated commentary]
func (l *Layout) Assess() ([]string, error) {
	var assessmentErrors []string
	allPaths := l.GetAllPaths()

	for _, path := range allPaths {
		fullPath := path.String()
		info, err := os.Stat(fullPath)

		if os.IsNotExist(err) {
			if path.IsMandatory {
				assessmentErrors = append(assessmentErrors, fmt.Sprintf("mandatory path does not exist: %s", fullPath))
			}
			continue
		}

		if err != nil {
			assessmentErrors = append(assessmentErrors, fmt.Sprintf("cannot access path %s: %v", fullPath, err))
			continue
		}

		// Check path type match
		if (path.PathType == FileType && info.IsDir()) || (path.PathType == DirectoryType && !info.IsDir()) {
			assessmentErrors = append(assessmentErrors, fmt.Sprintf("path type mismatch: %s is %s but expected %s",
				fullPath, getActualType(info), getExpectedType(path.PathType)))
			continue
		}

		// Check permissions
		actualPerm := info.Mode().Perm()
		expectedPerm := os.FileMode(path.DefaultPerm)
		if actualPerm == 0000 {
			assessmentErrors = append(assessmentErrors, fmt.Sprintf("invalid permissions 0000 for %s: file is completely inaccessible", fullPath))
		}
		if actualPerm != expectedPerm {
			assessmentErrors = append(assessmentErrors, fmt.Sprintf("permissions mismatch for %s: has %04o, expected %04o",
				fullPath, actualPerm, expectedPerm))
		}

		// Check file size for files
		if path.PathType == FileType && path.MaxSize > 0 {
			if info.Size() > int64(path.MaxSize) {
				assessmentErrors = append(assessmentErrors, fmt.Sprintf("file size exceeds limit for %s: %d > %d",
					fullPath, info.Size(), path.MaxSize))
			}
		}

		// Check directory children count
		if path.PathType == DirectoryType && path.MaxChildren > 0 {
			entries, err := os.ReadDir(fullPath)
			if err != nil {
				assessmentErrors = append(assessmentErrors, fmt.Sprintf("cannot read directory %s: %v", fullPath, err))
			} else if uint32(len(entries)) > path.MaxChildren {
				assessmentErrors = append(assessmentErrors, fmt.Sprintf("directory children count exceeded for %s: %d > %d",
					fullPath, len(entries), path.MaxChildren))
			}
		}
	}

	var overallError error
	if len(assessmentErrors) > 0 {
		overallError = fmt.Errorf("assessment errors: %d", len(assessmentErrors))
	}

	return assessmentErrors, overallError
}

// getActualType returns string representation of actual filesystem object type
// [ai generated commentary]
func getActualType(info os.FileInfo) string {
	if info.IsDir() {
		return "directory"
	}
	return "file"
}

// getExpectedType returns string representation of expected Path type
// [ai generated commentary]
func getExpectedType(pathType PathType) string {
	switch pathType {
	case FileType:
		return "file"
	case DirectoryType:
		return "directory"
	case SymlinkType:
		return "symlink"
	default:
		return "unknown"
	}
}
