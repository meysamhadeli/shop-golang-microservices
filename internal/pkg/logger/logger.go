package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	LogLevel string `mapstructure:"level"`
}

// NewAppLogger App Logger constructor
func NewAppLogger(cfg *Config) *appLogger {
	return &appLogger{level: cfg.LogLevel}
}

// Logger methods interface
type ILogger interface {
	InitLogger()
	getLevel() log.Level
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
}

// Application logger
type appLogger struct {
	level  string
	logger *log.Logger
}

// For mapping config logger to email_service logger levels
var loggerLevelMap = map[string]log.Level{
	"debug": log.DebugLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"panic": log.PanicLevel,
	"fatal": log.FatalLevel,
	"trace": log.TraceLevel,
}

func (l *appLogger) getLevel() log.Level {

	level, exist := loggerLevelMap[l.level]
	if !exist {
		return log.DebugLevel
	}

	return level
}

// InitLogger Init logger
func (l *appLogger) InitLogger() {

	l.logger = log.StandardLogger()

	logLevel := l.getLevel()

	env := os.Getenv("APP_ENV")

	if env == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		log.SetFormatter(&log.TextFormatter{
			DisableColors: false,
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	log.SetLevel(logLevel)
}

func (l *appLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *appLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *appLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *appLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *appLogger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *appLogger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

func (l *appLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *appLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *appLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *appLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *appLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *appLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *appLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *appLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
