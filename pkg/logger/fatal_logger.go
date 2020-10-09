package logger

// FatalLeveledLogger has log levels Debug, Info, Error, Fatal
type FatalLeveledLogger interface {
	debugLogger
	errorLogger
	infoLogger
	fatalLogger

	syncer

	Indent(name string) FatalLeveledLogger
	LeveledLogger() LeveledLogger
}

type fatalLeveledLogger struct {
	zapLogger
}

// NewFatalLogger instantiates a zap.SugaredLogger to satisfy the FatalLeveledLogger interface.
func NewFatalLogger(logLevel string) FatalLeveledLogger {
	sugaredLogger := newZapSugaredLogger(logLevel)
	return &fatalLeveledLogger{sugaredLogger}
}

// Indent is an alias for Named, see https://godoc.org/go.uber.org/zap#Logger.Named
func (l fatalLeveledLogger) Indent(s string) FatalLeveledLogger {
	newLogger := l.Named(s)
	return &fatalLeveledLogger{
		newLogger,
	}
}

// LeveledLogger drops the fatal methods from the logger's interface.
func (l fatalLeveledLogger) LeveledLogger() LeveledLogger {
	return &leveledLogger{l}
}
