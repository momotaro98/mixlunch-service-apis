package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type ErrorLevel string

const (
	Debug = ErrorLevel("Debug")
	Info  = ErrorLevel("Info")
	Warn  = ErrorLevel("Warn")
	Error = ErrorLevel("Error")
)

type Logger interface {
	Log(errorLevel ErrorLevel, requestId, message string)
}

func ProvideLogger(conf *Config) Logger {
	l := log.StandardLogger()
	// logrus default New function of StandardLogger
	// &Logger{
	//	Out:          os.Stderr,
	//	Formatter:    new(TextFormatter),
	//	Hooks:        make(LevelHooks),
	//	Level:        InfoLevel,
	//	ExitFunc:     os.Exit,
	//	ReportCaller: false,
	//}

	// Log as JSON instead of the default ASCII formatter.
	l.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	l.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	switch conf.ErrorLevel {
	case Debug:
		l.SetLevel(log.DebugLevel)
	case Info:
		l.SetLevel(log.InfoLevel)
	case Warn:
		l.SetLevel(log.WarnLevel)
	case Error:
		l.SetLevel(log.ErrorLevel)
	default:
		l.SetLevel(log.InfoLevel)
	}

	return &logger{
		logger: l,
	}
}

type logger struct {
	logger *log.Logger
}

func (l *logger) Log(errorLevel ErrorLevel, requestId, message string) {
	switch errorLevel {
	case Debug:
		l.debug(requestId, message)
	case Info:
		l.info(requestId, message)
	case Warn:
		l.warn(requestId, message)
	case Error:
		l.error(requestId, message)
	}
}

var (
	reqIdKeyText      = "request-id"
	stackTraceKeyText = "stack_trace"
)

func (l *logger) debug(requestId, message string) {
	//// stack trace
	//buf := make([]byte, 1024)
	//stackTrace := runtime.Stack(buf, true)
	l.logger.WithFields(log.Fields{
		reqIdKeyText: requestId,
		//stackTraceKeyText: buf[:stackTrace],
	}).Debug(message)
}

func (l *logger) info(requestId, message string) {
	l.logger.WithFields(log.Fields{
		reqIdKeyText: requestId,
	}).Info(message)
}

func (l *logger) warn(requestId, message string) {
	l.logger.WithFields(log.Fields{
		reqIdKeyText: requestId,
	}).Warn(message)
}

func (l *logger) error(requestId, message string) {
	// TODO: Set Stack Trace by using runtime or error
	// We would address it by using new error handling in Go1.13 era.
	//// stack trace
	//buf := make([]byte, 1024)
	//stackTrace := runtime.Stack(buf, true)
	l.logger.WithFields(log.Fields{
		reqIdKeyText: requestId,
		//stackTraceKeyText: buf[:stackTrace],
	}).Error(message)
}
