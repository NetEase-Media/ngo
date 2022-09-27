package optimus

import "github.com/NetEase-Media/ngo/pkg/log"

// region ngoLoggerWrapper
type ngoLoggerWrapper struct {
	log.Logger
}

func newNgoLoggerWrapper(ngoLogger log.Logger) *ngoLoggerWrapper {
	return &ngoLoggerWrapper{Logger: ngoLogger}
}

func (wrapper *ngoLoggerWrapper) Error(msg string) {
	wrapper.Logger.Error(msg)
}
