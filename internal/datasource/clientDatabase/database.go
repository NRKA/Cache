package clientDatabase

import (
	"context"
	"fmt"
	"github.com/NRKA/Cache/pkg/database"
	"time"
)

type Client struct {
	db *database.Database
}

func NewClient(db *database.Database) *Client {
	return &Client{db: db}
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := c.db.Begin(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	err = c.db.Set(ctx, key, value, expiration)
	if err != nil {
		err = c.db.Rollback(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to rollback: %v", err)
		}
		return fmt.Errorf("failed to set data to db: %v", err)
	}
	c.db.Commit(key)
	return nil
}

func (c *Client) Get(ctx context.Context, key string) (any, error) {
	value, err := c.db.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get value from database: %v", err)
	}
	return value, nil
}
