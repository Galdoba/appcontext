package configmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Galdoba/appcontext/xdg"
	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"
)

const (
	LibVersion = "1.0.0"
)

type Manager[T any] struct {
	mu                sync.RWMutex
	config            *T
	path              string
	serializationType SerializationFormat
}

type SerializationFormat string

const (
	JSON SerializationFormat = "json"
	YAML SerializationFormat = "yaml"
	TOML SerializationFormat = "toml"
)

func New[T any](appName string, defaultConfig T, options ...ManagerOption) (*Manager[T], error) {
	m := &Manager[T]{
		config: &defaultConfig,
	}
	mo := managerOptions{
		forceAlternativePath: "",
		skipSaveOnCreation:   false,
		serializationFormat:  TOML,
	}
	for _, modify := range options {
		modify(&mo)
	}
	switch mo.forceAlternativePath {
	case "":
		m.path = filepath.Join(xdg.New(appName).ConfigDir(), "config."+string(m.serializationType))
	default:
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

type ManagerOption func(*managerOptions)

type managerOptions struct {
	forceAlternativePath string
	skipSaveOnCreation   bool
	serializationFormat  SerializationFormat
}

func ForcePath(path string) ManagerOption {
	return func(mo *managerOptions) {
		mo.forceAlternativePath = path
	}
}

func SaveOnCreation(save bool) ManagerOption {
	return func(mo *managerOptions) {
		mo.skipSaveOnCreation = !save
	}
}

func WithSerializationFormat(format SerializationFormat) ManagerOption {
	return func(mo *managerOptions) {
		mo.serializationFormat = format
	}
}

type Validator interface {
	Validate() error
}

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
	switch m.serializationType {
	case JSON:
		if err := json.Unmarshal(data, m.config); err != nil {
			return fmt.Errorf("failed to unmarshal data: %v", err)
		}
	case YAML:
		if err := yaml.Unmarshal(data, m.config); err != nil {
			return fmt.Errorf("failed to unmarshal data: %v", err)
		}
	case TOML:
		if err := toml.Unmarshal(data, m.config); err != nil {
			return fmt.Errorf("failed to unmarshal data: %v", err)
		}
	}

	if v, ok := any(m.config).(Validator); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}
	}
	m.path = selectedPath
	return nil
}

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

func (m *Manager[T]) Config() *T {
	return m.config
}

func (m *Manager[T]) Path() string {
	return m.path
}
