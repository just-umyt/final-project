package zap

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	myLog "stocks/internal/observability/log"
)

const (
	filePerm             = 0644
	errorFailedToSyncLog = "failed to sync logger: %v"
	errorCloseLogFile    = "failed to close log file: %v"
)

var _ myLog.Logger = &Logger{}

type Logger struct {
	L *zap.Logger
}

func NewLogger(path string) (*Logger, func(), error) {
	//production Logger
	prodLogCoreEncoderConfig := zap.NewProductionEncoderConfig()
	prodLogCoreEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, filePerm)
	if err != nil {
		return nil, nil, err
	}

	prodLogCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(prodLogCoreEncoderConfig),
		zapcore.AddSync(file),
		zapcore.InfoLevel,
	)

	//developmet Logger
	devLogCoreEncoderConfig := zap.NewDevelopmentEncoderConfig()
	devLogCoreEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	devLogCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(devLogCoreEncoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	//new zap logger
	zapLogger := zap.New(zapcore.NewTee(prodLogCore, devLogCore))

	newLog := &Logger{
		L: zapLogger,
	}

	cleanup := func() {
		err := zapLogger.Sync()
		if err != nil {
			newLog.Errorf(errorFailedToSyncLog, err)
		}

		err = file.Close()
		if err != nil {
			newLog.Errorf(errorCloseLogFile, err)
		}
	}

	return newLog, cleanup, nil
}

// Trace logs at Trace log level using fields.
func (l *Logger) Trace(msg string, fields ...myLog.Field) {
	if ce := l.L.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Tracef logs at Trace log level using fmt formatter.
func (l *Logger) Tracef(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.DebugLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Debug logs at Debug log level using fields.
func (l *Logger) Debug(msg string, fields ...myLog.Field) {
	if ce := l.L.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Debugf logs at Debug log level using fmt formatter.
func (l *Logger) Debugf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.DebugLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Info logs at Info log level using fields.
func (l *Logger) Info(msg string, fields ...myLog.Field) {
	if ce := l.L.Check(zap.InfoLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Infof logs at Info log level using fmt formatter.
func (l *Logger) Infof(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.InfoLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Warn logs at Warn log level using fields.
func (l *Logger) Warn(msg string, fields ...myLog.Field) {
	if ce := l.L.Check(zap.WarnLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Warnf logs at Warn log level using fmt formatter.
func (l *Logger) Warnf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.WarnLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Error logs at Error log level using fields.
func (l *Logger) Error(msg string, fields ...myLog.Field) {
	if ce := l.L.Check(zap.ErrorLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Errorf logs at Error log level using fmt formatter.
func (l *Logger) Errorf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.ErrorLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Fatal logs at Fatal log level using fields.
func (l *Logger) Fatal(msg string, fields ...myLog.Field) {
	if ce := l.L.Check(zap.FatalLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Fatalf logs at Fatal log level using fmt formatter.
func (l *Logger) Fatalf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.FatalLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}
