package web

import (
	"golang.org/x/time/rate"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	rateLimiters map[string]*rate.Limiter
	mtx          sync.Mutex
)

func GetRateLimiter(name string, limit int, duration time.Duration) *rate.Limiter {
	// acquire lock
	mtx.Lock()
	defer mtx.Unlock()

	// init map
	if rateLimiters == nil {
		rateLimiters = make(map[string]*rate.Limiter)
		log.Trace("Initialized rateLimiters map")
	}

	// retrieve or create new ratelimit
	var rl *rate.Limiter
	ok := false
	lowerName := strings.ToLower(name)

	rl, ok = rateLimiters[lowerName]
	if !ok {
		l := rate.Every(duration / time.Duration(limit))
		rl = rate.NewLimiter(l, 1)

		rateLimiters[lowerName] = rl

		log.WithFields(logrus.Fields{
			"name":  name,
			"limit": limit,
		}).Trace("Created new ratelimit")
	}

	return rl
}
