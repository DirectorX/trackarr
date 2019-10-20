package cache

import (
	"github.com/ReneKroon/ttlcache"
	"github.com/l3uddz/trackarr/logger"
	"time"
)

var (
	cache *ttlcache.Cache
	log   = logger.GetLogger("cache")
)

/* Public */
func Close() {
	cache.Close()
}

func AddItem(key string, value *CacheItem) {
	cache.Set(key, *value)
}

func GetItem(key string) (*CacheItem, bool) {
	result, ok := cache.Get(key)
	if !ok {
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

	cache.SetTTL(60 * time.Second)
	cache.SetExpirationCallback(cacheItemExpired)

	return nil
}

/* Private - Callbacks */
func cacheItemExpired(key string, value interface{}) {
	log.Tracef("Cleared item: %s", key)
}
