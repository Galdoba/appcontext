package logmanager

import (
	"os"
)

const (
	Stderr = "stderr"
	Stdout = "stdout"
)

type MessageHandler struct {
	formatter   Formatter
	minLevel    LogLevel
	destination string
	descriptor  *os.File
}

var StderrHandler MessageHandler = MessageHandler{
	destination: Stderr,
	descriptor:  os.Stderr,
	formatter:   &PlainTextFormatter{},
	minLevel:    LevelInfo,
}

func (l *Logger) AddHandler(handler *MessageHandler) {
	l.handlers = append(l.handlers, handler)
}

func NewHandler(destination string, level LogLevel, formatter Formatter) *MessageHandler {
	switch destination {
	case Stderr:
		return &MessageHandler{
			destination: destination,
			descriptor:  os.Stderr,
			formatter:   formatter,
			minLevel:    level,
		}
	case Stdout:
		return &MessageHandler{
			destination: destination,
			descriptor:  os.Stdout,
			formatter:   formatter,
			minLevel:    level,
		}
	default:
		return &MessageHandler{
			formatter:   formatter,
			minLevel:    level,
			destination: destination,
		}

	}
}
