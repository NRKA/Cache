package database

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/NRKA/Cache/pkg/dbErrors"
	"os"
	"strings"
	"sync"
	"time"
)

type Database struct {
	mu      sync.RWMutex
	changed map[string]any
	muKey   map[string]*sync.Mutex
	dir     string
}

func NewDb() *Database {
	db := &Database{
		changed: make(map[string]any),
		dir:     "databaseSource",
		muKey:   make(map[string]*sync.Mutex),
	}
	return db
}

func (db *Database) Begin(ctx context.Context, key string) error {
	db.mu.Lock()
	_, ok := db.muKey[key]
	if !ok {
		db.muKey[key] = &sync.Mutex{}
	}
	mu := db.muKey[key]
	db.mu.Unlock()

	mu.Lock()
	value, err := db.Get(ctx, key)
	if err != nil {
		if errors.Is(err, dbErrors.ErrFileNotFound) {
			return nil
		}
		mu.Unlock()
		return fmt.Errorf("failed to get initial value: %v", err)
	}
	db.mu.Lock()
	db.changed[key] = value
	db.mu.Unlock()
	return nil
}

func (db *Database) Rollback(ctx context.Context, key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	defer func() {
		delete(db.changed, key)
		db.muKey[key].Unlock()
	}()

	_, ok := db.changed[key]
	if ok {
		return db.Set(ctx, key, db.changed[key], 0)
	}

	filePath := fmt.Sprintf("%s/%s.txt", db.dir, key)
	_, err := os.Stat(filePath)
	if err != nil {
		return nil
	}

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

func (db *Database) Commit(key string) {
	db.mu.Lock()
	delete(db.changed, key)
	db.muKey[key].Unlock()
	db.mu.Unlock()
}

func (db *Database) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, err := os.Stat(db.dir); os.IsNotExist(err) {
		err := os.MkdirAll(db.dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory")
		}
	}

	file, err := os.Create(fmt.Sprintf("%s/%s.txt", db.dir, key))
	if err != nil {
		return fmt.Errorf("failed to create file")
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%v\n", value))
	if err != nil {
		return fmt.Errorf("failed to write to file")
	}
	return nil
}

func (db *Database) Get(ctx context.Context, key string) (any, error) {
	file, err := os.Open(fmt.Sprintf("%s/%s.txt", db.dir, key))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, dbErrors.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileReader := bufio.NewReader(file)
	val, err := fileReader.ReadString('\n')

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	val = strings.TrimSpace(val)

	return val, nil
}
