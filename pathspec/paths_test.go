package pathspec

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestPath_getBaseDirPath(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path Path
		want string
	}{
		{
			name: "config directory with XDG_CONFIG_HOME set",
			path: Path{
				AppName: "testapp",
				BaseDir: Config,
			},
			want: func() string {
				if path := os.Getenv("XDG_CONFIG_HOME"); path != "" {
					return path
				}
				return filepath.Join(os.Getenv("HOME"), ".config")
			}(),
		},
		{
			name: "data directory",
			path: Path{
				AppName: "testapp",
				BaseDir: Data,
			},
			want: func() string {
				if path := os.Getenv("XDG_DATA_HOME"); path != "" {
					return path
				}
				return filepath.Join(os.Getenv("HOME"), ".local", "share")
			}(),
		},
		{
			name: "temp directory with app name and user id",
			path: Path{
				AppName: "testapp",
				BaseDir: Temp,
			},
			want: func() string {
				if path := os.Getenv("XDG_RUNTIME_DIR"); path != "" {
					return path
				}
				return filepath.Join("/tmp", "testapp-"+strconv.Itoa(os.Getuid()))
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.path.getBaseDirPath()
			if got != tt.want {
				t.Errorf("Path.getBaseDirPath() = %v, want %v", got, tt.want)
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
