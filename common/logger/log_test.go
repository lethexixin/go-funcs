package logger

import (
	"testing"
)

import (
	"go.uber.org/zap/zapcore"
)

func TestLogConsole(t *testing.T) {
	_ = SetLogger(Env(LogReleaseEnv), Level(zapcore.DebugLevel))

	Debugf("hello %s", "world !!!")
	Infof("hello %s", "world !!!")
	Errorf("hello %s", "world !!!")
	Warnf("hello %s", "world !!!")

	Debug("hi:", "name:xin")
	Info("hi:", "name:xin")
	Error("hi:", "name:xin")
	Warn("hi:", "name:xin")
}

func TestLogFile(t *testing.T) {
	_ = SetLogger(Level(zapcore.DebugLevel), FileLog(&FileLogger{
		Filename:   "test_file.log",
		MaxSize:    1,
		MaxAge:     1,
		MaxBackups: 2,
		LocalTime:  true,
		Compress:   true,
	}))

	Debugf("hello %s", "world !!!")
	Infof("hello %s", "world !!!")
	Errorf("hello %s", "world !!!")
	Warnf("hello %s", "world !!!")

	Debug("hi:", "name:xin")
	Info("hi:", "name:xin")
	Error("hi:", "name:xin")
	Warn("hi:", "name:xin")
}
