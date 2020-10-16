package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type debugLogger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}

type errorLogger interface {
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

type infoLogger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
}

type fatalLogger interface {
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
}

type syncer interface {
	Sync() error
}

type zapLogger interface {
	debugLogger
	infoLogger
	errorLogger
	fatalLogger

	syncer

	Named(name string) *zap.SugaredLogger
}

func newZapSugaredLogger(logLevel string) *zap.SugaredLogger {
	// In the case of the logger, we always want to return a working one. So we will
	// default the log level to INFO, and if an invalid log level was passed in, we
	// will log it as an error using the INFO level logger.
	atom := zap.NewAtomicLevelAt(zap.InfoLevel)

	logLevel = strings.ToUpper(logLevel)
	var err error

	switch logLevel {
	case "DEBUG", "INFO", "ERROR", "FATAL":
		_ = (&atom).UnmarshalText([]byte(logLevel))
	default:
		// We don't have the logger to log this error built yet, save it for later.
		err = fmt.Errorf("log level %v is not a valid log level, defaulting to INFO", logLevel)
	}

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && lvl >= atom.Level()
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl >= atom.Level()
	})

	encoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())

	stdoutWriter := zapcore.Lock(os.Stdout)
	stderrWriter := zapcore.Lock(os.Stderr)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, stdoutWriter, lowPriority),
		zapcore.NewCore(encoder, stderrWriter, highPriority),
	)

	logger := zap.New(core)
	logger = logger.WithOptions(
		zap.AddCaller(),
	)

	sugaredLogger := logger.Sugar()

	// the error from the switch, we can log it now...
	if err != nil {
		sugaredLogger.Error(err)
	}

	// Always good to let everyone know what the log level is set to.
	// sugaredLogger.Infof("Log level set to: %s", atom.Level())

	return sugaredLogger
}
