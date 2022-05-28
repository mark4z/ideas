package log

import (
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
)

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	log = sugar
}

func Infof(fmt string, args ...interface{}) {
	log.Infof(fmt, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(fmt string, args ...interface{}) {
	log.Errorf(fmt, args...)
}

func Warnf(fmt string, args ...interface{}) {
	log.Warnf(fmt, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}
