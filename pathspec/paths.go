package pathspec

import (
	"errors"
	"fmt"
)

const (
	LibVersion = "1.1.0"
)

// String returns the absolute filesystem path for the Path
// [ai generated commentary]
func (p Path) String() string {
	return BuildPath(p)
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
