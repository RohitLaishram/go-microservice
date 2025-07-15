package pubsub

import (
	"context"
	"encoding/json"
	"log"
	"order/repository"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
)

type OrderEventWithStatus struct {
	OrderID uuid.UUID `json:"order_id"`
	UserID  uint      `json:"user_id"`
	Amount  float64   `json:"amount"`
	Status  string    `json:"status"`
}

func ListenPaymentSuccess(client *pubsub.Client, repo *repository.OrderRepository) {
	subName := os.Getenv("PAYMENT_SUBSCRIPTION")
	if subName == "" {
		log.Println("❌ PUBSUB_SUBSCRIPTION not set")
		return
	}
	ctx := context.Background()

	sub := client.Subscription(subName)
	go func() {
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			var event OrderEventWithStatus
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Printf("❌ Failed to unmarshal message: %v", err)
				msg.Nack()
				return
			}

			if event.Status == "success" {
				if err := repo.UpdateStatus(event.OrderID, "PAID"); err != nil {
					log.Printf("❌ Failed to update order status: %v", err)
					msg.Nack()
					return
				}
				log.Printf("✅ Order %s marked as PAID", event.OrderID)
			}
			msg.Ack()
		})
		if err != nil {
			log.Printf("❌ sub.Receive error: %v", err)
		}
	}()
}
