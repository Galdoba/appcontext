package xdg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LibVersion represents the current version of the XDG library.
const (
	LibVersion = "1.0.0"
)

// PathOption defines a function type that modifies path configuration.
type PathOption func(*pathConfig)

// pathConfig holds the configuration for building application paths.
type pathConfig struct {
	programName  string   // Name of the application
	projectGroup string   // Optional group/organization name
	baseDir      string   // Base directory type (config, data, cache, state)
	subDir       []string // Additional subdirectories
	fileName     string   // Optional filename
}

// Location constructs a full path based on XDG Base Directory specification
// and the provided configuration options. Returns an empty string if required
// parameters are missing or invalid. If filename is not set path will be ended
// with separator to mark it as a directory.
func Location(opts ...PathOption) string {
	config := &pathConfig{}
	for _, opt := range opts {
		opt(config)
	}

	if config.programName == "" {
		return ""
	}

	basePath := getBaseDir(config.baseDir)
	if basePath == "" {
		return ""
	}

	path := basePath
	if config.projectGroup != "" {
		path = filepath.Join(path, config.projectGroup)
	}
	path = filepath.Join(path, config.programName)

	if len(config.subDir) != 0 {
		sections := append([]string{path}, config.subDir...)
		path = filepath.Join(sections...)
	}

	switch config.fileName {
	case "":
		fn := "tmpName"
		path = filepath.Join(path, fn)
		path = strings.TrimSuffix(path, "tmpName")
	default:
		path = filepath.Join(path, config.fileName)
	}
	if config.fileName != "" {
	}

	return path
}

// WithProgramName sets the application name for the path configuration.
func WithProgramName(name string) PathOption {
	return func(pc *pathConfig) {
		pc.programName = name
	}
}

// WithProjectGroup sets the project group/organization for the path configuration.
func WithProjectGroup(group string) PathOption {
	return func(pc *pathConfig) {
		pc.projectGroup = group
	}
}

// WithBaseDir sets the base directory type for the path configuration.
func WithBaseDir(dirType string) PathOption {
	return func(pc *pathConfig) {
		pc.baseDir = dirType
	}
}

// WithSubDir sets additional subdirectories for the path configuration.
func WithSubDir(subDir []string) PathOption {
	return func(pc *pathConfig) {
		pc.subDir = subDir
	}
}

// WithFileName sets the filename for the path configuration.
func WithFileName(fileName string) PathOption {
	return func(pc *pathConfig) {
		pc.fileName = fileName
	}
}

// ForConfig returns a PathOption that sets the base directory to config.
func ForConfig() PathOption {
	return WithBaseDir("config")
}

// ForData returns a PathOption that sets the base directory to data.
func ForData() PathOption {
	return WithBaseDir("data")
}

// ForCache returns a PathOption that sets the base directory to cache.
func ForCache() PathOption {
	return WithBaseDir("cache")
}

// ForState returns a PathOption that sets the base directory to state.
func ForState() PathOption {
	return WithBaseDir("state")
}

// ForRuntime returns a PathOption that sets the base directory to runtime.
func ForRuntime() PathOption {
	return WithBaseDir("runtime")
}

// ForTemp returns a PathOption that sets the base directory to temp.
func ForTemp() PathOption {
	return WithBaseDir("temp")
}

// runtimeHome returns the path to the runtime directory.
func runtimeHome() string {
	if path := os.Getenv("XDG_RUNTIME_DIR"); path != "" {
		return path
	}
	return filepath.Join(home(), ".local", "run")
}

// getBaseDir returns the appropriate base directory path based on directory type.
func getBaseDir(dirType string) string {
	switch dirType {
	case "config":
		return configHome()
	case "data":
		return dataHome()
	case "cache":
		return cacheHome()
	case "state":
		return stateHome()
	case "runtime":
		return runtimeHome()
	case "temp":
		return tempHome()
	default:
		return ""
	}
}

// configHome returns the path to the config home directory.
func configHome() string {
	return filepath.Join(home(), ".config")
}

// dataHome returns the path to the data home directory.
func dataHome() string {
	return filepath.Join(home(), ".local", "share")
}

// cacheHome returns the path to the cache home directory.
func cacheHome() string {
	return filepath.Join(home(), ".cache")
}

// stateHome returns the path to the state home directory.
func stateHome() string {
	return filepath.Join(home(), ".local", "state")
}

// home returns the user's home directory.
func home() string {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("home not found: %v", err))
	}
	return h
}

// tempHome returns the path to the temp directory.
func tempHome() string {
	return os.TempDir()
}
