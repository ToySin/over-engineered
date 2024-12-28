package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	env "github.com/caarlos0/env/v11"

	"github.com/Toysin/terraform-practice/database"
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

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read message from client
		message := r.FormValue("message")
		if message == "" {
			http.Error(w, "message cannot be empty", http.StatusBadRequest)
			return
		}

		// Publish to Redis
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := rdb.PublishMessage(ctx, message)
		if err != nil {
			http.Error(w, "failed to publish message", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "message published: %s", message)
	})

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get messages from database
		messages, err := db.GetMessages()
		if err != nil {
			http.Error(w, "failed to get messages", http.StatusInternalServerError)
			return
		}

		for _, msg := range messages {
			fmt.Fprintf(w, "Sequence: %d, Message: %s\n", msg.Sequence, msg.Message)
		}
	})

	slog.Info("server started", slog.String("port", "8080"))
	http.ListenAndServe(":8080", nil)
}
