package cache

import (
	"context"
	"github.com/NRKA/Cache/pkg/dbErrors"
	"github.com/NRKA/Cache/pkg/structs"
	"sync"
	"time"
)

type Cache struct {
	mu     sync.RWMutex
	source map[string]structs.Value
}

func NewCache() *Cache {
	return &Cache{source: make(map[string]structs.Value)}
}
func (c *Cache) Get(ctx context.Context, key string) (any, error) {
	c.mu.RLock()
	val, exists := c.source[key]
	c.mu.RUnlock()
	if exists {
		if time.Now().After(val.Expiration) {
			c.mu.Lock()
			delete(c.source, key)
			c.mu.Unlock()
			return nil, dbErrors.ErrKeyNotFound
		}
		return val.Val, nil
	}
	return nil, dbErrors.ErrKeyNotFound
}

func (c *Cache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.source[key] = structs.NewValue(value, expiration)
	return nil
}
