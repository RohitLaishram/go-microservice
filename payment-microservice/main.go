package main

import (
	"context"
	"log"
	"payment/config"
	"payment/pubsub"
	"payment/repository"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

type OrderEvent struct {
	OrderID uuid.UUID `json:"order_id"`
	UserID  uint      `json:"user_id"`
	Amount  float64   `json:"amount"`
}

func main() {
	app := iris.New()
	app.Use(iris.Compression)

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
	app.Get("/test", func(ctx iris.Context) {
		ctx.HTML("server is running")
	})

	app.Listen(":8081")
}
