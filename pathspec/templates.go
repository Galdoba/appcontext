package pathspec

// JSONStorageTemplate defines a template for JSON data storage files
// [ai generated commentary]
var JSONStorageTemplate = Path{
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
	Subcategory:    SubcategoryStorage,
}

// LogFileTemplate defines a template for primary application log files
// [ai generated commentary]
var LogFileTemplate = Path{
	BaseDir:       Runtime,
	PathType:      FileType,
	Category:      CategoryRuntime,
	Priority:      PriorityMedium,
	DefaultPerm:   0644,
	OwnerOnly:     true,
	IsAutoCreated: true,
	IsBackedUp:    false,
	Format:        "text",
	MaxSize:       10 * 1024 * 1024,
	RetentionDays: 30,
	Subcategory:   SubcategoryLogs,
}

// ConfigFileTemplate defines a template for primary configuration files
// [ai generated commentary]
var ConfigFileTemplate = Path{
	BaseDir:     Config,
	PathType:    FileType,
	Category:    CategoryConfig,
	Priority:    PriorityCritical,
	DefaultPerm: 0644,
	OwnerOnly:   true,
	IsMandatory: true,
	IsBackedUp:  true,
	IsVersioned: true,
	Format:      "toml",
	Subcategory: SubcategoryConfig,
}

// ProcessStateTemplate defines a template for process state directories
// [ai generated commentary]
var ProcessStateTemplate = Path{
	BaseDir:       Data,
	PathType:      DirectoryType,
	Category:      CategoryData,
	Priority:      PriorityHigh,
	DefaultPerm:   0755,
	OwnerOnly:     true,
	IsAutoCreated: true,
	IsBackedUp:    true,
	Subcategory:   SubcategoryProcesses,
	HasSubdirs:    true,
	MaxChildren:   1000,
	RetentionDays: 90,
}

// ProjectsTemplate defines a template for user project directories
// [ai generated commentary]
var ProjectsTemplate = Path{
	BaseDir:       Data,
	PathType:      DirectoryType,
	Category:      CategoryData,
	Priority:      PriorityHigh,
	DefaultPerm:   0755,
	OwnerOnly:     true,
	IsAutoCreated: true,
	IsBackedUp:    true,
	Subcategory:   SubcategoryProjects,
	HasSubdirs:    true,
	MaxChildren:   100,
	RetentionDays: 0, // Permanent storage
}

// StatsTemplate defines a template for statistics files
// [ai generated commentary]
var StatsTemplate = Path{
	BaseDir:       Runtime,
	PathType:      FileType,
	Category:      CategoryRuntime,
	Priority:      PriorityMedium,
	DefaultPerm:   0644,
	OwnerOnly:     true,
	IsAutoCreated: true,
	IsBackedUp:    false,
	Subcategory:   SubcategoryStats,
	Format:        "json",
	MaxSize:       5 * 1024 * 1024, // 5MB
	RetentionDays: 365,
}

// BackupStorageTemplate defines a template for backup storage directories
// [ai generated commentary]
var BackupStorageTemplate = Path{
	BaseDir:        Data,
	PathType:       DirectoryType,
	Category:       CategoryData,
	Priority:       PriorityMedium,
	DefaultPerm:    0700,
	OwnerOnly:      true,
	IsAutoCreated:  false,
	IsBackedUp:     false, // Backups shouldn't be backed up recursively
	IsCompressible: true,
	Subcategory:    SubcategoryBackups,
	HasSubdirs:     true,
	MaxChildren:    50, // Limit number of backup sets
	RetentionDays:  30, // Keep backups for 30 days
	CleanupAge:     90, // Start cleanup after 90 days
}

// UploadCacheTemplate defines a template for uploaded file cache
// [ai generated commentary]
var UploadCacheTemplate = Path{
	BaseDir:        Cache,
	PathType:       DirectoryType,
	Category:       CategoryCache,
	Priority:       PriorityLow,
	DefaultPerm:    0750,
	OwnerOnly:      true,
	IsAutoCreated:  true,
	IsBackedUp:     false,
	IsCompressible: false,
	Subcategory:    SubcategoryUploads,
	HasSubdirs:     true,
	MaxChildren:    1000, // Limit number of uploaded files
	RetentionDays:  7,    // Keep uploads for 7 days
	CleanupAge:     30,   // Aggressive cleanup after 30 days
}

// ExportTemplate defines a template for export file directories
// [ai generated commentary]
var ExportTemplate = Path{
	BaseDir:        Data,
	PathType:       DirectoryType,
	Category:       CategoryData,
	Priority:       PriorityMedium,
	DefaultPerm:    0755,
	OwnerOnly:      false, // Exports might be shared
	IsAutoCreated:  true,
	IsBackedUp:     false, // Exports are generated data
	IsCompressible: true,
	Subcategory:    SubcategoryExports,
	HasSubdirs:     true,
	MaxChildren:    100, // Limit number of export sets
	RetentionDays:  14,  // Keep exports for 14 days
	CleanupAge:     60,  // Cleanup old exports after 60 days
}
