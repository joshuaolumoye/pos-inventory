package infrastructure

import (
	"time"

	"github.com/allegro/bigcache/v3"
)

func NewBigCache() (*bigcache.BigCache, error) {
	return bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
}
