package cache

import (
	"gitlab.com/cloudb0x/trackarr/config"
)

/* Structs */

type CacheItem struct {
	Name    string
	Data    []byte
	Release *config.ReleaseInfo
}
