package jsonstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

const (
	LibVersion                  = "2.0.0"
	Compact    MarshalingMethod = 0
	Indent     MarshalingMethod = 1
	Hybrid     MarshalingMethod = 2
)

type MarshalingMethod int

var ErrRecordNotFound = errors.New("record not found")
var ErrRecordExist = errors.New("record already exist")

type JsonDB[T any] struct {
	mu               sync.RWMutex
	data             map[string]T
	path             string
	autoSave         bool
	marshalingMethod MarshalingMethod
	prefix           string
	indent           string
}

type options struct {
	autoSave         bool
	marshalingMethod MarshalingMethod
	prefix           string
	indent           string
}

type DB_Option func(*options)

// WithAutoSave enables or disables automatic saving to file after each modification
func WithAutoSave(autoSave bool) DB_Option {
	return func(o *options) {
		o.autoSave = autoSave
	}
}

// WithCompactMarshaling configures the database to use compact JSON formatting
func WithCompactMarshaling() DB_Option {
	return func(o *options) {
		o.marshalingMethod = Compact
	}
}

// WithIndentMarshaling configures the database to use indented JSON formatting
func WithIndentMarshaling(prefix, indent string) DB_Option {
	return func(o *options) {
		o.marshalingMethod = Indent
		o.prefix = prefix
		o.indent = indent
	}
}

// New creates a new JSON database instance
// If the file exists, data is loaded from it. Otherwise, an empty database is created
func New[T any](path string, opts ...DB_Option) (*JsonDB[T], error) {
	db := &JsonDB[T]{
		path:             path,
		data:             make(map[string]T),
		mu:               sync.RWMutex{},
		marshalingMethod: Hybrid,
		indent:           "  ",
	}
	optionSet := options{}
	for _, modify := range opts {
		modify(&optionSet)
	}
	db.autoSave = optionSet.autoSave
	db.marshalingMethod = optionSet.marshalingMethod
	db.indent = optionSet.indent
	db.prefix = optionSet.prefix

	file, err := os.ReadFile(path)
	switch err {
	case nil:
		if err := json.Unmarshal(file, &db.data); err != nil {
			return nil, err
		}
	default:
		if os.IsNotExist(err) {
			return db, nil
		}
		return nil, err
	}

	return db, nil
}

// Path returns the file system path where the database is stored
func (db *JsonDB[T]) Path() string {
	return db.path
}

// Load creates a database instance from an existing file
// Returns an error if the file doesn't exist or cannot be read
func Load[T any](path string, opts ...DB_Option) (*JsonDB[T], error) {
	db := &JsonDB[T]{
		path:             path,
		data:             make(map[string]T),
		marshalingMethod: Hybrid,
		indent:           "  ",
	}

	optionSet := options{}
	for _, modify := range opts {
		modify(&optionSet)
	}
	db.autoSave = optionSet.autoSave
	db.marshalingMethod = optionSet.marshalingMethod
	db.indent = optionSet.indent
	db.prefix = optionSet.prefix

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &db.data); err != nil {
		return nil, err
	}

	return db, nil
}

// Save writes the current database state to file atomically
// Uses a temporary file and atomic rename to ensure data consistency
func (db *JsonDB[T]) Save() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.internalSave()
}

func (db *JsonDB[T]) internalSave() error {
	data, err := db.Marshal()
	if err != nil {
		return fmt.Errorf("failed db marshaling: %v", err)
	}

	if err := writeToFile(data, db.path); err != nil {
		return fmt.Errorf("failed to write storage: %v", err)
	}

	return nil
}

// Marshal returns the JSON representation of the database
// Format depends on the configured marshaling method
func (db *JsonDB[T]) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	switch db.marshalingMethod {
	case Compact:
		return json.Marshal(db.data)
	case Indent:
		return json.MarshalIndent(db.data, db.prefix, db.indent)
	case Hybrid:
		if _, err := buf.WriteString("{\n"); err != nil {
			return nil, err
		}

		keys := make([]string, 0, len(db.data))
		for k := range db.data {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i, key := range keys {
			keyJSON, err := json.Marshal(key)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal entry '%v': %v", key, err)
			}
			if _, err := buf.WriteString("  "); err != nil {
				return nil, err
			}
			if _, err := buf.Write(keyJSON); err != nil {
				return nil, err
			}
			if _, err := buf.WriteString(": "); err != nil {
				return nil, err
			}

			valJSON, err := json.Marshal(db.data[key])
			if err != nil {
				return nil, err
			}
			if _, err := buf.Write(valJSON); err != nil {
				return nil, err
			}

			if i < len(keys)-1 {
				if _, err := buf.WriteString(","); err != nil {
					return nil, err
				}
			}
			if _, err := buf.WriteString("\n"); err != nil {
				return nil, err
			}
		}
		if _, err := buf.WriteString("}"); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("unexpected Marshaling conclusion")
}

// Insert adds a new record to the database
// Returns error if ID is empty or record already exists
func (db *JsonDB[T]) Insert(id string, value T) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if id == "" {
		return fmt.Errorf("empty entry id")
	}
	if _, exists := db.data[id]; exists {
		return ErrRecordExist
	}

	oldData := maps.Clone(db.data)

	db.data[id] = value
	if db.autoSave {
		if err := db.internalSave(); err != nil {
			db.data = oldData
			return fmt.Errorf("failed to save db: %v", err)
		}
	}
	return nil
}

// Contains checks if a record with the given ID exists in the database
func (db *JsonDB[T]) Contains(id string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	_, exists := db.data[id]
	return exists
}

// Count returns the number of records in the database
func (db *JsonDB[T]) Count() int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return len(db.data)
}

// Get retrieves a record by ID
// Returns error if record is not found
func (db *JsonDB[T]) Get(id string) (T, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	val, exists := db.data[id]
	if !exists {
		var noneRecord T
		return noneRecord, ErrRecordNotFound
	}
	return val, nil
}

// Update modifies an existing record in the database
// Returns error if ID is empty or record doesn't exist
func (db *JsonDB[T]) Update(id string, value T) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if id == "" {
		return fmt.Errorf("empty entry id")
	}
	if _, exists := db.data[id]; !exists {
		return ErrRecordNotFound
	}

	oldData := maps.Clone(db.data)

	db.data[id] = value
	if db.autoSave {
		if err := db.internalSave(); err != nil {
			db.data = oldData
			return fmt.Errorf("failed to save db: %v", err)
		}
	}
	return nil
}

// Delete removes a record from the database
// Returns error if record doesn't exist
func (db *JsonDB[T]) Delete(id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.data[id]; !exists {
		return ErrRecordNotFound
	}

	oldData := maps.Clone(db.data)

	delete(db.data, id)
	if db.autoSave {
		if err := db.internalSave(); err != nil {
			db.data = oldData
			return fmt.Errorf("failed to save db: %v", err)
		}
	}
	return nil
}

// GetAll returns a copy of all records in the database
func (db *JsonDB[T]) GetAll() (map[string]T, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result := make(map[string]T, len(db.data))
	maps.Copy(result, db.data)
	return result, nil
}

func writeToFile(data []byte, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(dir, "tmp-")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return err
	}

	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpName, path); err != nil {
		return err
	}

	return nil
}
