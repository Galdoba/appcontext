package xdg_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Galdoba/appcontext/xdg/v2"
)

var testAppName = "myapp"
var testUserHome = testHome()

func testHome() string {

	h, err := os.UserHomeDir()
	if err != nil {
		return fmt.Sprintf("panic: home not found: %v", err)
	}
	return h
}

func TestProjectPath(t *testing.T) {
	fmt.Println("empty:", xdg.Location())
	tests := []struct {
		name string
		opts []xdg.PathOption
		want string
	}{
		{
			name: "config directory",
			opts: []xdg.PathOption{xdg.ForConfig(), xdg.WithProgramName(testAppName)},
			want: filepath.Join(testUserHome, ".config", testAppName) + string(filepath.Separator),
		},
		{
			name: "data directory",
			opts: []xdg.PathOption{xdg.ForData(), xdg.WithProgramName(testAppName)},
			want: filepath.Join(testUserHome, ".local", "share", testAppName) + string(filepath.Separator),
		},
		{
			name: "cache directory",
			opts: []xdg.PathOption{xdg.ForCache(), xdg.WithProgramName(testAppName)},
			want: filepath.Join(testUserHome, ".cache", testAppName) + string(filepath.Separator),
		},
		{
			name: "state directory",
			opts: []xdg.PathOption{xdg.ForState(), xdg.WithProgramName(testAppName)},
			want: filepath.Join(testUserHome, ".local", "state", testAppName) + string(filepath.Separator),
		},
		{
			name: "with project group",
			opts: []xdg.PathOption{
				xdg.ForConfig(),
				xdg.WithProgramName(testAppName),
				xdg.WithProjectGroup("mycompany"),
			},
			want: filepath.Join(testUserHome, ".config", "mycompany", testAppName) + string(filepath.Separator),
		},
		{
			name: "with subdirectories",
			opts: []xdg.PathOption{
				xdg.ForConfig(),
				xdg.WithProgramName(testAppName),
				xdg.WithSubDir([]string{"logs", "2024"}),
			},
			want: filepath.Join(testUserHome, ".config", testAppName, "logs", "2024") + string(filepath.Separator),
		},
		{
			name: "with filename",
			opts: []xdg.PathOption{
				xdg.ForConfig(),
				xdg.WithProgramName(testAppName),
				xdg.WithFileName("config.json"),
			},
			want: filepath.Join(testUserHome, ".config", testAppName, "config.json"),
		},
		{
			name: "complete configuration",
			opts: []xdg.PathOption{
				xdg.ForData(),
				xdg.WithProgramName(testAppName),
				xdg.WithProjectGroup("mycompany"),
				xdg.WithSubDir([]string{"database", "backups"}),
				xdg.WithFileName("backup.db"),
			},
			want: filepath.Join(testUserHome, ".local", "share", "mycompany", testAppName, "database", "backups", "backup.db"),
		},
		{
			name: "empty program name",
			opts: []xdg.PathOption{xdg.ForConfig()},
			want: "",
		},
		{
			name: "invalid base directory",
			opts: []xdg.PathOption{
				xdg.WithBaseDir("invalid"),
				xdg.WithProgramName(testAppName),
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xdg.Location(tt.opts...)
			if got != tt.want {
				t.Errorf("Test '%s' failed:\n   got = %v\n  want = %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestPathOptionsCombinations(t *testing.T) {
	t.Run("multiple subdirectories", func(t *testing.T) {
		opts := []xdg.PathOption{
			xdg.ForCache(),
			xdg.WithProgramName("testapp"),
			xdg.WithSubDir([]string{"level1", "level2", "level3"}),
		}
		expected := filepath.Join(testUserHome, ".cache", "testapp", "level1", "level2", "level3") + string(filepath.Separator)
		result := xdg.Location(opts...)
		if result != expected {
			t.Errorf("Multiple subdirectories failed:\n   got = %v\n  want = %v", result, expected)
		}
	})

	t.Run("project group without subdirectories", func(t *testing.T) {
		opts := []xdg.PathOption{
			xdg.ForData(),
			xdg.WithProgramName("app"),
			xdg.WithProjectGroup("org"),
			xdg.WithFileName("data.bin"),
		}
		expected := filepath.Join(testUserHome, ".local", "share", "org", "app", "data.bin")
		result := xdg.Location(opts...)
		if result != expected {
			t.Errorf("Project group without subdirectories failed:\n   got = %v\n  want = %v", result, expected)
		}
	})
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		opts []xdg.PathOption
		want string
	}{
		{
			name: "empty subdirectory slice",
			opts: []xdg.PathOption{
				xdg.ForConfig(),
				xdg.WithProgramName(testAppName),
				xdg.WithSubDir([]string{}),
			},
			want: filepath.Join(testUserHome, ".config", testAppName) + string(filepath.Separator),
		},
		{
			name: "single element subdirectory",
			opts: []xdg.PathOption{
				xdg.ForConfig(),
				xdg.WithProgramName(testAppName),
				xdg.WithSubDir([]string{"single"}),
			},
			want: filepath.Join(testUserHome, ".config", testAppName, "single") + string(filepath.Separator),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xdg.Location(tt.opts...)
			if got != tt.want {
				t.Errorf("Edge case '%s' failed:\n   got = %v\n  want = %v", tt.name, got, tt.want)
			}
		})
	}
}
