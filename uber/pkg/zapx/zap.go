package zapx

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func NewLogger() (*Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"

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
