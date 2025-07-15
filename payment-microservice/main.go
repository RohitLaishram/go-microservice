package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"payment/config"
	"payment/pubsub"
	"payment/repository"
	"syscall"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type OrderEvent struct {
	OrderID uuid.UUID `json:"order_id"`
	UserID  uint      `json:"user_id"`
	Amount  float64   `json:"amount"`
}

func main() {
	_ = godotenv.Load()
	db, err := config.ConnectWithGORM()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	repo := repository.NewPaymentRepository(db)

	go func() {
		if err := pubsub.StartSubscriber(ctx, repo); err != nil {
			log.Fatalf("❌ Pub/Sub error: %v", err)
		}
	}()

	// Block forever (should never reach here)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("🛑 Gracefully shutting down...")
	cancel()
}
