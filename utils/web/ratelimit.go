package web

import (
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"
)

var (
	rateLimiters map[string]ratelimit.Limiter
	mtx          sync.Mutex
)

func GetRateLimiter(name string, newRateLimit int, limitSeconds int) *ratelimit.Limiter {
	// acquire lock
	mtx.Lock()
	defer mtx.Unlock()

	// init map
	if rateLimiters == nil {
		rateLimiters = make(map[string]ratelimit.Limiter)
		log.Trace("Initialized rateLimiters map")
	}

	// retrieve or create new ratelimit
	var rl ratelimit.Limiter
	ok := false
	lowerName := strings.ToLower(name)

	rl, ok = rateLimiters[lowerName]
	if !ok {
		rl = ratelimit.New(newRateLimit, ratelimit.WithoutSlack, ratelimit.Per(time.Duration(limitSeconds)*time.Second))
		rateLimiters[lowerName] = rl

		log.WithFields(logrus.Fields{
			"name":  name,
			"limit": newRateLimit,
		}).Trace("Created new ratelimit")
	}

	return &rl
}
