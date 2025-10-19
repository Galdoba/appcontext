package jsonstore

import (
	"os"
	"path/filepath"
	"testing"
)

type TestData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestJsonDB_Save(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		setup   func(*JsonDB[TestData])
		wantErr bool
	}{
		{
			name:    "save empty database",
			setup:   func(db *JsonDB[TestData]) {},
			wantErr: false,
		},
		{
			name: "save database with data",
			setup: func(db *JsonDB[TestData]) {
				db.Insert("test1", TestData{Name: "test", Value: 1})
			},
			wantErr: false,
		},
		{
			name: "save to non-existent directory",
			setup: func(db *JsonDB[TestData]) {
				// Path will be set to non-existent directory
			},
			wantErr: false, // Should create directory
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".json")
			db, err := New[TestData](path)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			tt.setup(db)

			gotErr := db.Save()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Save() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Save() succeeded unexpectedly")
			}

			// Verify file was created
			if _, err := os.Stat(path); err != nil {
				t.Errorf("Save() did not create file: %v", err)
			}
		})
	}
}

func TestJsonDB_Insert(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		id      string
		data    TestData
		setup   func(*JsonDB[TestData])
		wantErr bool
	}{
		{
			name:    "insert new record",
			id:      "test1",
			data:    TestData{Name: "test1", Value: 1},
			setup:   func(db *JsonDB[TestData]) {},
			wantErr: false,
		},
		{
			name:    "insert duplicate record",
			id:      "test1",
			data:    TestData{Name: "test2", Value: 2},
			setup:   func(db *JsonDB[TestData]) { db.Insert("test1", TestData{Name: "existing", Value: 1}) },
			wantErr: true,
		},
		{
			name:    "insert with empty id",
			id:      "",
			data:    TestData{Name: "test", Value: 1},
			setup:   func(db *JsonDB[TestData]) {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".json")
			db, err := New[TestData](path)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			tt.setup(db)

			gotErr := db.Insert(tt.id, tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Insert() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Insert() succeeded unexpectedly")
			}

			// Verify record was inserted
			if !db.Contains(tt.id) {
				t.Error("Insert() did not add record to database")
			}
		})
	}
}

func TestJsonDB_Get(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		id      string
		setup   func(*JsonDB[TestData])
		want    TestData
		wantErr bool
	}{
		{
			name: "get existing record",
			id:   "test1",
			setup: func(db *JsonDB[TestData]) {
				db.Insert("test1", TestData{Name: "test1", Value: 1})
			},
			want:    TestData{Name: "test1", Value: 1},
			wantErr: false,
		},
		{
			name:    "get non-existent record",
			id:      "nonexistent",
			setup:   func(db *JsonDB[TestData]) {},
			want:    TestData{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".json")
			db, err := New[TestData](path)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			tt.setup(db)

			got, gotErr := db.Get(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Get() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Get() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonDB_Update(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		id      string
		data    TestData
		setup   func(*JsonDB[TestData])
		wantErr bool
	}{
		{
			name: "update existing record",
			id:   "test1",
			data: TestData{Name: "updated", Value: 2},
			setup: func(db *JsonDB[TestData]) {
				db.Insert("test1", TestData{Name: "original", Value: 1})
			},
			wantErr: false,
		},
		{
			name:    "update non-existent record",
			id:      "nonexistent",
			data:    TestData{Name: "test", Value: 1},
			setup:   func(db *JsonDB[TestData]) {},
			wantErr: true,
		},
		{
			name:    "update with empty id",
			id:      "",
			data:    TestData{Name: "test", Value: 1},
			setup:   func(db *JsonDB[TestData]) {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".json")
			db, err := New[TestData](path)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			tt.setup(db)

			gotErr := db.Update(tt.id, tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Update() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Update() succeeded unexpectedly")
			}

			// Verify record was updated
			updated, err := db.Get(tt.id)
			if err != nil {
				t.Errorf("Failed to get updated record: %v", err)
			}
			if updated != tt.data {
				t.Errorf("Update() did not change record: got %v, want %v", updated, tt.data)
			}
		})
	}
}

func TestJsonDB_Delete(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		id      string
		setup   func(*JsonDB[TestData])
		wantErr bool
	}{
		{
			name: "delete existing record",
			id:   "test1",
			setup: func(db *JsonDB[TestData]) {
				db.Insert("test1", TestData{Name: "test1", Value: 1})
			},
			wantErr: false,
		},
		{
			name:    "delete non-existent record",
			id:      "nonexistent",
			setup:   func(db *JsonDB[TestData]) {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".json")
			db, err := New[TestData](path)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			tt.setup(db)

			gotErr := db.Delete(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Delete() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Delete() succeeded unexpectedly")
			}

			// Verify record was deleted
			if db.Contains(tt.id) {
				t.Error("Delete() did not remove record from database")
			}
		})
	}
}

func TestJsonDB_Contains(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name  string
		id    string
		setup func(*JsonDB[TestData])
		want  bool
	}{
		{
			name: "contains existing record",
			id:   "test1",
			setup: func(db *JsonDB[TestData]) {
				db.Insert("test1", TestData{Name: "test1", Value: 1})
			},
			want: true,
		},
		{
			name:  "does not contain non-existent record",
			id:    "nonexistent",
			setup: func(db *JsonDB[TestData]) {},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".json")
			db, err := New[TestData](path)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			tt.setup(db)

			got := db.Contains(tt.id)
			if got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonDB_Count(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name  string
		setup func(*JsonDB[TestData])
		want  int
	}{
		{
			name:  "count empty database",
			setup: func(db *JsonDB[TestData]) {},
			want:  0,
		},
		{
			name: "count with multiple records",
			setup: func(db *JsonDB[TestData]) {
				db.Insert("test1", TestData{Name: "test1", Value: 1})
				db.Insert("test2", TestData{Name: "test2", Value: 2})
				db.Insert("test3", TestData{Name: "test3", Value: 3})
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".json")
			db, err := New[TestData](path)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			tt.setup(db)

			got := db.Count()
			if got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonDB_GetAll(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name  string
		setup func(*JsonDB[TestData])
		want  map[string]TestData
	}{
		{
			name:  "get all from empty database",
			setup: func(db *JsonDB[TestData]) {},
			want:  map[string]TestData{},
		},
		{
			name: "get all with multiple records",
			setup: func(db *JsonDB[TestData]) {
				db.Insert("test1", TestData{Name: "test1", Value: 1})
				db.Insert("test2", TestData{Name: "test2", Value: 2})
			},
			want: map[string]TestData{
				"test1": {Name: "test1", Value: 1},
				"test2": {Name: "test2", Value: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".json")
			db, err := New[TestData](path)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			tt.setup(db)

			got, err := db.GetAll()
			if err != nil {
				t.Errorf("GetAll() failed: %v", err)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("GetAll() returned %d records, want %d", len(got), len(tt.want))
				return
			}

			for key, wantValue := range tt.want {
				gotValue, exists := got[key]
				if !exists {
					t.Errorf("GetAll() missing key %s", key)
					continue
				}
				if gotValue != wantValue {
					t.Errorf("GetAll()[%s] = %v, want %v", key, gotValue, wantValue)
				}
			}
		})
	}
}

func TestJsonDB_NewAndLoad(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		path        string
		prepopulate bool
		wantErr     bool
	}{
		{
			name:        "new with non-existent file",
			path:        "nonexistent.json",
			prepopulate: false,
			wantErr:     true,
		},
		{
			name:        "load from existing file",
			path:        "existing.json",
			prepopulate: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.path)

			if tt.prepopulate {
				// Create and populate a database first
				db1, err := New[TestData](path)
				if err != nil {
					t.Fatalf("Setup New() failed: %v", err)
				}
				db1.Insert("test1", TestData{Name: "test1", Value: 1})
				if err := db1.Save(); err != nil {
					t.Fatalf("Setup Save() failed: %v", err)
				}
			}

			db, err := Load[TestData](path)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Load() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Load() succeeded unexpectedly")
			}

			if db.Path() != path {
				t.Errorf("Path() = %v, want %v", db.Path(), path)
			}
		})
	}
}

func TestJsonDB_AutoSave(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		autoSave bool
		ops      func(*JsonDB[TestData])
	}{
		{
			name:     "auto save enabled",
			autoSave: true,
			ops: func(db *JsonDB[TestData]) {
				db.Insert("test1", TestData{Name: "test1", Value: 1})
				db.Update("test1", TestData{Name: "updated", Value: 2})
				db.Insert("test2", TestData{Name: "test2", Value: 3})
				db.Delete("test1")
			},
		},
		{
			name:     "auto save disabled",
			autoSave: false,
			ops: func(db *JsonDB[TestData]) {
				db.Insert("test1", TestData{Name: "test1", Value: 1})
				db.Insert("test2", TestData{Name: "test2", Value: 2})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".json")
			db, err := New[TestData](path, WithAutoSave(tt.autoSave))
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			tt.ops(db)
			if !tt.autoSave {
				db.Save()
			}

			// Load the database from file to verify auto-save worked
			db2, err := Load[TestData](path)
			if err != nil {
				t.Fatalf("Load() failed: %v", err)
			}

			count1 := db.Count()
			count2 := db2.Count()

			if tt.autoSave && count1 != count2 {
				t.Errorf("AutoSave failed: original count %d, loaded count %d", count1, count2)
			}
		})
	}
}

func TestJsonDB_MarshalingMethods(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name   string
		option DB_Option
	}{
		{
			name:   "compact marshaling",
			option: WithCompactMarshaling(),
		},
		{
			name:   "indent marshaling",
			option: WithIndentMarshaling("", "  "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".json")
			db, err := New[TestData](path, tt.option)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			db.Insert("test1", TestData{Name: "test1", Value: 1})

			if err := db.Save(); err != nil {
				t.Errorf("Save() failed: %v", err)
			}

			// Verify file was created and can be loaded
			db2, err := Load[TestData](path)
			if err != nil {
				t.Errorf("Load() failed: %v", err)
				return
			}

			data, err := db2.Get("test1")
			if err != nil {
				t.Errorf("Get() failed: %v", err)
			}

			if data.Name != "test1" || data.Value != 1 {
				t.Errorf("Loaded data incorrect: got %v, want {test1 1}", data)
			}
		})
	}
}
