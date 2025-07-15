package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"payment/model"
	"payment/repository"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
)

type PaymentEvent struct {
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
}

type OrderEvent struct {
	OrderID uuid.UUID `json:"order_id"`
	UserID  uint      `json:"user_id"`
	Amount  float64   `json:"amount"`
}

func PublishOrderCreated(ctx context.Context, client *pubsub.Client, order PaymentEvent) {
	topicName := os.Getenv("PAYMENT_TOPIC")
	data, _ := json.Marshal(order)
	topic := client.Topic(topicName)
	result := topic.Publish(ctx, &pubsub.Message{Data: data})
	_, err := result.Get(ctx)
	if err != nil {
		log.Printf("Error publishing order: %v", err)
	}
}

// publishPaymentSuccess publishes a payment success message to a Pub/Sub topic
func publishPaymentSuccess(ctx context.Context, client *pubsub.Client, event OrderEvent) error {
	topicName := os.Getenv("PAYMENT_TOPIC")
	if topicName == "" {
		log.Println("❌ PUBSUB_TOPIC not set")
		return nil
	}
	topic := client.Topic(topicName)
	msgData, err := json.Marshal(map[string]interface{}{
		"type":     "payment-success",
		"order_id": event.OrderID,
		"user_id":  event.UserID,
		"amount":   event.Amount,
		"status":   "success",
	})
	if err != nil {
		return err
	}
	result := topic.Publish(ctx, &pubsub.Message{
		Data: msgData,
	})
	_, err = result.Get(ctx)
	if err != nil {
		log.Printf("❌ Failed to publish payment success: %v", err)
	}
	log.Printf("✅ Payment success published for Order ID: %s", event.OrderID)
	return err
}

// StartSubscriber continuously listens for messages from the subscription
func StartSubscriber(ctx context.Context, repo *repository.PaymentRepository) error {
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	subID := os.Getenv("PUBSUB_SUBSCRIPTION")

	log.Printf("📡 Subscribing to %s in project %s...\n", subID, projectID)

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("PubSub client error: %w", err)
	}

	sub := client.Subscription(subID)

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("✅ Payment Service received: %s\n", string(msg.Data))

		var event OrderEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("❌ Failed to parse message: %v", err)
			msg.Nack()
			return
		}

		newPayment := &model.Payment{
			OrderID: event.OrderID,
			Amount:  event.Amount,
			Status:  "success",
		}

		if err := repo.CreatePayment(newPayment); err != nil {
			log.Printf("❌ Failed to create payment: %v", err)
			msg.Nack()
			return
		}
		log.Printf("✅ Payment created successfully for Order ID: %s", event.OrderID)

		_ = publishPaymentSuccess(ctx, client, event)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("❌ Failed to receive messages: %w", err)
	}

	return nil
}
