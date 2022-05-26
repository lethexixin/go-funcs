package logger

// Debug is debug level
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info is info level
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn is warning level
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error is error level
func Error(args ...interface{}) {
	log.Error(args...)
}

// DPanic is d_panic level
func DPanic(args ...interface{}) {
	log.DPanic(args...)
}

// Panic is panic level
func Panic(args ...interface{}) {
	log.Panic(args...)
}

// Fatal is fatal level
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Debugf is format debug level
func Debugf(fmt string, args ...interface{}) {
	log.Debugf(fmt, args...)
}

// Infof is format info level
func Infof(fmt string, args ...interface{}) {
	log.Infof(fmt, args...)
}

// Warnf is format warning level
func Warnf(fmt string, args ...interface{}) {
	log.Warnf(fmt, args...)
}

// Errorf is format error level
func Errorf(fmt string, args ...interface{}) {
	log.Errorf(fmt, args...)
}

// DPanicf is format d_panic level
func DPanicf(fmt string, args ...interface{}) {
	log.DPanicf(fmt, args...)
}

// Panicf is format panic level
func Panicf(fmt string, args ...interface{}) {
	log.Panicf(fmt, args...)
}

// Fatalf is format fatal level
func Fatalf(fmt string, args ...interface{}) {
	log.Fatalf(fmt, args...)
}
