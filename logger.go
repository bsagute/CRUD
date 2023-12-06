package logger

import (
	"digi-model-engine/utils/constants"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Logger = logrus.New()

// this is the entry point for the logger as it initializes the logger
func InitLogger() {

	// Set logger format
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     true,
		ForceQuote:      true,
	})

	// Create log file
	currDate := time.Now().Format(constants.DateMonthYearFormat)
	filePath := constants.LogPath + currDate + ".log"
	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	// Set logger output to file
	if err == nil {
		Logger.SetOutput(logFile)
	} else {
		Logger.Info("Failed to set file as logging output" + err.Error())
	}

	setLoggerLevel(viper.GetString("LOG_LEVEL"))
	ConfigureSentry()
}

func setLoggerLevel(l string) {
	l = strings.ToUpper(l)
	switch l {
	case "PANIC":
		Logger.SetLevel(logrus.PanicLevel)
	case "FATAL":
		Logger.SetLevel(logrus.FatalLevel)
	case "ERROR":
		Logger.SetLevel(logrus.ErrorLevel)
	case "WARN":
		Logger.SetLevel(logrus.WarnLevel)
	case "INFO":
		Logger.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		Logger.SetLevel(logrus.DebugLevel)
	case "TRACE":
		Logger.SetLevel(logrus.TraceLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}
}

func Info(msg string, fields map[string]interface{}) {
	if fields == nil {
		Logger.Info(msg)
	} else {
		Logger.Info(msg, fields)
	}
}

func Debug(msg string, fields map[string]interface{}) {
	if fields == nil {
		Logger.Debug(msg)
	} else {
		Logger.Debug(msg, fields)
	}
}

func Error(err error, msg interface{}) {
	if msg == "" {
		msg = err.Error()
	}
	NotifySentry(err)
	Logger.WithError(errors.WithStack(err)).Error(msg)
}

func Fatal(err error, msg string, fields map[string]interface{}) {
	if msg == "" {
		msg = err.Error()
	}
	NotifySentry(err)
	Logger.WithError(errors.WithStack(err)).Fatal(msg)
}

func Panic(err error, msg string, fields map[string]interface{}) {
	if msg == "" {
		msg = err.Error()
	}
	NotifySentry(err)
	Logger.WithError(errors.WithStack(err)).Panic(msg)
}
