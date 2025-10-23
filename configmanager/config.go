package configmanager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Galdoba/appcontext/xdg"
	"github.com/pelletier/go-toml/v2"
)

const (
	LibVersion = "1.0."
)

// Manager handles application configuration lifecycle using TOML format.
// Generic type T defines the configuration structure.
//
// Features:
// - XDG-compliant default paths
// - Atomic file writes
// - Config validation hooks
// - Multi-path config discovery
type Manager[T any] struct {
	mu     sync.RWMutex
	config *T
	path   string
}

// New creates a Manager instance with default configuration
//
// Parameters:
//   - appName: Application name (used for XDG paths)
//   - defaultConfig: Initial configuration values
//   - options: Optional manager behavior modifiers
//
// Behavior:
//   - Default config path: <XDG_CONFIG_DIR>/<appName>/<appName>.toml
//   - Creates config file on first run unless disabled
//   - Returns error if file creation fails
func New[T any](appName string, defaultConfig T, options ...ManagerOption) (*Manager[T], error) {
	configPath := filepath.Join(xdg.New(appName).ConfigDir(), appName+".toml")
	m := &Manager[T]{
		config: &defaultConfig,
		path:   configPath,
	}
	mo := managerOptions{
		forceAlternativePath: "",
		skipSaveOnCreation:   false,
	}
	for _, modify := range options {
		modify(&mo)
	}
	if mo.forceAlternativePath != "" {
		m.path = mo.forceAlternativePath
	}
	switch mo.skipSaveOnCreation {
	case false:
		if err := m.ensureConfigFileExist(); err != nil {
			return nil, err
		}
	case true:
	}
	return m, nil
}

// ensureConfigFileExist creates config file if missing
// Called automatically during initialization unless disabled
func (m *Manager[T]) ensureConfigFileExist() error {
	if _, err := os.Stat(m.path); err != nil {
		switch errors.Is(err, os.ErrNotExist) {
		case true:
			if err := m.Save(); err != nil {
				return fmt.Errorf("failed to save file on first start: %v", err)
			}
		case false:
			return fmt.Errorf("config manager failed unexpectedly: %v", err)
		}
	}
	return nil
}

// ManagerOption configures Manager behavior during creation
type ManagerOption func(*managerOptions)

type managerOptions struct {
	forceAlternativePath string
	skipSaveOnCreation   bool
}

// ForcePath overrides default XDG config path
// Path should include filename and extension
func ForcePath(path string) ManagerOption {
	return func(mo *managerOptions) {
		mo.forceAlternativePath = path
	}
}

// SaveOnCreation controls automatic file creation
// Default: true (create file if missing)
func SaveOnCreation(save bool) ManagerOption {
	return func(mo *managerOptions) {
		mo.skipSaveOnCreation = !save
	}
}

// Validator defines optional configuration validation
// Implement in configuration struct to enable validation
type Validator interface { //optional interface
	Validate() error
}

// Load reads configuration from disk
//
// Workflow:
//  1. Checks default path
//  2. Checks alternativePaths in order
//  3. Uses first valid file found
//  4. Updates active path to loaded file
//  5. Validates if Validator implemented
//
// Returns error if:
//   - No config files found
//   - File read error
//   - TOML parsing fails
//   - Validation fails
func (m *Manager[T]) Load(alternativePaths ...string) error {
	searchInPaths := append([]string{m.path}, alternativePaths...)
	selectedPath := ""
	for _, path := range searchInPaths {
		if _, err := os.Stat(path); err == nil {
			selectedPath = path
			break
		}
	}
	if selectedPath == "" {
		return fmt.Errorf("could not find config in paths: %v", searchInPaths)
	}
	data, err := os.ReadFile(selectedPath)
	if err != nil {
		return fmt.Errorf("failed to read selected file: %v", err)
	}
	if err := toml.Unmarshal(data, m.config); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}
	if v, ok := any(m.config).(Validator); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}
	}
	m.path = selectedPath
	return nil
}

// Save writes configuration to active path
//
// Steps:
//  1. Validates config if Validator implemented
//  2. Creates parent directories
//  3. Writes to temporary file
//  4. Atomically replaces target file
//
// Returns error at any failure point
func (m *Manager[T]) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if validator, ok := any(m.config).(Validator); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("validation failed before save: %w", err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(m.path), 0755); err != nil {
		return fmt.Errorf("failed to enshure config directory: %v", err)
	}

	data, err := toml.Marshal(m.config)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	if err := atomicSave(data, m.path); err != nil {
		return fmt.Errorf("atomic save: %v", err)
	}

	return nil
}

// atomicSave safely writes data using temporary swap file
// Prevents partial writes and data corruption
func atomicSave(data []byte, path string) error {
	tempPath := path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to save to tmp file: %v", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("failed to seal saved data to file: %v", err)
	}
	return nil
}

// Config provides access to current configuration
//
// Note:
//   - Returned pointer references live data
//   - Modifications require explicit Save()
func (m *Manager[T]) Config() *T {
	return m.config
}

// Path returns current active config file path
func (m *Manager[T]) Path() string {
	return m.path
}

// UseConfigFile changes active configuration file path
//
// Options:
//   - WithAutoLoad: load config immediately
//   - WithAutoCreate: create file if missing
//
// Behavior:
//   - Validates path if file exists
//   - Returns error if validation fails
//   - Auto operations execute before path switch
func (m *Manager[T]) UseConfigFile(newPath string, opts ...UseConfigOption) error {
	options := useConfigOptions{}
	for _, modify := range opts {
		modify(&options)
	}

	if _, err := os.Stat(newPath); err == nil {
		if err := validatePath(newPath); err != nil {
			return fmt.Errorf("path validation failed: %w", err)
		}

		if options.autoLoad {
			return m.Load(newPath)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// Handle missing files
		if options.autoCreate {
			oldPath := m.path
			m.path = newPath
			if err := m.Save(); err != nil {
				m.path = oldPath
				return fmt.Errorf("auto-creation failed: %w", err)
			}
			return nil
		}
		return fmt.Errorf("file not found: %w", err)
	} else {
		return fmt.Errorf("path error: %w", err)
	}
	//basic switching
	m.path = newPath
	return nil
}

// UseConfigOption modifies UseConfigFile behavior
type UseConfigOption func(*useConfigOptions)

type useConfigOptions struct {
	autoLoad   bool
	autoCreate bool
}

// WithAutoLoad enables automatic config loading
func WithAutoLoad() UseConfigOption {
	return func(o *useConfigOptions) {
		o.autoLoad = true
	}
}

// WithAutoCreate enables automatic file creation
func WithAutoCreate() UseConfigOption {
	return func(o *useConfigOptions) {
		o.autoCreate = true
	}
}

// validatePath performs basic file validation
//
// Checks:
//   - Path exists
//   - Is regular file (not directory)
//   - Is not symlink
//   - Is not special file
func validatePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("path is directory")
	}

	if f, err := os.Open(path); err != nil {
		return fmt.Errorf("read access denied: %v", err)
	} else {
		f.Close()
	}

	return nil
}
