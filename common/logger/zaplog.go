package logger

import (
	"os"
	"os/signal"
	"syscall"
)

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// this zap logger refer to dubbo-go logger.go

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	DPanic(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})

	Debugf(fmt string, args ...interface{})
	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	DPanicf(fmt string, args ...interface{})
	Panicf(fmt string, args ...interface{})
	Fatalf(fmt string, args ...interface{})
}

type EnvLogger string

type FileLogger struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

const (
	LogReleaseEnv EnvLogger = "release"
	LogDevelopEnv EnvLogger = "develop"
)

var (
	log       Logger
	zapLogger *zap.Logger

	zapLoggerConfig        = zap.NewDevelopmentConfig()
	zapLoggerEncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
)

func init() {
	zapLoggerConfig.EncoderConfig = zapLoggerEncoderConfig
	zapLogger, _ = zapLoggerConfig.Build()
	log = zapLogger.Sugar()

	// flushes buffer when redirect log to file.
	var exitSignal = make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-exitSignal
		// Sync calls the underlying Core's Sync method, flushing any buffered log entries.
		// Applications should take care to call Sync before exiting.
		err := zapLogger.Sync() // flushes buffer, if any
		if err != nil {
			log.Infof("zapLogger sync err: %s", err.Error())
		}
		os.Exit(0)
	}()
}

type Options struct {
	appName    string
	env        EnvLogger
	level      zapcore.Level
	callerSkip int
	fileLog    *FileLogger
}

type Option func(*Options)

var (
	DefaultEnv        = LogDevelopEnv
	DefaultLevel      = zapcore.DebugLevel
	DefaultCallerSkip = 1
)

func Env(env EnvLogger) Option {
	return func(o *Options) {
		o.env = env
	}
}

func Level(level zapcore.Level) Option {
	return func(o *Options) {
		o.level = level
	}
}

func CallerSkip(callerSkip int) Option {
	return func(o *Options) {
		o.callerSkip = callerSkip
	}
}

func FileLog(fileLog *FileLogger) Option {
	return func(o *Options) {
		o.fileLog = fileLog
	}
}

// SetLogger: customize yourself logger.
func SetLogger(options ...Option) (err error) {
	opts := Options{
		env:        DefaultEnv,
		level:      DefaultLevel,
		callerSkip: DefaultCallerSkip,
	}

	for _, o := range options {
		o(&opts)
	}

	if opts.fileLog != nil {
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zapLoggerEncoderConfig),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   opts.fileLog.Filename,
				MaxSize:    opts.fileLog.MaxSize,
				MaxAge:     opts.fileLog.MaxAge,
				MaxBackups: opts.fileLog.MaxBackups,
				LocalTime:  opts.fileLog.LocalTime,
				Compress:   opts.fileLog.Compress,
			}),
			zap.NewAtomicLevelAt(opts.level),
		)
		zapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(opts.callerSkip))
	} else {
		if opts.env == LogReleaseEnv {
			zapLoggerConfig = zap.NewProductionConfig()
			zapLoggerConfig.EncoderConfig = zapLoggerEncoderConfig
		}

		zapLoggerConfig.Level = zap.NewAtomicLevelAt(opts.level)
		zapLogger, err = zapLoggerConfig.Build(zap.AddCaller(), zap.AddCallerSkip(opts.callerSkip))
		if err != nil {
			return err
		}
	}

	log = zapLogger.Sugar()
	return nil
}

// GetLogger get logger
func GetLogger() Logger {
	return log
}

// SetLoggerCallerDisable: disable caller info in production env for performance improve.
// It is highly recommended that you execute this method in a production environment.
func SetLoggerCallerDisable() (err error) {
	zapLoggerConfig.Development = false
	zapLoggerConfig.DisableCaller = true
	zapLogger, err = zapLoggerConfig.Build()
	if err != nil {
		return err
	}
	log = zapLogger.Sugar()
	return nil
}
