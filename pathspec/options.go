package pathspec

// PathOption defines a functional option for modifying Path
// [ai generated commentary]
type PathOption func(*Path)

// NewCustomPath creates a customized Path based on template and options
// [ai generated commentary]
func NewCustomPath(template Path, options ...PathOption) Path {
	result := template
	for _, option := range options {
		option(&result)
	}
	return result
}

func WithAppName(appName string) PathOption {
	return func(p *Path) {
		p.AppName = appName
	}
}

func WithName(name string) PathOption {
	return func(p *Path) {
		p.Name = name
	}
}

func WithBaseDir(baseDir BaseDirType) PathOption {
	return func(p *Path) {
		p.BaseDir = baseDir
	}
}

func WithSubcategory(subcategory PathSubcategory) PathOption {
	return func(p *Path) {
		p.Subcategory = subcategory
	}
}

func WithNameGroup(nameGroup string) PathOption {
	return func(p *Path) {
		p.Groupcategory = nameGroup
	}
}

func WithPathType(pathType PathType) PathOption {
	return func(p *Path) {
		p.PathType = pathType
	}
}

func WithCategory(category PathCategory) PathOption {
	return func(p *Path) {
		p.Category = category
	}
}

func WithPriority(priority PathPriority) PathOption {
	return func(p *Path) {
		p.Priority = priority
	}
}

func WithDescription(description string) PathOption {
	return func(p *Path) {
		p.Description = description
	}
}

func WithPattern(pattern string) PathOption {
	return func(p *Path) {
		p.Pattern = pattern
	}
}

func WithDefaultPerm(perm uint32) PathOption {
	return func(p *Path) {
		p.DefaultPerm = perm
	}
}

func WithOwnerOnly(ownerOnly bool) PathOption {
	return func(p *Path) {
		p.OwnerOnly = ownerOnly
	}
}

func WithIsMandatory(mandatory bool) PathOption {
	return func(p *Path) {
		p.IsMandatory = mandatory
	}
}

func WithIsAutoCreated(autoCreated bool) PathOption {
	return func(p *Path) {
		p.IsAutoCreated = autoCreated
	}
}

func WithIsBackedUp(backedUp bool) PathOption {
	return func(p *Path) {
		p.IsBackedUp = backedUp
	}
}

func WithIsVersioned(versioned bool) PathOption {
	return func(p *Path) {
		p.IsVersioned = versioned
	}
}

func WithIsCompressible(compressible bool) PathOption {
	return func(p *Path) {
		p.IsCompressible = compressible
	}
}

func WithMaxSize(maxSize uint64) PathOption {
	return func(p *Path) {
		p.MaxSize = maxSize
	}
}

func WithFormat(format string) PathOption {
	return func(p *Path) {
		p.Format = format
	}
}

func WithMaxChildren(maxChildren uint32) PathOption {
	return func(p *Path) {
		p.MaxChildren = maxChildren
	}
}

func WithHasSubdirs(hasSubdirs bool) PathOption {
	return func(p *Path) {
		p.HasSubdirs = hasSubdirs
	}
}

func WithRetentionDays(retentionDays uint16) PathOption {
	return func(p *Path) {
		p.RetentionDays = retentionDays
	}
}

func WithCleanupAge(cleanupAge uint16) PathOption {
	return func(p *Path) {
		p.CleanupAge = cleanupAge
	}
}
