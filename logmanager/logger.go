package logmanager

import (
	"sync"
	"time"
)

const (
	LibVersion = "1.0.0"
)

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelCritical
	LevelAlert
	LevelEmergency
	LevelFatal
)

type Logger struct {
	mu        sync.Mutex
	level     LogLevel
	startTime time.Time
	handlers  []*MessageHandler
}

func New(options ...LoggerOption) *Logger {
	l := Logger{
		mu:        sync.Mutex{},
		level:     LevelInfo,
		startTime: time.Now(),
	}

	for _, modify := range options {
		modify(&l)
	}

	return &l
}
