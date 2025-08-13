package xdg

import (
	"fmt"
	"os"
	"path/filepath"
)

type ProgramPaths struct {
	programName string
}

func New(programName string) *ProgramPaths {
	return &ProgramPaths{programName: programName}
}

// Базовые функции для получения корневых директорий XDG

func (p *ProgramPaths) configHome() string {
	return filepath.Join(home(), ".config")
}

func (p *ProgramPaths) dataHome() string {
	return filepath.Join(home(), ".local", "share")
}

func (p *ProgramPaths) cacheHome() string {
	return filepath.Join(home(), ".cache")
}

func (p *ProgramPaths) stateHome() string {
	return filepath.Join(home(), ".local", "state")
}

// func (p *ProgramPaths) runtimeDir() string {
// 	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
// 		return dir
// 	}
// 	return filepath.Join("/run/user", strconv.Itoa(os.Getuid()))
// }

func home() string {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("home not found: %v", err))
	}
	return h
}

// Основные пути приложения

// ConfigDir возвращает путь к директории конфигурации
// Пример: /home/user/.config/myapp
func (p *ProgramPaths) ConfigDir() string {
	return filepath.Join(p.configHome(), p.programName)
}

// LogDir возвращает путь к директории логов
// Пример: /home/user/.local/state/myapp/log
func (p *ProgramPaths) LogDir() string {
	return filepath.Join(p.stateHome(), p.programName, "log")
}

// CacheDir возвращает путь к операционному кэшу
// Пример: /home/user/.cache/myapp
func (p *ProgramPaths) CacheDir() string {
	return filepath.Join(p.cacheHome(), p.programName)
}

// UserProfilesDir возвращает путь к профилям пользователей
// Пример: /home/user/.local/share/myapp/profiles
func (p *ProgramPaths) UserProfilesDir() string {
	return filepath.Join(p.dataHome(), p.programName, "profiles")
}

// PersistentDataDir возвращает путь к файлам длительного хранения
// Пример: /home/user/.local/share/myapp/data
func (p *ProgramPaths) PersistentDataDir() string {
	return filepath.Join(p.dataHome(), p.programName, "data")
}

// InboxDir возвращает путь к входящим файлам для обмена
// Пример: /home/user/.local/share/myapp/inbox
func (p *ProgramPaths) InboxDir() string {
	return filepath.Join(p.dataHome(), p.programName, "inbox")
}

// OutboxDir возвращает путь к исходящим файлам для обмена
// Пример: /home/user/.local/share/myapp/outbox
func (p *ProgramPaths) OutboxDir() string {
	return filepath.Join(p.dataHome(), p.programName, "outbox")
}

// Дополнительные пути

// StateDir возвращает общую директорию состояния программы
// Пример: /home/user/.local/state/myapp
func (p *ProgramPaths) StateDir() string {
	return filepath.Join(p.stateHome(), p.programName)
}

// PluginsDir возвращает путь к плагинам
// Пример: /home/user/.local/share/myapp/plugins
func (p *ProgramPaths) PluginsDir() string {
	return filepath.Join(p.dataHome(), p.programName, "plugins")
}

// ThemesDir возвращает путь к темам оформления
// Пример: /home/user/.local/share/myapp/themes
func (p *ProgramPaths) ThemesDir() string {
	return filepath.Join(p.dataHome(), p.programName, "themes")
}

// BinDir возвращает путь к исполняемым файлам программы
// Пример: /home/user/.local/share/myapp/bin
func (p *ProgramPaths) BinDir() string {
	return filepath.Join(p.dataHome(), p.programName, "bin")
}

// BackupsDir возвращает путь к бэкапам
// Пример: /home/user/.local/share/myapp/backups
func (p *ProgramPaths) BackupsDir() string {
	return filepath.Join(p.dataHome(), p.programName, "backups")
}

// RuntimeDir возвращает путь для временных файлов сессии
// Пример: /run/user/1000/myapp
// func (p *ProgramPaths) RuntimeDir() string {
// 	return filepath.Join(p.runtimeDir(), p.programName)
// }

//стандартные пути к файлам

// ConfigFile - возвращат стандартный путь конфига
func (p *ProgramPaths) ConfigFile() string {
	return filepath.Join(p.ConfigDir(), p.programName+".conf")
}

// LogFile  - возвращает стандартный путь логфайла
func (p *ProgramPaths) LogFile() string {
	return filepath.Join(p.LogDir(), p.programName+".log")
}
