package logger

import (
	"github.com/dunzoit/dunzo_commons/go_commons/common-utils/common_dtos"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

func SetUpLogging() {

	log.Formatter = &logrus.JSONFormatter{
		//PrettyPrint: true,
	}
	log.SetLevel(logrus.InfoLevel)
	//log.SetReportCaller(true)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func WithError(err error) *logrus.Entry {
	return log.WithError(err)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, value)
}

func WithContext(context *common_dtos.Context) *logrus.Entry {
	return log.WithField("corelation-id", *context.CorrelationId)
}
