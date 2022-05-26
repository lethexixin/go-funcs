package graceful

import (
	"os"
	"syscall"
)

var (
	// ShutdownSignals receives shutdown signals to process
	ShutdownSignals = []os.Signal{
		os.Interrupt, os.Kill, syscall.SIGKILL,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT, syscall.SIGTERM,
	}

	// DumpHeapShutdownSignals receives shutdown signals to process
	DumpHeapShutdownSignals = []os.Signal{syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT}
)
