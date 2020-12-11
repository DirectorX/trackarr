package cache

import (
	"fmt"
	"time"

	"gitlab.com/cloudb0x/trackarr/logger"

	"github.com/ReneKroon/ttlcache/v2"
)

var (
	cache *ttlcache.Cache
	log   = logger.GetLogger("cache")
)

/* Public */
func Close() {
	_ = cache.Close()
}

func AddItem(key string, value *CacheItem) {
	_ = cache.Set(key, *value)
}

func GetItem(key string) (*CacheItem, bool) {
	result, err := cache.Get(key)
	if err != nil {
		return nil, false
	}

	// assert result type
	item, ok := result.(CacheItem)
	if !ok {
		_ = cache.Remove(key)
		return nil, false
	}

	return &item, true
}

/* Private */
func Init() error {
	cache = ttlcache.NewCache()

	if err := cache.SetTTL(60 * time.Second); err != nil {
		return fmt.Errorf("set ttl: %w", err)
	}

	cache.SetExpirationCallback(cacheItemExpired)
	return nil
}

/* Private - Callbacks */
func cacheItemExpired(key string, value interface{}) {
	log.Tracef("Cleared item: %s", key)
}
