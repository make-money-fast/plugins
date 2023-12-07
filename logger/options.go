package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

type Option struct {
	output    io.Writer
	closer    io.Closer
	level     logrus.Level
	format    string
	formatter logrus.Formatter
}

func DefaultOptions() *Option {
	return &Option{
		output: os.Stdout,
		level:  logrus.DebugLevel,
		format: "json",
	}
}

type Options func(option *Option)

func Stdout() Options {
	return func(o *Option) {
		o.output = os.Stdout
	}
}

func File(file string) Options {
	return func(option *Option) {
		fi, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
		if err != nil {
			panic(err)
		}
		option.closer = fi
		option.output = fi
	}
}

func RotateFile(file string, maxMB int, maxBackUp int, maxAge int) Options {
	return func(option *Option) {
		output := &lumberjack.Logger{
			Filename:   file,
			MaxSize:    maxMB, // megabytes
			MaxBackups: maxBackUp,
			MaxAge:     maxAge, //days
			Compress:   true,   // disabled by default
		}

		option.output = output
	}
}

func Level(s string) Options {
	return func(option *Option) {
		lvl, err := logrus.ParseLevel(s)
		if err != nil {
			panic(err)
		}
		option.level = lvl
	}
}

// Format text || json , default is json
func Format(s string) Options {
	return func(option *Option) {
		option.format = s
	}
}

func Formatter(formatter logrus.Formatter) Options {
	return func(option *Option) {
		option.formatter = formatter
	}
}
