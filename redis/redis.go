package redis

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

const (
	messageTopic = "message"
)

// Config is a struct that contains Redis configuration
type Config struct {
	Addr string `envconfig:"REDIS_ADDR" envDefault:"localhost:6379"`
}

func createRedisConnection(cfg *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
	status := client.Ping(context.Background())
	if status.Err() != nil {
		return nil, status.Err()
	}
	return client, nil
}

// PubSubClient is a struct that provides Pub/Sub functions
type PubSubClient struct {
	client *redis.Client
}

// NewPubSubClient is a function to create a new PubSubClient with the given configuration
func NewPubSubClient(cfg *Config) (*PubSubClient, error) {
	client, err := createRedisConnection(cfg)
	return &PubSubClient{client: client}, err
}

func (p *PubSubClient) Close() error {
	return p.client.Close()
}

// PublishMessage is a function to publish a message to the Redis Pub/Sub channel
func (p *PubSubClient) PublishMessage(ctx context.Context, message string) error {
	return p.client.Publish(ctx, messageTopic, message).Err()
}

// SubscribeMessage is a function to subscribe to the Redis Pub/Sub channel
func (p *PubSubClient) SubscribeMessage(ctx context.Context) chan string {
	pubsub := p.client.Subscribe(ctx, messageTopic)
	ch := make(chan string)
	go func() {
		defer func() {
			pubsub.Close()
			close(ch)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-pubsub.Channel():
				ch <- msg.Payload
			}
		}
	}()
	return ch
}
