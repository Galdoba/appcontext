package pathspec

import (
	"testing"
)

func TestNewCustomPath(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		template Path
		options  []PathOption
		want     Path
	}{
		{
			name:     "customize with name and description",
			template: ConfigFileTemplate,
			options: []PathOption{
				WithName("custom_config.yaml"),
				WithDescription("Custom configuration file"),
			},
			want: Path{
				BaseDir:       Config,
				PathType:      FileType,
				Category:      CategoryConfig,
				Priority:      PriorityCritical,
				DefaultPerm:   0644,
				OwnerOnly:     true,
				IsMandatory:   true,
				IsBackedUp:    true,
				IsVersioned:   true,
				Format:        "yaml",
				Subcategory:   SubcategoryConfig,
				Name:          "custom_config.yaml",
				Description:   "Custom configuration file",
			},
		},
		{
			name:     "change permissions and category",
			template: JSONStorageTemplate,
			options: []PathOption{
				WithDefaultPerm(0600),
				WithCategory(CategoryData),
				WithSubcategory(SubcategoryDatabase),
			},
			want: Path{
				BaseDir:        Data,
				PathType:       FileType,
				Category:       CategoryData,
				Priority:       PriorityHigh,
				DefaultPerm:    0600,
				OwnerOnly:      true,
				IsAutoCreated:  true,
				IsBackedUp:     true,
				IsCompressible: true,
				Format:         "json",
				RetentionDays:  0,
				Subcategory:    SubcategoryDatabase,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCustomPath(tt.template, tt.options...)

			if got.Name != tt.want.Name {
				t.Errorf("NewCustomPath() Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Description != tt.want.Description {
				t.Errorf("NewCustomPath() Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.DefaultPerm != tt.want.DefaultPerm {
				t.Errorf("NewCustomPath() DefaultPerm = %v, want %v", got.DefaultPerm, tt.want.DefaultPerm)
			}
			if got.Category != tt.want.Category {
				t.Errorf("NewCustomPath() Category = %v, want %v", got.Category, tt.want.Category)
			}
			if got.Subcategory != tt.want.Subcategory {
				t.Errorf("NewCustomPath() Subcategory = %v, want %v", got.Subcategory, tt.want.Subcategory)
			}
		})
	}
}
