package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	// Logger logger
	Logger *zap.SugaredLogger
)

func init() {
	New(&Config{
		Level:       -1,
		Development: true,
		Sampling: Sampling{
			Initial:    10,
			Thereafter: 10,
		},
		OutputPath:      []string{"stderr"},
		ErrorOutputPath: []string{"stderr"},
	})
}

// New 创建日志
func New(conf *Config) error {
	return newLogger(zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(conf.Level)),
		Development: conf.Development,
		Sampling: &zap.SamplingConfig{
			Initial:    conf.Sampling.Initial,
			Thereafter: conf.Sampling.Thereafter,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      conf.OutputPath,
		ErrorOutputPaths: conf.ErrorOutputPath,
	})
}

func newLogger(conf zap.Config) error {
	var err error
	logger, err = conf.Build(
		zap.AddStacktrace(zap.ErrorLevel),
		zap.WithCaller(true),
	)

	if err != nil {
		return err
	}

	Logger = logger.Sugar()
	return nil
}

// Sync sync
func Sync() {
	if logger != nil {
		logger.Sync()
	}
}
