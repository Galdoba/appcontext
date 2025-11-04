package pathspec

type BaseDirType uint8

const (
	Config BaseDirType = iota
	// Configuration files directory (XDG_CONFIG_HOME)
	// [ai generated commentary]
	Data
	// Data files directory (XDG_DATA_HOME)
	// [ai generated commentary]
	Cache
	// Cache files directory (XDG_CACHE_HOME)
	// [ai generated commentary]
	Runtime
	// Runtime files directory (XDG_STATE_HOME)
	// [ai generated commentary]
	Temp
	// Temporary files directory
	// [ai generated commentary]
)

type PathType uint8

const (
	FileType PathType = iota
	// Regular file
	// [ai generated commentary]
	DirectoryType
	// Directory
	// [ai generated commentary]
	SymlinkType
	// Symbolic link
	// [ai generated commentary]
)

type PathCategory uint8

const (
	CategoryConfig PathCategory = iota
	// Configuration files category
	// [ai generated commentary]
	CategoryData
	// Data files category
	// [ai generated commentary]
	CategoryCache
	// Cache files category
	// [ai generated commentary]
	CategoryRuntime
	// Runtime files category
	// [ai generated commentary]
	CategoryTemp
	// Temporary files category
	// [ai generated commentary]
)

type PathSubcategory string

const (
	SubcategoryDatabase PathSubcategory = "database"
	// Database files subcategory
	// [ai generated commentary]
	SubcategoryStorage PathSubcategory = "storage"
	// Storage files subcategory
	// [ai generated commentary]
	SubcategoryLogs PathSubcategory = "logs"
	// Log files subcategory
	// [ai generated commentary]
	SubcategoryConfig PathSubcategory = ""
	// Configuration files subcategory
	// [ai generated commentary]
	SubcategoryTemplates PathSubcategory = "templates"
	// Template files subcategory
	// [ai generated commentary]
	SubcategoryPlugins PathSubcategory = "plugins"
	// Plugin files subcategory
	// [ai generated commentary]
	SubcategoryState PathSubcategory = "state"
	// Application state subcategory
	// [ai generated commentary]
	SubcategoryResources PathSubcategory = "resources"
	// Resource files subcategory
	// [ai generated commentary]
	SubcategoryLocks PathSubcategory = "locks"
	// Lock files subcategory
	// [ai generated commentary]
	SubcategorySockets PathSubcategory = "sockets"
	// Socket files subcategory
	// [ai generated commentary]
	SubcategoryProcessing PathSubcategory = "processing"
	// Processing files subcategory
	// [ai generated commentary]
	SubcategoryThumbnails PathSubcategory = "thumbnails"
	// Thumbnail files subcategory
	// [ai generated commentary]
	SubcategoryCacheData PathSubcategory = "cache"
	// Cache data subcategory
	// [ai generated commentary]
	SubcategoryProcesses PathSubcategory = "processes"
	// Process state subcategory
	// [ai generated commentary]
	SubcategoryProjects PathSubcategory = "projects"
	// User projects subcategory
	// [ai generated commentary]
	SubcategoryStats PathSubcategory = "stats"
	// Statistics subcategory
	// [ai generated commentary]
	SubcategoryBackups PathSubcategory = "backups"
	// Backup files subcategory
	// [ai generated commentary]
	SubcategoryUploads PathSubcategory = "uploads"
	// Uploaded files subcategory
	// [ai generated commentary]
	SubcategoryExports PathSubcategory = "exports"
	// Export files subcategory
	// [ai generated commentary]
)

type PathPriority uint8

const (
	PriorityCritical PathPriority = iota
	// Critical priority for essential files
	// [ai generated commentary]
	PriorityHigh
	// High priority for important files
	// [ai generated commentary]
	PriorityMedium
	// Medium priority for standard files
	// [ai generated commentary]
	PriorityLow
	// Low priority for temporary files
	// [ai generated commentary]
)

// ValidSubcategories defines allowed subcategories for each path category
// [ai generated commentary]
var ValidSubcategories = map[PathCategory]map[PathSubcategory]bool{
	CategoryConfig: {
		SubcategoryConfig:    true,
		SubcategoryTemplates: true,
		SubcategoryPlugins:   true,
	},
	CategoryData: {
		SubcategoryDatabase:  true,
		SubcategoryStorage:   true,
		SubcategoryState:     true,
		SubcategoryResources: true,
		SubcategoryProcesses: true,
		SubcategoryProjects:  true,
		SubcategoryBackups:   true,
		SubcategoryUploads:   true,
		SubcategoryExports:   true,
	},
	CategoryCache: {
		SubcategoryCacheData:  true,
		SubcategoryThumbnails: true,
	},
	CategoryRuntime: {
		SubcategoryLogs:  true,
		SubcategoryStats: true,
	},
	CategoryTemp: {
		SubcategoryLocks:      true,
		SubcategorySockets:    true,
		SubcategoryProcessing: true,
	},
}

// Path represents a single file or directory used by the application
// [ai generated commentary]
type Path struct {
	AppName string `json:"app_name,omitempty"`
	// Application name for path construction
	// [ai generated commentary]
	Name string `json:"name"`
	// File or directory name
	// [ai generated commentary]
	BaseDir BaseDirType `json:"base_dir"`
	// Base directory type according to XDG
	// [ai generated commentary]
	Groupcategory string `json:"groupcategory,omitempty"`
	// Groupcategory - upper level category for detailed path organization
	// [ai generated commentary]
	Subcategory PathSubcategory `json:"subcategory,omitempty"`
	// Subcategory - lower level category for detailed path organization
	// [ai generated commentary]
	PathType PathType `json:"path_type"`
	// Type of filesystem object
	// [ai generated commentary]
	Category PathCategory `json:"category"`
	// Category for lifecycle management
	// [ai generated commentary]
	Priority PathPriority `json:"priority"`
	// Priority for management operations
	// [ai generated commentary]
	Description string `json:"description,omitempty"`
	// Purpose and description
	// [ai generated commentary]
	Pattern string `json:"pattern,omitempty"`
	// Name pattern for dynamic files
	// [ai generated commentary]
	DefaultPerm uint32 `json:"default_perm"`
	// Default permissions
	// [ai generated commentary]
	OwnerOnly bool `json:"owner_only"`
	// Accessible only by owner
	// [ai generated commentary]
	IsMandatory bool `json:"is_mandatory"`
	// Required for application operation
	// [ai generated commentary]
	IsAutoCreated bool `json:"is_auto_created"`
	// Automatically created during initialization
	// [ai generated commentary]
	IsBackedUp bool `json:"is_backed_up"`
	// Included in backup operations
	// [ai generated commentary]
	IsVersioned bool `json:"is_versioned"`
	// Version controlled
	// [ai generated commentary]
	IsCompressible bool `json:"is_compressible"`
	// Can be compressed during archiving
	// [ai generated commentary]
	MaxSize uint64 `json:"max_size,omitempty"`
	// Maximum file size in bytes
	// [ai generated commentary]
	Format string `json:"format,omitempty"`
	// Data format (yaml, json, binary, text)
	// [ai generated commentary]
	MaxChildren uint32 `json:"max_children,omitempty"`
	// Maximum number of files in directory
	// [ai generated commentary]
	HasSubdirs bool `json:"has_subdirs"`
	// Contains subdirectories
	// [ai generated commentary]
	RetentionDays uint16 `json:"retention_days,omitempty"`
	// Storage duration in days
	// [ai generated commentary]
	CleanupAge uint16 `json:"cleanup_age,omitempty"`
	// Age for automatic cleanup in days
	// [ai generated commentary]
}

// Layout represents complete application file structure
// [ai generated commentary]
type Layout struct {
	AppName string `json:"app_name"`
	// Application name
	// [ai generated commentary]
	AppVersion string `json:"app_version,omitempty"`
	// Application version
	// [ai generated commentary]
	ConfigPaths []Path `json:"config_paths,omitempty"`
	// Configuration file paths
	// [ai generated commentary]
	DataPaths []Path `json:"data_paths,omitempty"`
	// Data file paths
	// [ai generated commentary]
	CachePaths []Path `json:"cache_paths,omitempty"`
	// Cache file paths
	// [ai generated commentary]
	RuntimePaths []Path `json:"runtime_paths,omitempty"`
	// Runtime file paths
	// [ai generated commentary]
	TempPaths []Path `json:"temp_paths,omitempty"`
	// Temporary file paths
	// [ai generated commentary]
}
