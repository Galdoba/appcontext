package logmanager

import (
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	logg := New(WithLevel(LevelTrace))
	// logg.AddHandler(&Stderr)
	logg.AddHandler(NewHandler(os.Stderr, LevelInfo,
		NewTextFormatter(
			WithTimePrecision(3),
			WithTimerPrecision(3),
			WithLevelTag(true),
			WithColor(true),
		)))
	logg.Emergencyf("this is %v", "EMERGENCY")
	logg.Alertf("this is %v", "ALERT")
	logg.Criticalf("this is %v", "CRITICAL")
	logg.Errorf("this is %v", "ERROR")
	logg.Warnf("this is %v", "WARN")
	logg.Noticef("this is %v", "NOTICE")
	logg.Infof("this is %v", "INFO")
	logg.Debugf("this is %v", "DEBUG")
	logg.Tracef("this is %v", "TRACE")
	logg.Fatalf("this is %v", "FATAL")
}
