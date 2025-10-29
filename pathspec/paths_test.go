package pathspec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildPath(t *testing.T) {
	tests := []struct {
		name string
		path Path
		want string
	}{
		{
			name: "config file with subcategory",
			path: Path{
				AppName:     "testapp",
				Name:        "config.yaml",
				BaseDir:     Config,
				Subcategory: SubcategoryConfig,
				PathType:    FileType,
			},
			want: func() string {
				// Ожидаемый путь через xdg
				home, _ := os.UserHomeDir()
				return filepath.Join(home, ".config", "testapp", "config", "config.yaml")
			}(),
		},
		{
			name: "data directory without subcategory",
			path: Path{
				AppName:  "testapp",
				Name:     "data",
				BaseDir:  Data,
				PathType: DirectoryType,
			},
			want: func() string {
				home, _ := os.UserHomeDir()
				return filepath.Join(home, ".local", "share", "testapp", "data") + string(filepath.Separator)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildPath(tt.path)
			if got != tt.want {
				t.Errorf("BuildPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestIsValid(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		category    PathCategory
		subcategory PathSubcategory
		want        bool
	}{
		{
			name:        "valid config subcategory",
			category:    CategoryConfig,
			subcategory: SubcategoryConfig,
			want:        true,
		},
		{
			name:        "invalid config subcategory",
			category:    CategoryConfig,
			subcategory: SubcategoryLogs,
			want:        false,
		},
		{
			name:        "valid data subcategory",
			category:    CategoryData,
			subcategory: SubcategoryDatabase,
			want:        true,
		},
		{
			name:        "valid runtime subcategory",
			category:    CategoryRuntime,
			subcategory: SubcategoryStats,
			want:        true,
		},
		{
			name:        "non-existent category",
			category:    PathCategory(99), // Non-existent category
			subcategory: SubcategoryConfig,
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValid(tt.category, tt.subcategory)
			if got != tt.want {
				t.Errorf("isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
