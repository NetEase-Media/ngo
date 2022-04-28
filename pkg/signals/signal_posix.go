//go:build !windows
// +build !windows

package signals

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM}
