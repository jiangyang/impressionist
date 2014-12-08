package main

import (
	"log"
	"os"
)

const name string = "impressionist"
const (
	_trace_ = iota
	_debug_
	_info_
	_error_
)
const loglevel = _info_

type newlogger struct {
	L *log.Logger
}

func (l *newlogger) trace(format string, v ...interface{}) {
	if loglevel <= _trace_ {
		l.L.Printf(format, v...)
	}
}

func (l *newlogger) debug(format string, v ...interface{}) {
	if loglevel <= _debug_ {
		l.L.Printf(format, v...)
	}
}

func (l *newlogger) info(format string, v ...interface{}) {
	if loglevel <= _info_ {
		l.L.Printf(format, v...)
	}
}

func (l *newlogger) error(format string, v ...interface{}) {
	if loglevel <= _error_ {
		l.L.Printf(format, v...)
	}
}

var l *newlogger = &newlogger{log.New(os.Stdout, name+":", 0)}
