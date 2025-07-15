package pubsub

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
)

type OrderEvent struct {
	OrderID uuid.UUID `json:"order_id"`
	UserID  uint      `json:"user_id"`
	Amount  float64   `json:"amount"`
}

func PublishOrderEvent(event OrderEvent) error {
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	topicID := os.Getenv("PUBSUB_TOPIC")

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	topic := client.Topic(topicID)
	defer topic.Stop()

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	log.Println("📦 Publishing order event:", string(data))

	result := topic.Publish(ctx, &pubsub.Message{Data: data})
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	log.Println("✅ Published with ID:", id)
	return nil
}
