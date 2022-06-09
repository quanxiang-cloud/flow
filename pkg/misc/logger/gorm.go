package logger

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	lg "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Gorm gorm logger
type Gorm struct {
	writer *zap.SugaredLogger
	lg.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// NewGorm new gorm
func NewGorm(writer *zap.SugaredLogger, config lg.Config) lg.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	return &Gorm{
		writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (l *Gorm) LogMode(level lg.LogLevel) lg.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func msssage(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// Info print info
func (l Gorm) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lg.Info {
		l.writer.Infof(msssage(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

// Warn print warn messages
func (l Gorm) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lg.Warn {
		l.writer.Warnf(msssage(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

// Error print error messages
func (l Gorm) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lg.Error {
		l.writer.Errorf(msssage(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

// Trace print sql message
func (l Gorm) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= lg.Error:
			sql, rows := fc()
			if rows == -1 {
				l.writer.Debugf(msssage(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql))
			} else {
				l.writer.Debugf(msssage(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql))
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= lg.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.writer.Debugf(msssage(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql))
			} else {
				l.writer.Debugf(msssage(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql))
			}
		case l.LogLevel >= lg.Info:
			sql, rows := fc()
			if rows == -1 {
				l.writer.Debugf(msssage(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql))
			} else {
				l.writer.Debugf(msssage(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql))
			}
		}
	}
}
