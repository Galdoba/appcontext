package pathspec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	LibVersion = "1.0.0"
)

// String returns the absolute filesystem path for the Path
// [ai generated commentary]
func (p Path) String() string {
	basePath := p.getBaseDirPath()
	if p.Subcategory != "" {
		return filepath.Join(basePath, p.AppName, string(p.Subcategory), p.Name)
	}
	return filepath.Join(basePath, p.AppName, p.Name)
}

// getBaseDirPath returns the base directory path based on XDG specifications
// [ai generated commentary]
func (p Path) getBaseDirPath() string {
	switch p.BaseDir {
	case Config:
		if path := os.Getenv("XDG_CONFIG_HOME"); path != "" {
			return path
		}
		return filepath.Join(os.Getenv("HOME"), ".config")
	case Data:
		if path := os.Getenv("XDG_DATA_HOME"); path != "" {
			return path
		}
		return filepath.Join(os.Getenv("HOME"), ".local", "share")
	case Cache:
		if path := os.Getenv("XDG_CACHE_HOME"); path != "" {
			return path
		}
		return filepath.Join(os.Getenv("HOME"), ".cache")
	case Runtime:
		if path := os.Getenv("XDG_STATE_HOME"); path != "" {
			return path
		}
		return filepath.Join(os.Getenv("HOME"), ".local", "state")
	case Temp:
		if path := os.Getenv("XDG_RUNTIME_DIR"); path != "" {
			return path
		}
		return filepath.Join("/tmp", p.AppName+"-"+strconv.Itoa(os.Getuid()))
	default:
		return "/tmp"
	}
}

// isValid checks if a subcategory is valid for the given category
// [ai generated commentary]
func isValid(c PathCategory, s PathSubcategory) bool {
	if subcategories, exists := ValidSubcategories[c]; exists {
		return subcategories[s]
	}
	return false
}

// validate checks Path for conflicting parameters and path constructability
// [ai generated commentary]
func validate(p Path) error {
	if p.AppName == "" {
		return errors.New("AppName cannot be empty")
	}

	if p.Name == "" {
		return errors.New("Name cannot be empty")
	}

	if p.PathType == FileType && p.MaxChildren > 0 {
		return errors.New("MaxChildren cannot be set for file type")
	}

	if p.PathType == FileType && p.HasSubdirs {
		return errors.New("HasSubdirs cannot be true for file type")
	}

	if p.PathType == DirectoryType && p.MaxSize > 0 {
		return errors.New("MaxSize cannot be set for directory type")
	}

	if p.PathType == DirectoryType && p.Format != "" {
		return errors.New("Format cannot be set for directory type")
	}

	if p.RetentionDays > 0 && p.CleanupAge > 0 && p.CleanupAge > p.RetentionDays {
		return errors.New("CleanupAge cannot be greater than RetentionDays")
	}

	if p.IsVersioned && p.Category != CategoryConfig {
		return errors.New("Versioning is only applicable for config files")
	}

	if p.Subcategory != "" && !isValid(p.Category, p.Subcategory) {
		return fmt.Errorf("subcategory %s is not valid for category %d", p.Subcategory, p.Category)
	}

	if p.DefaultPerm == 0000 {
		return errors.New("permissions 0000 are invalid and make files inaccessible")
	}

	return nil
}
