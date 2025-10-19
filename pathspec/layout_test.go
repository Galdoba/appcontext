package pathspec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLayout(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		appname string
		paths   []Path
		want    *Layout
		wantErr bool
	}{
		{
			name:    "valid layout creation",
			appname: "testapp",
			paths: []Path{
				{
					AppName:     "testapp",
					Name:        "config.yaml",
					BaseDir:     Config,
					PathType:    FileType,
					Category:    CategoryConfig,
					Priority:    PriorityCritical,
					DefaultPerm: 0644,
					OwnerOnly:   true,
					IsMandatory: true,
				},
			},
			want: &Layout{
				AppName: "testapp",
				ConfigPaths: []Path{
					{
						AppName:     "testapp",
						Name:        "config.yaml",
						BaseDir:     Config,
						PathType:    FileType,
						Category:    CategoryConfig,
						Priority:    PriorityCritical,
						DefaultPerm: 0644,
						OwnerOnly:   true,
						IsMandatory: true,
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty app name in path uses layout app name",
			appname: "testapp",
			paths: []Path{
				{
					Name:        "config.yaml",
					BaseDir:     Config,
					PathType:    FileType,
					Category:    CategoryConfig,
					Priority:    PriorityCritical,
					DefaultPerm: 0644,
				},
			},
			want: &Layout{
				AppName: "testapp",
				ConfigPaths: []Path{
					{
						AppName:     "testapp",
						Name:        "config.yaml",
						BaseDir:     Config,
						PathType:    FileType,
						Category:    CategoryConfig,
						Priority:    PriorityCritical,
						DefaultPerm: 0644,
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid path validation fails",
			appname: "testapp",
			paths: []Path{
				{
					Name:        "", // Empty name should fail validation
					BaseDir:     Config,
					PathType:    FileType,
					Category:    CategoryConfig,
					Priority:    PriorityCritical,
					DefaultPerm: 0644,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "paths distributed to correct categories",
			appname: "testapp",
			paths: []Path{
				{
					Name:        "config.yaml",
					BaseDir:     Config,
					PathType:    FileType,
					Category:    CategoryConfig,
					DefaultPerm: 0644,
				},
				{
					Name:        "data.json",
					BaseDir:     Data,
					PathType:    FileType,
					Category:    CategoryData,
					DefaultPerm: 0644,
				},
				{
					Name:        "cache.dat",
					BaseDir:     Cache,
					PathType:    FileType,
					Category:    CategoryCache,
					DefaultPerm: 0644,
				},
			},
			want: &Layout{
				AppName: "testapp",
				ConfigPaths: []Path{
					{
						AppName:  "testapp",
						Name:     "config.yaml",
						BaseDir:  Config,
						PathType: FileType,
						Category: CategoryConfig,
					},
				},
				DataPaths: []Path{
					{
						AppName:  "testapp",
						Name:     "data.json",
						BaseDir:  Data,
						PathType: FileType,
						Category: CategoryData,
					},
				},
				CachePaths: []Path{
					{
						AppName:  "testapp",
						Name:     "cache.dat",
						BaseDir:  Cache,
						PathType: FileType,
						Category: CategoryCache,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewLayout(tt.appname, tt.paths)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewLayout() failed: %v (test name = %v)", gotErr, tt.name)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewLayout() succeeded unexpectedly")
			}

			if got.AppName != tt.want.AppName {
				t.Errorf("NewLayout() AppName = %v, want %v", got.AppName, tt.want.AppName)
			}

			if len(got.ConfigPaths) != len(tt.want.ConfigPaths) {
				t.Errorf("NewLayout() ConfigPaths count = %v, want %v", len(got.ConfigPaths), len(tt.want.ConfigPaths))
			}

			if len(got.DataPaths) != len(tt.want.DataPaths) {
				t.Errorf("NewLayout() DataPaths count = %v, want %v", len(got.DataPaths), len(tt.want.DataPaths))
			}
		})
	}
}

func TestLayout_Generate(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		layout   *Layout
		wantErr  bool
		setupEnv func()
	}{
		{
			name: "generate directories and files successfully",
			layout: &Layout{
				AppName: "testapp",
				ConfigPaths: []Path{
					{
						AppName:     "testapp",
						Name:        "config.yaml",
						BaseDir:     Config,
						Subcategory: SubcategoryConfig,
						PathType:    FileType,
						DefaultPerm: 0644,
					},
				},
				DataPaths: []Path{
					{
						AppName:     "testapp",
						Name:        "data",
						BaseDir:     Data,
						Category:    CategoryData,
						Subcategory: SubcategoryStorage,
						PathType:    DirectoryType,
						DefaultPerm: 0755,
					},
				},
			},
			wantErr: false,
			setupEnv: func() {
				os.Setenv("XDG_CONFIG_HOME", tmpDir)
				os.Setenv("XDG_DATA_HOME", tmpDir)
			},
		},
		{
			name: "fail on invalid permissions",
			layout: &Layout{
				AppName: "testapp",
				ConfigPaths: []Path{
					{
						AppName:     "testapp",
						Name:        "config.yaml",
						BaseDir:     Config,
						Subcategory: SubcategoryConfig,
						PathType:    FileType,
						DefaultPerm: 0000, // No permissions
					},
				},
			},
			wantErr: true,
			setupEnv: func() {
				os.Setenv("XDG_CONFIG_HOME", tmpDir)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			gotErr := tt.layout.Generate()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Layout.Generate() failed: %v (test name = %v)", gotErr, tt.name)
				}
				return
			}
			if tt.wantErr {
				t.Fatalf("Layout.Generate() succeeded unexpectedly: %v", tt.name)
			}

			// Verify files were created
			allPaths := tt.layout.GetAllPaths()
			for _, path := range allPaths {
				fullPath := path.String()
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					t.Errorf("Layout.Generate() failed to create path: %s", fullPath)
				}
			}
		})
	}
}

func TestLayout_Assess(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		layout   *Layout
		setupFS  func()
		wantErr  bool
		wantMsgs int // number of expected assessment messages
	}{
		{
			name: "all paths compliant",
			layout: &Layout{
				AppName: "testapp",
				ConfigPaths: []Path{
					{
						AppName:     "testapp",
						Name:        "config.yaml",
						BaseDir:     Config,
						Subcategory: SubcategoryConfig,
						PathType:    FileType,
						DefaultPerm: 0644,
						IsMandatory: true,
					},
				},
			},
			setupFS: func() {
				os.Setenv("XDG_CONFIG_HOME", tmpDir)
				configPath := filepath.Join(tmpDir, "testapp", "config", "config.yaml")
				os.MkdirAll(filepath.Dir(configPath), 0755)
				os.WriteFile(configPath, []byte("test"), 0644)
			},
			wantErr:  false,
			wantMsgs: 0,
		},
		{
			name: "missing mandatory file",
			layout: &Layout{
				AppName: "testapp",
				ConfigPaths: []Path{
					{
						AppName:     "testapp",
						Name:        "missing.yaml",
						BaseDir:     Config,
						Subcategory: SubcategoryConfig,
						PathType:    FileType,
						DefaultPerm: 0644,
						IsMandatory: true,
					},
				},
			},
			setupFS: func() {
				os.Setenv("XDG_CONFIG_HOME", tmpDir)
			},
			wantErr:  true,
			wantMsgs: 1, // One error for missing mandatory file
		},
		{
			name: "permission mismatch",
			layout: &Layout{
				AppName: "testapp",
				ConfigPaths: []Path{
					{
						AppName:     "testapp",
						Name:        "config.yaml",
						BaseDir:     Config,
						Subcategory: SubcategoryConfig,
						PathType:    FileType,
						DefaultPerm: 0600, // Expecting 0600
						IsMandatory: true,
					},
				},
			},
			setupFS: func() {
				os.Setenv("XDG_CONFIG_HOME", tmpDir)
				configPath := filepath.Join(tmpDir, "testapp", "config", "config.yaml")
				os.MkdirAll(filepath.Dir(configPath), 0755)
				os.WriteFile(configPath, []byte("test"), 0644) // Actual 0644
			},
			wantErr:  true,
			wantMsgs: 1, // One error for permission mismatch
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFS()
			gotMsgs, gotErr := tt.layout.Assess()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Layout.Assess() failed: %v", gotErr)
				}
			} else if tt.wantErr {
				t.Fatal("Layout.Assess() succeeded unexpectedly")
			}

			if len(gotMsgs) != tt.wantMsgs {
				t.Errorf("Layout.Assess() messages count = %v, want %v", len(gotMsgs), tt.wantMsgs)
			}
		})
	}
}

func TestImport(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		jsonData string
		want     *Layout
		wantErr  bool
	}{
		{
			name: "valid layout import",
			jsonData: `{
				"app_name": "testapp",
				"app_version": "1.0.0",
				"config_paths": [
					{
						"name": "config.yaml",
						"base_dir": 0,
						"subcategory": "config",
						"path_type": 0,
						"category": 0,
						"priority": 0,
						"default_perm": 420,
						"owner_only": true,
						"is_mandatory": true,
						"is_backed_up": true,
						"is_versioned": true,
						"format": "yaml"
					}
				]
			}`,
			want: &Layout{
				AppName:    "testapp",
				AppVersion: "1.0.0",
				ConfigPaths: []Path{
					{
						AppName:     "testapp",
						Name:        "config.yaml",
						BaseDir:     Config,
						Subcategory: SubcategoryConfig,
						PathType:    FileType,
						Category:    CategoryConfig,
						Priority:    PriorityCritical,
						DefaultPerm: 0644,
						OwnerOnly:   true,
						IsMandatory: true,
						IsBackedUp:  true,
						IsVersioned: true,
						Format:      "yaml",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid JSON data",
			jsonData: `{
				"app_name": "testapp",
				"config_paths": [
					{
						"name": "",  // Empty name should fail validation
						"base_dir": 0,
						"path_type": 0,
						"category": 0
					}
				]
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:     "non-existent file",
			jsonData: "",
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string

			if tt.name == "non-existent file" {
				filePath = filepath.Join(tmpDir, "nonexistent.json")
			} else {
				filePath = filepath.Join(tmpDir, "layout.json")
				err := os.WriteFile(filePath, []byte(tt.jsonData), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			got, gotErr := Import(filePath)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Import() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Import() succeeded unexpectedly")
			}

			if got.AppName != tt.want.AppName {
				t.Errorf("Import() AppName = %v, want %v", got.AppName, tt.want.AppName)
			}
			if len(got.ConfigPaths) != len(tt.want.ConfigPaths) {
				t.Errorf("Import() ConfigPaths count = %v, want %v", len(got.ConfigPaths), len(tt.want.ConfigPaths))
			}
		})
	}
}

func TestLayout_GetAllPaths(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		layout *Layout
		want   int // expected total number of paths
	}{
		{
			name: "multiple categories",
			layout: &Layout{
				AppName: "testapp",
				ConfigPaths: []Path{
					{Name: "config1.yaml", BaseDir: Config, PathType: FileType},
					{Name: "config2.yaml", BaseDir: Config, PathType: FileType},
				},
				DataPaths: []Path{
					{Name: "data1.json", BaseDir: Data, PathType: FileType},
				},
				CachePaths: []Path{
					{Name: "cache1", BaseDir: Cache, PathType: DirectoryType},
				},
			},
			want: 4,
		},
		{
			name: "empty layout",
			layout: &Layout{
				AppName: "testapp",
			},
			want: 0,
		},
		{
			name: "single category",
			layout: &Layout{
				AppName: "testapp",
				ConfigPaths: []Path{
					{Name: "config.yaml", BaseDir: Config, PathType: FileType},
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.layout.GetAllPaths()
			if len(got) != tt.want {
				t.Errorf("Layout.GetAllPaths() count = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestGetActualType(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		setup    func() os.FileInfo
		expected string
	}{
		{
			name: "directory type",
			setup: func() os.FileInfo {
				dirPath := filepath.Join(tmpDir, "testdir")
				os.MkdirAll(dirPath, 0755)
				info, _ := os.Stat(dirPath)
				return info
			},
			expected: "directory",
		},
		{
			name: "file type",
			setup: func() os.FileInfo {
				filePath := filepath.Join(tmpDir, "testfile.txt")
				os.WriteFile(filePath, []byte("test"), 0644)
				info, _ := os.Stat(filePath)
				return info
			},
			expected: "file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := tt.setup()
			got := getActualType(info)
			if got != tt.expected {
				t.Errorf("getActualType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetExpectedType(t *testing.T) {
	tests := []struct {
		name     string
		pathType PathType
		expected string
	}{
		{
			name:     "file type",
			pathType: FileType,
			expected: "file",
		},
		{
			name:     "directory type",
			pathType: DirectoryType,
			expected: "directory",
		},
		{
			name:     "symlink type",
			pathType: SymlinkType,
			expected: "symlink",
		},
		{
			name:     "unknown type",
			pathType: PathType(99),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getExpectedType(tt.pathType)
			if got != tt.expected {
				t.Errorf("getExpectedType() = %v, want %v", got, tt.expected)
			}
		})
	}
}
