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
	if !l.running {
		return nil
	}

	// Increment WaitGroup
	l.wg.Add(1)
	// push entry to be processed
	_ = l.Push(entry)
	// WaitGroup done
	l.wg.Done()

	return nil
}
