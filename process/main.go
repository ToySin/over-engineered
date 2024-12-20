package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	env "github.com/caarlos0/env/v11"

	"github.com/Toysin/terraform-practice/process/database"
	"github.com/Toysin/terraform-practice/redis"
)

func main() {
	// Create a database client
	var dbConfig database.Config
	if err := env.Parse(&dbConfig); err != nil {
		slog.Error("failed to parse database config", slog.String("reason", err.Error()))
		return
	}
	db, err := database.NewSQLClient(&dbConfig)
	if err != nil {
		slog.Error("failed to create database client", slog.String("reason", err.Error()))
		return
	}
	defer db.Close()

	// Create a redis client
	var redisConfig redis.Config
	if err := env.Parse(&redisConfig); err != nil {
		slog.Error("failed to parse Redis config", slog.String("reason", err.Error()))
		return
	}
	rdb, err := redis.NewPubSubClient(&redisConfig)
	if err != nil {
		slog.Error("failed to create Redis client", slog.String("reason", err.Error()))
		return
	}
	defer rdb.Close()

	// Exit signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := rdb.SubscribeMessage(ctx)
	for {
		select {
		case message := <-ch:
			slog.Info("received message", slog.String("message", message))
			if err := db.SaveMessage(message); err != nil {
				slog.Error("failed to save message", slog.String("reason", err.Error()))
			}
		case <-sigChan:
			slog.Info("received signal, exiting...")
			return
		}
	}
}
