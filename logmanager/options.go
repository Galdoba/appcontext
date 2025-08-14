package logmanager

type LoggerOption func(*Logger)

func WithLevel(level LogLevel) LoggerOption {
	return func(l *Logger) {
		l.level = level
	}
}

func WithHandlers(handlers ...*MessageHandler) LoggerOption {
	return func(l *Logger) {
		for _, hand := range handlers {
			l.AddHandler(hand)
		}
	}
}

// func WithOutput(writers ...io.Writer) LoggerOption {
// 	return func(l *Logger) {
// 		l.outputs = append(l.outputs, writers...)
// 	}
// }

// func WithPrefix(prefix string) LoggerOption {
// 	return func(l *Logger) {
// 		l.prefix = prefix
// 	}
// }

// func WithColor(enabled bool) LoggerOption {
// 	return func(l *Logger) {
// 		l.useColors = enabled
// 	}
// }

// func WithFlags(flags int) LoggerOption {
// 	return func(l *Logger) {
// 		l.flags = flags
// 	}
// }
