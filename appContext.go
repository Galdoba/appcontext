package appcontext

import (
	"fmt"

	"github.com/Galdoba/gogacon"
	"github.com/Galdoba/golog"
	"github.com/Galdoba/xdgpaths"
)

type AppContext struct {
	AppName string
	Path    *xdgpaths.ProgramPaths
	Config  *gogacon.ConfigManager
	Logger  *golog.Logger
	err     error
}

func New(appName string, optional ...OptionalContext) *AppContext {
	actx := AppContext{AppName: appName}
	actx.Path = xdgpaths.New(actx.AppName)
	for _, modify := range optional {
		modify(&actx)
	}
	return &actx
}

type OptionalContext func(*AppContext)

func WithConfig(defaultConfig gogacon.Serializer) OptionalContext {
	return func(ac *AppContext) {
		cm, err := gogacon.NewConfigManager(gogacon.Defaults{
			AppName:             ac.AppName,
			DefaultConfigValues: defaultConfig,
		})
		if err != nil {
			ac.err = err
			return
		}
		ac.Config = cm
	}
}

func WithLogger(log *golog.Logger) OptionalContext {
	return func(ac *AppContext) {
		ac.Logger = log
	}
}

func (actx *AppContext) LoadConfig(cfg gogacon.Serializer, paths ...string) error {
	collectedErrors := []error{}
	paths = append(paths, actx.Path.ConfigDir())
	for _, path := range paths {
		if err := actx.Config.LoadConfig(path, cfg); err != nil {
			collectedErrors = append(collectedErrors, fmt.Errorf("failed to load from %v: %v", err))
		} else {
			return nil
		}
	}
	return fmt.Errorf("colected errors: %v", collectedErrors)
}
