package jsonstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
)

const (
	LibVersion = "1.0.0"
)

// Predefined errors for common operations
var ErrRecordNotFound = errors.New("record not found")
var ErrRecordExist = errors.New("record already exist")

// JsonDB represents a JSON-based database for storing items of type T.
// It is safe for concurrent use by multiple goroutines.
type JsonDB[T any] struct {
	mu       sync.RWMutex
	data     map[string]*T
	path     string
	autoSave bool
}

// options holds configuration options for JsonDB.
type options struct {
	autoSave bool
}

// DB_Option defines a function type for configuring JsonDB options.
type DB_Option func(*options)

// WithAutoSave returns a DB_Option that sets the autoSave flag.
// When enabled, the database will automatically save to file after every modification.
func WithAutoSave(autoSave bool) DB_Option {
	return func(o *options) {
		o.autoSave = autoSave
	}
}

// New creates a new JsonDB instance. If the file exists, it loads the data from the file.
// If the file does not exist, it creates an empty database.
// opts are optional parameters to configure the DB.
func New[T any](path string, opts ...DB_Option) (*JsonDB[T], error) {
	db := &JsonDB[T]{
		path: path,
		data: make(map[string]*T),
	}
	optionSet := options{}
	for _, modify := range opts {
		modify(&optionSet)
	}

	db.autoSave = optionSet.autoSave

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

// Path returns the file path of the database.
func (db *JsonDB[T]) Path() string {
	return db.path
}

// Load creates a JsonDB instance from an existing file.
// Returns an error if the file doesn't exist or cannot be read.
// opts are optional parameters to configure the DB.
func Load[T any](path string, opts ...DB_Option) (*JsonDB[T], error) {
	db := &JsonDB[T]{
		path: path,
		data: make(map[string]*T),
	}

	// Apply options
	optionSet := options{}
	for _, modify := range opts {
		modify(&optionSet)
	}
	db.autoSave = optionSet.autoSave

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &db.data); err != nil {
		return nil, err
	}

	return db, nil
}

// Save writes the current state of the database to file atomically.
// It uses a temporary file and then renames it to the target file to ensure atomicity.
func (db *JsonDB[T]) Save() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.save()
}

// save writes the database to file with custom formatting without locking.
// It is the caller's responsibility to ensure proper locking.
func (db *JsonDB[T]) save() error {
	var buf bytes.Buffer
	buf.WriteString("{\n")

	keys := make([]string, 0, len(db.data))
	for k := range db.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, key := range keys {
		keyJSON, err := json.Marshal(key)
		if err != nil {
			return fmt.Errorf("failed to marshal entry '%v': %v", key, err)
		}
		buf.WriteString("  ")
		buf.Write(keyJSON)
		buf.WriteString(": ")

		valJSON, err := json.Marshal(db.data[key])
		if err != nil {
			return err
		}
		buf.Write(valJSON)

		if i < len(keys)-1 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
	}

	buf.WriteString("}")

	// Write to a temporary file first to ensure atomicity
	tmpPath := db.path + ".tmp"
	if err := os.WriteFile(tmpPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %v", err)
	}

	if err := os.Rename(tmpPath, db.path); err != nil {
		return fmt.Errorf("failed to rename temp file: %v", err)
	}
	return nil
}

// autoSaveCheck triggers save if autoSave is enabled.
// This method should be called only when the mutex is already locked.
func (db *JsonDB[T]) autoSaveCheck() error {
	if db.autoSave {
		return db.save()
	}
	return nil
}

// Insert adds a new record to the database.
// Returns an error if the id is empty or if a record with the same id already exists.
func (db *JsonDB[T]) Insert(id string, value *T) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if id == "" {
		return fmt.Errorf("empty entry id")
	}
	if _, exists := db.data[id]; exists {
		return ErrRecordExist
	}

	db.data[id] = value
	return db.autoSaveCheck()
}

// Contains checks if a record with the given id exists in the database.
func (db *JsonDB[T]) Contains(id string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	_, exists := db.data[id]
	return exists
}

// Count returns the number of entries in the database.
func (db *JsonDB[T]) Count() int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return len(db.data)
}

// Get retrieves a record from the database by ID.
// Returns an error if the record is not found.
func (db *JsonDB[T]) Get(id string) (*T, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	val, exists := db.data[id]
	if !exists {
		return nil, ErrRecordNotFound
	}
	return val, nil
}

// Update modifies an existing record in the database.
// Returns an error if the id is empty or the record does not exist.
func (db *JsonDB[T]) Update(id string, value *T) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if id == "" {
		return fmt.Errorf("empty entry id")
	}
	if _, exists := db.data[id]; !exists {
		return ErrRecordNotFound
	}

	db.data[id] = value
	return db.autoSaveCheck()
}

// Delete removes a record from the database.
// Returns an error if the record does not exist.
func (db *JsonDB[T]) Delete(id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.data[id]; !exists {
		return ErrRecordNotFound
	}

	delete(db.data, id)
	return db.autoSaveCheck()
}

// GetAll returns a copy of all records in the database.
// Note that the values are pointers to the actual data, so modifying them will affect the database.
func (db *JsonDB[T]) GetAll() (map[string]*T, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result := make(map[string]*T, len(db.data))
	for k, v := range db.data {
		result[k] = v
	}
	return result, nil
}
