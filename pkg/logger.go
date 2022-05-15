package pkg

import (
	"fmt"
	"log"
)

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
}

type defaultLogger struct {
	Level Level
}

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (d *defaultLogger) log(level Level, msg string) {
	if level >= d.Level {
		switch level {
		case DEBUG:
			log.Printf("DEBUG - %s", msg)
		case INFO:
			log.Printf("INFO - %s", msg)
		case WARN:
			log.Printf("WARN - %s", msg)
		case ERROR:
			log.Printf("ERROR - %s", msg)
		}
	}
}

func (d *defaultLogger) logf(level Level, msg string, v ...any) {
	d.log(level, fmt.Sprintf(msg, v...))
}

func (d *defaultLogger) Debug(s string) {
	d.log(DEBUG, s)
}

func (d *defaultLogger) Info(s string) {
	d.log(INFO, s)
}

func (d *defaultLogger) Warn(s string) {
	d.log(WARN, s)
}

func (d *defaultLogger) Error(s string) {
	d.log(ERROR, s)
}

func (d *defaultLogger) Debugf(s string, i ...interface{}) {
	d.logf(DEBUG, s, i...)
}

func (d *defaultLogger) Infof(s string, i ...interface{}) {
	d.logf(INFO, s, i...)

}

func (d defaultLogger) Warnf(s string, i ...interface{}) {
	d.logf(WARN, s, i...)

}

func (d *defaultLogger) Errorf(s string, i ...interface{}) {
	d.logf(ERROR, s, i...)

}

func DefaultLogger(level Level) Logger {
	return &defaultLogger{Level: level}
}
