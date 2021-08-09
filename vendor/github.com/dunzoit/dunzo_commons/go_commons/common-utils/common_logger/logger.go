package common_logger

import (
	"net/http"

	"go.uber.org/zap"
)

var log *zapLogger
var logConfig *LogConfig

type Level int8

// -1 ensures default value will always be zero and hence info level
const (
	DEBUG Level = iota - 1
	INFO
	WARN
	ERROR
	FATAL
)

func init() {
	cnfg := configureZap(Level(1))
	zl, _ := cnfg.Build(zap.AddCallerSkip(1))
	log = &zapLogger{zl}
}

// The Logger interface denotes the necessary functions to be provided
// by a logging struct. For the first version, we are using a zap logger
type Logger interface {
	Info(msg string, fields ...ContextKV)
	Warn(msg string, fields ...ContextKV)
	Error(msg string, fields ...ContextKV)
	Debug(msg string, fields ...ContextKV)
	Fatal(msg string, fields ...ContextKV)
	With(fields ...ContextKV) Logger
	WithOptions(options ...Option) Logger
}

func Int(key string, val int) ContextKV {
	return ContextKV{zap.Int(key, val)}
}

func String(key string, val string) ContextKV {
	return ContextKV{zap.String(key, val)}
}

func Float(key string, val float64) ContextKV {
	return ContextKV{zap.Float64(key, val)}
}

func Bool(key string, val bool) ContextKV {
	return ContextKV{zap.Bool(key, val)}
}

func Error(err error) ContextKV {
	return ContextKV{zap.Error(err)}
}

// CreateLoggerForRequest returns a logger with context extracted from the request
// as specified by the config
func CreateLoggerForRequest(request *http.Request) Logger {
	contextFields := ExtractContextFieldFromRequest(request)
	var contextKVs []ContextKV
	for field, value := range contextFields {
		contextKVs = append(contextKVs, String(field, value))
	}
	return log.With(contextKVs...)
}

// GlobalLogger returns a raw logger to be used when no context is present
func GlobalLogger() Logger {
	return log
}

// DONOT USE EXCEPT FOR LEGACY SUPPORT
func BaseSugaredLogger() *zap.SugaredLogger {
	return log.Sugar()
}
