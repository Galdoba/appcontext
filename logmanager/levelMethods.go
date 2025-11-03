package logmanager

import "strings"

func (l *Logger) Tracef(msg string, args ...interface{}) {
	l.log(LevelTrace, msg, args...)
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.log(LevelDebug, msg, args...)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.log(LevelInfo, msg, args...)
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	l.log(LevelInfo, msg, args...)
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.log(LevelWarn, msg, args...)
}

func (l *Logger) Noticef(msg string, args ...interface{}) {
	l.log(LevelNotice, msg, args...)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.log(LevelError, msg, args...)
}

func (l *Logger) Criticalf(msg string, args ...interface{}) {
	l.log(LevelCritical, msg, args...)
}

func (l *Logger) Alertf(msg string, args ...interface{}) {
	l.log(LevelAlert, msg, args...)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.log(LevelFatal, msg, args...)
}
func (l *Logger) Emergencyf(msg string, args ...interface{}) {
	l.log(LevelEmergency, msg, args...)
}

var levelNames = map[LogLevel]string{
	LevelTrace:     "trace",
	LevelDebug:     "debug",
	LevelInfo:      "info",
	LevelNotice:    "notice",
	LevelWarn:      "warning",
	LevelError:     "error",
	LevelCritical:  "critical",
	LevelAlert:     "alert",
	LevelFatal:     "fatal",
	LevelEmergency: "emergency",
}

func levelToString(l LogLevel) string {
	return levelNames[l]
}

// ListLevels - convinience function list log levels and it's string representations.
func ListLevels() ([]LogLevel, []string) {
	lvls := []LogLevel{}
	asString := []string{}
	for l := LevelTrace; l <= LevelFatal; l++ {
		lvls = append(lvls, l)
		asString = append(asString, levelToString(l))
	}
	return lvls, asString
}

// func StringToLevel(s string) LogLevel {
// 	l := LevelInfo
// 	for newL, str := range levelNames {
// 		if str == s {
// 			return newL
// 		}
// 	}
// 	return l
// }

func StringToLevel(s string) LogLevel {
	s = strings.ToLower(s)
	levels, strngs := ListLevels()
	for i, level := range levels {
		if strings.HasPrefix(strngs[i], s) {
			return level
		}
	}
	return LevelInfo
}
