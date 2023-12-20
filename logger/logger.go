package logger

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func NewLogger(opts ...Options) *Logger {
	var opt = DefaultOptions()
	for _, o := range opts {
		o(opt)
	}
	logger := logrus.New()
	logger.SetLevel(opt.level)

	if opt.formatter != nil {
		logger.SetFormatter(opt.formatter)
	} else {
		switch opt.format {
		case "json":
			logger.SetFormatter(&logrus.JSONFormatter{})
		default:
			logger.SetFormatter(&logrus.TextFormatter{})
		}
	}
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) exeract(ctx context.Context, fields ...*Field) *logrus.Entry {
	var logFields = make(logrus.Fields)
	if len(fields) > 0 {
		for _, item := range fields {
			logFields[item.key] = item.value
		}
	}
	logger := l.logger.WithContext(ctx)
	if logFields != nil {
		logger = logger.WithFields(logFields)
	}
	return logger
}

func (l *Logger) Info(ctx context.Context, message string, fields ...*Field) {
	l.exeract(ctx, fields...).Info(message)
}

func (l *Logger) Error(ctx context.Context, message string, fields ...*Field) {
	l.exeract(ctx, fields...).Error(message)
}

func (l *Logger) Debug(ctx context.Context, message string, fields ...*Field) {
	l.exeract(ctx, fields...).Debug(message)
}

func (l *Logger) Warning(ctx context.Context, message string, fields ...*Field) {
	l.exeract(ctx, fields...).Warning(message)
}

func (l *Logger) NewModule(name string) *Entry {
	entry := l.logger.WithField("module", name)
	return &Entry{
		logger: entry,
	}
}

type Entry struct {
	logger *logrus.Entry
}

func (l *Entry) exeract(ctx context.Context, fields ...*Field) *logrus.Entry {
	var logFields = make(logrus.Fields)
	if len(fields) > 0 {
		for _, item := range fields {
			logFields[item.key] = item.value
		}
	}
	logger := l.logger.WithContext(ctx)
	if logFields != nil {
		logger = logger.WithFields(logFields)
	}
	return logger
}

func (l *Entry) Info(ctx context.Context, message string, fields ...*Field) {
	l.exeract(ctx, fields...).Info(message)
}

func (l *Entry) Error(ctx context.Context, message string, fields ...*Field) {
	l.exeract(ctx, fields...).Error(message)
}

func (l *Entry) Debug(ctx context.Context, message string, fields ...*Field) {
	l.exeract(ctx, fields...).Debug(message)
}

func (l *Entry) Warning(ctx context.Context, message string, fields ...*Field) {
	l.exeract(ctx, fields...).Warning(message)
}

type Field struct {
	key   string
	value interface{}
}

func Any(k string, v interface{}) *Field {
	return &Field{
		key:   k,
		value: v,
	}
}

func Err(err error) *Field {
	return &Field{
		key:   "error",
		value: err,
	}
}

func Str(key string, format string, args ...interface{}) *Field {
	return &Field{
		key:   key,
		value: fmt.Sprintf(format, args...),
	}
}
