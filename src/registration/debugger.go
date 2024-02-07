package registration

import (
	"github.com/danbordeanu/go-logger"
)

type Debugger struct{}

func (d *Debugger) Debug(format string, v ...interface{}) {
	log := logger.SugaredLogger().With("package", "registration", "source", "debugger")
	log.Debugf(format, v...)
}
