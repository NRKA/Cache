package main

import (
	"context"
	"github.com/NRKA/Cache/internal/controller"
	"github.com/NRKA/Cache/internal/datasource/clientCache"
	"github.com/NRKA/Cache/internal/datasource/clientDatabase"
	"github.com/NRKA/Cache/pkg/cache"
	"github.com/NRKA/Cache/pkg/database"
	"log"
	"os"
	"os/signal"
	"time"
)

const cacheExpiration = 2 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	db := database.NewDb()

	cachee := cache.NewCache()

	clientDb := clientDatabase.NewClient(db)
	clientCachee := clientCache.NewClientCache(clientDb, cachee)
	client := controller.NewClient(clientCachee)

	var bestUser = "bestUser"

	// Создаём запись
	err := client.Set(ctx, "user-12345-profile", bestUser, cacheExpiration)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Получаем запись из кэша
	got, err := client.Get(ctx, "user-12345-profile")
	if err != nil {
		log.Fatalf(err.Error())
	}

	if got != bestUser {
		log.Fatalf("not the same value")
	}

	select {
	case <-time.After(cacheExpiration):
	case <-ctx.Done():
		return
	}

	// Получаем запись из базы данных и обновляем кэщ
	gotAgain, err := client.Get(ctx, "user-12345-profile")
	if err != nil {
		log.Fatalf("failed to get value from database: %v", err)
	}

	if gotAgain != bestUser {
		log.Fatalf("not the same value")
	}
}
