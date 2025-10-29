package configmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Galdoba/appcontext/xdg"
	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"
)

// Library version constant
const (
	LibVersion = "0.2.1"
)

// SerializationFormat represents supported configuration file formats
type SerializationFormat string

const (
	JSON SerializationFormat = "json"
	YAML SerializationFormat = "yaml"
	TOML SerializationFormat = "toml"
)

// Validator interface for configuration validation
type Validator interface {
	Validate() error
}

// Manager is the main configuration manager that handles loading, saving, and managing configuration files
type Manager[T any] struct {
	mu     sync.RWMutex
	config *T
	path   string
	format SerializationFormat
}

// managerOptions holds configuration options for the Manager
type managerOptions struct {
	forceAlternativePath string
	skipSaveOnCreation   bool
	format               SerializationFormat
}

// ManagerOption defines function type for configuring Manager options
type ManagerOption func(*managerOptions)

// ErrUnsupportedFormat represents an error for unsupported serialization formats
type ErrUnsupportedFormat struct {
	format SerializationFormat
}

// New creates a new configuration Manager with the specified application name and default configuration
func New[T any](appName string, defaultConfig T, options ...ManagerOption) (*Manager[T], error) {
	m := &Manager[T]{
		config: &defaultConfig,
	}
	mo := managerOptions{
		forceAlternativePath: "",
		skipSaveOnCreation:   false,
		format:               TOML,
	}
	for _, modify := range options {
		modify(&mo)
	}
	if err := validateFormat(mo.format); err != nil {
		return nil, err
	}
	m.format = mo.format

	switch mo.forceAlternativePath {
	case "":
		m.path = xdg.Location(xdg.ForConfig(), xdg.WithProgramName(appName), xdg.WithFileName(fmt.Sprintf("config.%v", m.format)))
	default:
		if fileExists(mo.forceAlternativePath) {
			if err := validatePath(mo.forceAlternativePath); err != nil {
				return nil, fmt.Errorf("forced path invalid: %v", err)
			}
		}
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

// ForcePath option forces the Manager to use a specific file path
func ForcePath(path string) ManagerOption {
	return func(mo *managerOptions) {
		mo.forceAlternativePath = path
	}
}

// SkipSaveOnCreation option controls whether to skip saving the config file on creation
func SkipSaveOnCreation(skipSave bool) ManagerOption {
	return func(mo *managerOptions) {
		mo.skipSaveOnCreation = skipSave
	}
}

// WithSerializationFormat option sets the serialization format for the configuration file
func WithSerializationFormat(format SerializationFormat) ManagerOption {
	return func(mo *managerOptions) {
		mo.format = format
	}
}

// Load reads and parses the configuration file from disk
func (m *Manager[T]) Load(alternativePaths ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	searchInPaths := append([]string{m.path}, alternativePaths...)
	selectedPath := ""
	for _, path := range searchInPaths {
		if err := validatePath(path); err != nil {
			continue
		}
		selectedPath = path
		break
	}
	if selectedPath == "" {
		return fmt.Errorf("could not find valid config paths in: %v", searchInPaths)
	}
	if ext := strings.TrimSuffix(filepath.Ext(selectedPath), "."); ext != string(m.format) {
		return fmt.Errorf("file is extention does not match serialization format (%v; %v)", ext, m.format)
	}
	data, err := os.ReadFile(selectedPath)
	if err != nil {
		return fmt.Errorf("failed to read selected file: %v", err)
	}

	old := m.Config()
	if err := m.unmarshal(data); err != nil {
		m.config = &old
		return err
	}

	if v, ok := any(m.config).(Validator); ok {
		if err := v.Validate(); err != nil {
			m.config = &old
			return fmt.Errorf("config validation failed: %w", err)
		}
	}
	m.path = selectedPath
	return nil
}

// Save writes the current configuration to disk
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

	data, err := m.marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	if err := atomicSave(data, m.path); err != nil {
		return fmt.Errorf("atomic save: %v", err)
	}

	return nil
}

// Config returns the current configuration
func (m *Manager[T]) Config() T {
	return *m.config
}

// Path returns the file path used for configuration storage
func (m *Manager[T]) Path() string {
	return m.path
}

// Error implements the error interface for ErrUnsupportedFormat
func (err *ErrUnsupportedFormat) Error() string {
	return fmt.Sprintf("unsupported serialization format: '%v'", err.format)
}

// ensureConfigFileExist creates the configuration file if it doesn't exist
func (m *Manager[T]) ensureConfigFileExist() error {
	if _, err := os.Stat(m.path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := m.Save(); err != nil {
				return fmt.Errorf("failed to save file on first start: %v", err)
			}
		}
	}
	return nil
}

// unmarshal deserializes data based on the configured format
func (m *Manager[T]) unmarshal(data []byte) error {
	switch m.format {
	case JSON:
		return json.Unmarshal(data, m.config)
	case YAML:
		return yaml.Unmarshal(data, m.config)
	case TOML:
		return toml.Unmarshal(data, m.config)
	default:
		return &ErrUnsupportedFormat{m.format}
	}
}

// marshal serializes the configuration based on the configured format
func (m *Manager[T]) marshal() ([]byte, error) {
	switch m.format {
	case JSON:
		return json.Marshal(m.config)
	case YAML:
		return yaml.Marshal(m.config)
	case TOML:
		return toml.Marshal(m.config)
	default:
		return nil, &ErrUnsupportedFormat{m.format}
	}
}

// validateFormat checks if the provided format is supported
func validateFormat(format SerializationFormat) error {
	switch format {
	case JSON, YAML, TOML:
		return nil
	default:
		return &ErrUnsupportedFormat{format}
	}
}

// validatePath checks if a file path is valid and accessible
func validatePath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("path is empty")
	}
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("path does not exist")
		}
		return err
	}

	if info.IsDir() {
		return errors.New("path is directory")
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			return errors.New("file permission does not allow read")
		}
		return err
	}
	defer file.Close()

	return nil
}

// fileExists checks if a file exists at the given path
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// atomicSave saves data to a temporary file then renames it to the target path
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

// SetPath sets new path for config file
func (m *Manager[T]) SetPath(newPath string) error {
	if !fileExists(newPath) {
		return fmt.Errorf("new path must exist")
	}
	if err := validatePath(newPath); err != nil {
		return fmt.Errorf("failed to set new path: %v", err)
	}
	return nil
}
