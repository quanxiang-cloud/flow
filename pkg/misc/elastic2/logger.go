package elastic2

import (
	"fmt"

	"go.uber.org/zap"
)

type logger struct {
	l *zap.SugaredLogger
}

func (l logger) Printf(format string, v ...interface{}) {
	l.l.Info(zap.String("elastic", fmt.Sprintf(format, v...)))
}
