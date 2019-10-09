package logger

import (
	"github.com/l3uddz/trackarr/utils/strings"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	prefixLen = 5
)

/* Public */

func Init(logLevel int, logFilePath string) error {
	var useLevel logrus.Level

	// determine logging level
	switch logLevel {
	case 0:
		useLevel = logrus.InfoLevel
	case 1:
		useLevel = logrus.DebugLevel
	default:
		useLevel = logrus.TraceLevel
	}

	// set rotating file hook
	fileLogFormatter := &prefixed.TextFormatter{}
	fileLogFormatter.FullTimestamp = true
	fileLogFormatter.QuoteEmptyFields = true
	fileLogFormatter.DisableColors = true
	fileLogFormatter.ForceFormatting = true

	rotateFileHook, err := NewRotateFileHook(RotateFileConfig{
		Filename:   logFilePath,
		MaxSize:    5,
		MaxBackups: 10,
		MaxAge:     90,
		Level:      useLevel,
		Formatter:  fileLogFormatter,
	})

	if err != nil {
		logrus.WithError(err).Errorf("Failed initializing rotating file log to %q", logFilePath)
		return err
	}
	logrus.AddHook(rotateFileHook)

	// set console formatter
	logFormatter := &prefixed.TextFormatter{}
	logFormatter.FullTimestamp = true
	logFormatter.QuoteEmptyFields = true
	logFormatter.ForceFormatting = true
	logrus.SetFormatter(logFormatter)

	// set logging level
	logrus.SetLevel(useLevel)

	// info
	log := logrus.WithField("prefix", strings.StringLeftJust("log", " ", prefixLen))
	log.Infof("Using %s = %s", strings.StringLeftJust("LOG_LEVEL", " ", 10), logrus.GetLevel().String())
	log.Infof("Using %s = %q", strings.StringLeftJust("LOG", " ", 10), logFilePath)

	return nil
}

func GetLogger(prefix string) *logrus.Entry {
	if len(prefix) > prefixLen {
		prefixLen = len(prefix)
	}

	return logrus.WithFields(logrus.Fields{"prefix": strings.StringLeftJust(prefix, " ", prefixLen)})
}