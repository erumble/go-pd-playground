package logger

// LeveledLogger has log levels Debug, Info, Error
type LeveledLogger interface {
	debugLogger
	errorLogger
	infoLogger

	syncer

	Indent(name string) LeveledLogger
}

type leveledLogger struct {
	zapLogger
}

// NewLeveledLogger instantiates a zap.SugaredLogger to satisfy the LeveledLogger interface.
func NewLeveledLogger(logLevel string) LeveledLogger {
	sugaredLogger := newZapSugaredLogger(logLevel)
	return &leveledLogger{sugaredLogger}
}

// Indent is an alias for Named, see https://godoc.org/go.uber.org/zap#Logger.Named
func (l leveledLogger) Indent(s string) LeveledLogger {
	newLogger := l.Named(s)
	return &leveledLogger{
		newLogger,
	}
}
