package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapSugaredLogger(w zapcore.WriteSyncer, extra ...interface{}) *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "datetime"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z"))
	}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.Lock(w),
		atom,
	), zap.AddCaller()).Sugar()
	atom.SetLevel(zap.InfoLevel)

	if extra != nil {
		logger = logger.With(extra...)
	}

	return logger
}
