package clientCache

import (
	"context"
	"errors"
	"fmt"
	"github.com/NRKA/Cache/internal/datasource"
	"github.com/NRKA/Cache/pkg/dbErrors"
	"time"
)

const cacheExpiration = 2 * time.Second

type Client struct {
	db    datasource.Datasource
	cache datasource.Datasource
}

func NewClientCache(db datasource.Datasource, cache datasource.Datasource) *Client {
	return &Client{db: db, cache: cache}
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := c.cache.Set(ctx, key, value, expiration)
	if err != nil {
		return fmt.Errorf("failed to set key to cache")
	}

	//updating database
	err = c.db.Set(ctx, key, value, 0)
	if err != nil {
		return fmt.Errorf("failed to update database: %v", err)
	}
	return nil
}

func (c *Client) Get(ctx context.Context, key string) (any, error) {
	value, err := c.cache.Get(ctx, key)
	if errors.Is(err, dbErrors.ErrKeyNotFound) {
		valueDb, err := c.db.Get(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("failed to get key from database: %v", err)
		}
		err = c.cache.Set(ctx, key, valueDb, cacheExpiration)
		if err != nil {
			return nil, fmt.Errorf("failed to set key to cache from db: %v", err)
		}
		return valueDb, nil
	}
	return value, nil
}
