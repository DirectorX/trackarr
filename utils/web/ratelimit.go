package web

import (
	"golang.org/x/time/rate"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	rateLimiters map[string]*RateLimiter
	mtx          sync.Mutex
)

type RateLimitCallback func(rl *rate.Limiter) error

type RateLimiter struct {
	limiter *rate.Limiter
	Take    func(rl *rate.Limiter) error
}

func GetRateLimiter(name string, limit int, duration time.Duration, burst int, check RateLimitCallback) *RateLimiter {
	// acquire lock
	mtx.Lock()
	defer mtx.Unlock()

	// init map
	if rateLimiters == nil {
		rateLimiters = make(map[string]*RateLimiter)
		log.Trace("Initialized rateLimiters map")
	}

	// retrieve or create new ratelimit
	var rl *RateLimiter
	ok := false
	lowerName := strings.ToLower(name)

	rl, ok = rateLimiters[lowerName]
	if !ok {
		rateLimiters[lowerName] = newRateLimiter(limit, duration, burst, check)

		log.WithFields(logrus.Fields{
			"name":  name,
			"limit": limit,
		}).Trace("Created new ratelimit")
	}

	return rl
}

func newRateLimiter(limit int, duration time.Duration, burst int, check RateLimitCallback) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Every(duration/time.Duration(limit)), burst),
		Take:    check,
	}
}
