package log

import (
	"fmt"
	"log"
	"os"
)

// LogLevel log level type.
type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
	Fatal
	Panic
)

var DEFAULT_LOG_LEVEL LogLevel

func (l LogLevel) String() string {
	return [...]string{"Debug", "Info", "Warn", "Error", "Fatal", "Panic"}[l]
}

func GetLevel(name string) LogLevel {
	return map[string]LogLevel{"Debug": Debug, "Info": Info, "Warn": Warn, "Error": Error, "Fatal": Fatal, "Panic": Panic}[name]
}

func init() {
	var env string = os.Getenv("LOGLEVEL")
	if env == "" {
		env = "Error"
	}

	DEFAULT_LOG_LEVEL = GetLevel(env)
}

// Logger logger interface, loggers check the env for LOGLEVEL flags. If LOGLEVEL env variable is set then that level and above will be logged. Fatal and Panic always fire.
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

// Logger simple log wrapper.
type Standard struct {
	debug    *log.Logger
	info     *log.Logger
	warn     *log.Logger
	error    *log.Logger
	fatal    *log.Logger
	panic    *log.Logger
	LogLevel LogLevel
}

// New create a new logger.
func New(tag string, logLevel LogLevel) *Standard {
	return &Standard{
		LogLevel: logLevel,
		debug:    log.New(os.Stderr, "DEBUG "+tag+": ", log.Ldate|log.Ltime|log.Lshortfile),
		info:     log.New(os.Stderr, "INFO "+tag+": ", log.Ldate|log.Ltime|log.Lshortfile),
		warn:     log.New(os.Stderr, "WARNING "+tag+": ", log.Ldate|log.Ltime|log.Lshortfile),
		error:    log.New(os.Stderr, "ERROR "+tag+": ", log.Ldate|log.Ltime|log.Lshortfile),
		fatal:    log.New(os.Stderr, "Fatal "+tag+": ", log.Ldate|log.Ltime|log.Lshortfile),
		panic:    log.New(os.Stderr, "Panic "+tag+": ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (s *Standard) Debug(v ...interface{}) {
	if s.LogLevel <= Debug {
		s.debug.Output(2, fmt.Sprint(v...))
	}
}
func (s *Standard) Debugf(format string, v ...interface{}) {
	if s.LogLevel <= Debug {
		s.debug.Output(2, fmt.Sprintf(format, v...))
	}
}

func (s *Standard) Info(v ...interface{}) {
	if s.LogLevel <= Info {
		s.info.Output(2, fmt.Sprint(v...))
	}
}
func (s *Standard) Infof(format string, v ...interface{}) {
	if s.LogLevel <= Info {
		s.info.Output(2, fmt.Sprintf(format, v...))
	}
}

func (s *Standard) Warn(v ...interface{}) {
	if s.LogLevel <= Warn {
		s.warn.Output(2, fmt.Sprint(v...))
	}
}
func (s *Standard) Warnf(format string, v ...interface{}) {
	if s.LogLevel <= Warn {
		s.warn.Output(2, fmt.Sprintf(format, v...))
	}
}

func (s *Standard) Error(v ...interface{}) {
	if s.LogLevel <= Error {
		s.error.Output(2, fmt.Sprint(v...))
	}
}
func (s *Standard) Errorf(format string, v ...interface{}) {
	if s.LogLevel <= Error {
		s.error.Output(2, fmt.Sprintf(format, v...))
	}
}

func (s *Standard) Fatal(v ...interface{}) {
	s.fatal.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}
func (s *Standard) Fatalf(format string, v ...interface{}) {
	s.fatal.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (s *Standard) Panic(v ...interface{}) {
	s.panic.Output(2, fmt.Sprint(v...))
	panic(fmt.Sprint(v...))
}
func (s *Standard) Panicf(format string, v ...interface{}) {
	s.panic.Output(2, fmt.Sprintf(format, v...))
	panic(fmt.Sprintf(format, v...))
}
