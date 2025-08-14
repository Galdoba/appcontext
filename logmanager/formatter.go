package logmanager

import (
	"fmt"
	"strings"
	"time"

	"github.com/Galdoba/appcontext/logmanager/colorizer"
)

type Formatter interface {
	Format(LogLevel, string, ...any) string
}

type PlainTextFormatter struct{}

func (pt *PlainTextFormatter) Format(level LogLevel, format string, args ...any) string {
	timestamp := time.Now().Format(time.DateTime + ".000")
	return fmt.Sprintf("[%s] [%s]: %s", timestamp, levelToString(level), fmt.Sprintf(format, args...))
}

type PlainTextTimerFormatter struct {
	start time.Time
}

func (pt *PlainTextTimerFormatter) Format(level LogLevel, format string, args ...any) string {
	ms := time.Since(pt.start).Milliseconds()
	timestamp := msToTimestamp(ms)
	return fmt.Sprintf("[%s] [%s]: %s", timestamp, levelToString(level), fmt.Sprintf(format, args...))
}

func msToTimestamp(ms int64) string {
	seconds := ms / 1000
	milli := ms - (seconds * 1000)
	return fmt.Sprintf("%v.%v", headZero(seconds, 4), headZero(milli, 3))
}

func headZero(n int64, z int) string {
	text := fmt.Sprintf("%v", n)
	for len(text) < z {
		text = "0" + text
	}
	return text
}

type ByTypeColorTextFormatter struct {
	colors *colorizer.Colorizer
}

func (col *ByTypeColorTextFormatter) Format(level LogLevel, format string, args ...any) string {
	prefixFgColorKey := colorizer.NewKey(colorizer.FG_KEY, levelToString(level))
	prefixBgColorKey := colorizer.NewKey(colorizer.BG_KEY, levelToString(level))
	timestamp := time.Now().Format(time.DateTime + ".000")

	text := fmt.Sprintf("[%v] [%v]: ", col.colors.ColorizeByKeys(timestamp, prefixFgColorKey, prefixBgColorKey), col.colors.ColorizeByKeys(levelToString_5(level), prefixFgColorKey, prefixBgColorKey))
	vals := []any{}
	for _, arg := range args {
		vals = append(vals, col.colors.ColorizeByType(arg))
	}

	return text + fmt.Sprintf(format, vals...)
}

func levelToString_5(level LogLevel) string {
	switch level {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info "
	case LevelNotice:
		return "notic"
	case LevelWarn:
		return "warn "
	case LevelError:
		return "error"
	case LevelCritical:
		return "critc"
	case LevelAlert:
		return "alert"
	case LevelFatal:
		return "fatal"
	case LevelEmergency:
		return "emerg"
	}
	return "unknw"
}

type TextFormatter struct {
	time       int
	startTime  time.Time
	tick       int
	level      int
	levelBlock bool
	dir        bool
	file       bool
	line       bool
	color      *colorizer.Colorizer
}

func NewTextFormatter(options ...TextFormatterOption) *TextFormatter {
	tf := TextFormatter{
		time:       -1,
		startTime:  time.Time{},
		tick:       -1,
		level:      -1,
		levelBlock: true,
		dir:        false,
		file:       false,
		line:       false,
		color:      nil,
	}
	for _, modify := range options {
		modify(&tf)
	}
	return &tf
}

type TextFormatterOption func(*TextFormatter)

func WithTimePrecision(tp int) TextFormatterOption {
	return func(tf *TextFormatter) {
		tf.time = tp
	}
}

func WithColor(enabled bool) TextFormatterOption {
	return func(tf *TextFormatter) {
		if enabled {
			tf.color = colorizer.DefaultColorizer()
		}
	}
}

func WithTimerPrecision(tp int) TextFormatterOption {
	return func(tf *TextFormatter) {
		tf.tick = tp
	}
}

func WithLevelTag(enabled bool) TextFormatterOption {
	return func(tf *TextFormatter) {
		tf.levelBlock = enabled
	}
}

func (tf *TextFormatter) Format(level LogLevel, format string, args ...any) string {
	timeBlock := timeStamp(tf.time)
	timer := ""
	if tf.tick > 0 {
		timer = fmt.Sprintf("%v", float64(time.Since(tf.startTime).Milliseconds())/1000)
	}
	lev := ""
	if tf.levelBlock {
		lev = levelToString(level)
	}
	blocks := []string{}
	for _, block := range []string{timeBlock, timer, lev} {
		if block != "" {
			blocks = append(blocks, block)
		}
	}
	switch tf.color {
	case nil:
		s := ""
		for _, bl := range blocks {
			s += fmt.Sprintf("[%v] ", bl)
		}
		s = strings.TrimSuffix(s, " ")
		if s != "" {
			s += ": "
		}
		s += fmt.Sprintf(format, args...)
		return s

	default:
		prefixFgColorKey := colorizer.NewKey(colorizer.FG_KEY, levelToString(level))
		prefixBgColorKey := colorizer.NewKey(colorizer.BG_KEY, levelToString(level))
		for i := range blocks {
			blocks[i] = tf.color.ColorizeByKeys(blocks[i], prefixFgColorKey, prefixBgColorKey)
		}

		vals := []any{}
		for _, arg := range args {
			vals = append(vals, tf.color.ColorizeByType(arg))
		}
		s := ""
		for _, bl := range blocks {
			s += fmt.Sprintf("[%v] ", bl)
		}
		s = strings.TrimSuffix(s, " ")
		if s != "" {
			s += ": "
		}
		s += fmt.Sprintf(format, vals...)
		return s
	}

}

func timeStamp(zeroes int) string {
	if zeroes > -1 {
		z := ""
		for len(z) < zeroes {
			z = "0" + z
		}
		if len(z) > 0 {
			z = "." + z
		}
		return time.Now().Format(time.DateTime + z)
	}
	return ""
}
