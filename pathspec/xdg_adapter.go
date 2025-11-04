package pathspec

import (
	"github.com/Galdoba/appcontext/xdg"
)

// xdgAdapter adapts xdg package options to PathOptions
type xdgAdapter struct {
	appName       string
	projectGroup  string
	baseDir       BaseDirType
	groupcategory string
	subcategory   PathSubcategory
	name          string
	pathType      PathType
}

func newXDGAdapter(path Path) *xdgAdapter {
	return &xdgAdapter{
		appName:       path.AppName,
		baseDir:       path.BaseDir,
		groupcategory: path.Groupcategory,
		subcategory:   path.Subcategory,
		name:          path.Name,
		pathType:      path.PathType,
	}
}

func (a *xdgAdapter) toXDGOptions() []xdg.PathOption {
	var opts []xdg.PathOption

	opts = append(opts, a.baseDirToXDOption())

	if a.groupcategory != "" {
		opts = append(opts, xdg.WithProjectGroup(a.groupcategory))
	}

	if a.appName != "" {
		opts = append(opts, xdg.WithProgramName(a.appName))
	}

	if a.subcategory != "" {
		opts = append(opts, xdg.WithSubDir([]string{string(a.subcategory)}))
	}

	if a.pathType == FileType {
		opts = append(opts, xdg.WithFileName(a.name))
	} else {
		// Для директорий добавляем имя как последний элемент пути
		if a.subcategory != "" {
			// Если есть подкатегория, добавляем имя как дополнительную поддиректорию
			opts = append(opts, xdg.WithSubDir([]string{string(a.subcategory), a.name}))
		} else {
			opts = append(opts, xdg.WithSubDir([]string{a.name}))
		}
	}

	return opts
}

func (a *xdgAdapter) baseDirToXDOption() xdg.PathOption {
	switch a.baseDir {
	case Config:
		return xdg.ForConfig()
	case Data:
		return xdg.ForData()
	case Cache:
		return xdg.ForCache()
	case Runtime:
		return xdg.ForState() // XDG_STATE_HOME для runtime
	case Temp:
		return xdg.ForTemp()
	default:
		return xdg.ForData() // fallback
	}
}

// BuildPath использует xdg для построения пути
func BuildPath(path Path) string {
	adapter := newXDGAdapter(path)
	opts := adapter.toXDGOptions()
	return xdg.Location(opts...)
}
