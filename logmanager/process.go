package logmanager

import (
	"fmt"
	"os"
)

func (l *Logger) log(level LogLevel, msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	mustShutdown := false
	for _, handler := range l.handlers {
		if level < handler.minLevel {
			continue
		}
		if level >= LevelFatal {
			mustShutdown = true
		}

		msg := handler.formatter.Format(level, msg, args...)
		switch handler.destination {
		case Stderr:
			fmt.Fprintln(os.Stderr, msg)
		case Stdout:
			fmt.Fprintln(os.Stdout, msg)
		default:
			f, err := os.OpenFile(handler.destination, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				continue
			}
			defer f.Close()
			if _, err := f.Write([]byte(msg + "\n")); err != nil {
				continue
			}
		}

	}
	if mustShutdown {
		os.Exit(1)
	}

}
