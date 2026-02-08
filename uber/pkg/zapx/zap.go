package zapx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(level ...zapcore.Level) (*Logger, error) {
	level = append(level, zapcore.InfoLevel)

	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.Level = zap.NewAtomicLevelAt(level[0])

	zl, err := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{zl}, nil
}

func (zl *Logger) Close() error {
	return zl.Logger.Sync()
}
