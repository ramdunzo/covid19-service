package common_logger

import (
	"go.uber.org/zap"
)

// SetUpLogging initiates the logger using the config provided
// and should be put into the app startup code
func SetUpLogging(lc *LogConfig) {
	logConfig = lc
	config := configureZap(Level(lc.Level))
	l, _ := config.Build(zap.AddCallerSkip(1))
	log = &zapLogger{l}
}

// config maps the level from the log config and builds the base logger using the config provided
func configureZap(level Level) zap.Config {
	defaultLevel := zap.InfoLevel
	switch level {
	case DEBUG:
		defaultLevel = zap.DebugLevel
		break
	case INFO:
		defaultLevel = zap.InfoLevel
	case ERROR:
		defaultLevel = zap.ErrorLevel
		break
	}
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(defaultLevel)
	config.Sampling = nil
	config.DisableStacktrace = true
	return config
}

type zapLogger struct {
	*zap.Logger
}

func (b *zapLogger) With(args ...ContextKV) Logger {
	var fields []zap.Field
	for _, el := range args {
		fields = append(fields, el.Field)
	}
	return &zapLogger{
		b.Logger.With(fields...),
	}
}

type ContextKV struct {
	zap.Field
}

type Option struct {
	zap.Option
}

func AddCallerSkip(i int) Option {
	return Option{
		zap.AddCallerSkip(i),
	}
}

func appendCallerFunctionFields(fields []zap.Field) []zap.Field {
	fields = append(fields, zap.String("functionName", Trace(3)))
	return fields
}

func (b *zapLogger) Info(msg string, args ...ContextKV) {
	var fields []zap.Field
	for _, el := range args {
		fields = append(fields, el.Field)
	}
	fields = appendCallerFunctionFields(fields)
	b.Logger.Info(msg, fields...)
}

func (b *zapLogger) Debug(msg string, args ...ContextKV) {
	var fields []zap.Field
	for _, el := range args {
		fields = append(fields, el.Field)
	}
	fields = appendCallerFunctionFields(fields)
	b.Logger.Debug(msg, fields...)
}

func (b *zapLogger) Warn(msg string, args ...ContextKV) {
	var fields []zap.Field
	for _, el := range args {
		fields = append(fields, el.Field)
	}
	fields = appendCallerFunctionFields(fields)
	b.Logger.Warn(msg, fields...)
}

func (b *zapLogger) Error(msg string, args ...ContextKV) {
	var fields []zap.Field
	for _, el := range args {
		fields = append(fields, el.Field)
	}
	fields = appendCallerFunctionFields(fields)
	b.Logger.Error(msg, fields...)
}

func (b *zapLogger) Fatal(msg string, args ...ContextKV) {
	var fields []zap.Field
	for _, el := range args {
		fields = append(fields, el.Field)
	}
	fields = appendCallerFunctionFields(fields)
	b.Logger.Fatal(msg, fields...)
}

func (b *zapLogger) WithOptions(options ...Option) Logger {
	var opts []zap.Option
	for _, el := range options {
		opts = append(opts, el)
	}

	return &zapLogger{
		b.Logger.WithOptions(opts...),
	}
}
