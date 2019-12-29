package loghook

import (
	"io/ioutil"
	"os"
	"testing"

	"gitlab.com/cloudb0x/trackarr/ws"
)

var (
	l *Loghooker
)

func TestMain(m *testing.M) {
	if err := ws.Init(); err != nil {
		log.WithError(err).Fatal()
	}

	l = NewLoghooker()

	if err := l.Start(); err != nil {
		log.WithError(err).Fatal()
	}

	log.Logger.SetOutput(ioutil.Discard)
	m.Run()
	log.Logger.SetOutput(os.Stdout)

	if err := l.Stop(); err != nil {
		log.WithError(err).Fatal()
	}

	os.Exit(0)
}

func benchmarkLoghook(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Info("benchmark")
	}
}

func BenchmarkLoghook(b *testing.B) { benchmarkLoghook(b) }
