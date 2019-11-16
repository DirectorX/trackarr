package loghook

import "github.com/sirupsen/logrus"

/* Public - Logrus */
func (l *Loghooker) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

func (l *Loghooker) Fire(entry *logrus.Entry) error {
	// dont process entry when loghook is not enabled
	if !l.running.Load() {
		return nil
	}

	// push entry to be processed
	_ = l.Push(entry)
	return nil
}
