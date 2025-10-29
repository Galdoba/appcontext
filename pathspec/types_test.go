package pathspec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPath_String(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("XDG_STATE_HOME", "")
	t.Setenv("XDG_RUNTIME_DIR", "")

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
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
			want: filepath.Join(os.Getenv("HOME"), ".config", "testapp", "config", "config.yaml"),
		},
		{
			name: "data file without subcategory",
			path: Path{
				AppName:  "testapp",
				Name:     "data.json",
				BaseDir:  Data,
				PathType: FileType,
			},
			want: filepath.Join(os.Getenv("HOME"), ".local", "share", "testapp", "data.json"),
		},
		{
			name: "cache directory with subcategory",
			path: Path{
				AppName:  "testapp",
				BaseDir:  Cache,
				Name:     "thumbnails",
				PathType: DirectoryType,
			},
			want: filepath.Join(os.Getenv("HOME"), ".cache", "testapp", "thumbnails") + string(filepath.Separator),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.path.String()
			if got != tt.want {
				t.Errorf("Path.String() = %v, want %v (test name = %v)", got, tt.want, tt.name)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path    Path
		wantErr bool
	}{
		{
			name: "valid path",
			path: Path{
				AppName:     "testapp",
				Name:        "config.yaml",
				BaseDir:     Config,
				PathType:    FileType,
				Category:    CategoryConfig,
				DefaultPerm: 0644,
			},
			wantErr: false,
		},
		{
			name: "empty app name",
			path: Path{
				Name:        "config.yaml",
				BaseDir:     Config,
				PathType:    FileType,
				Category:    CategoryConfig,
				DefaultPerm: 0644,
			},
			wantErr: true,
		},
		{
			name: "empty name",
			path: Path{
				AppName:     "testapp",
				BaseDir:     Config,
				PathType:    FileType,
				Category:    CategoryConfig,
				DefaultPerm: 0644,
			},
			wantErr: true,
		},
		{
			name: "file with max children should fail",
			path: Path{
				AppName:     "testapp",
				Name:        "file.txt",
				BaseDir:     Config,
				PathType:    FileType,
				Category:    CategoryConfig,
				DefaultPerm: 0644,
				MaxChildren: 10, // Invalid for files
			},
			wantErr: true,
		},
		{
			name: "directory with max size should fail",
			path: Path{
				AppName:     "testapp",
				Name:        "data",
				BaseDir:     Data,
				PathType:    DirectoryType,
				Category:    CategoryData,
				DefaultPerm: 0755,
				MaxSize:     1024, // Invalid for directories
			},
			wantErr: true,
		},
		{
			name: "invalid category-subcategory combination",
			path: Path{
				AppName:     "testapp",
				Name:        "file.txt",
				BaseDir:     Config,
				PathType:    FileType,
				Category:    CategoryConfig,
				Subcategory: SubcategoryLogs, // Logs not allowed for config category
				DefaultPerm: 0644,
			},
			wantErr: true,
		},
		{
			name: "valid category-subcategory combination",
			path: Path{
				AppName:     "testapp",
				Name:        "config.yaml",
				BaseDir:     Config,
				PathType:    FileType,
				Category:    CategoryConfig,
				Subcategory: SubcategoryConfig, // Valid combination
				DefaultPerm: 0644,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := validate(tt.path)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("validate() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("validate() succeeded unexpectedly")
			}
		})
	}
}
