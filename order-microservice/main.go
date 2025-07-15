package main

import (
	"context"
	"log"
	"order/config"
	pubsuborder "order/pub-sub"
	"order/repository"
	"order/routes"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}
}

func main() {
	app := iris.New()
	app.Use(iris.Compression)

	db, err := config.ConnectWithGORM()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		log.Fatalf("PubSub Client error: %v", err)
	}

	pubsuborder.ListenPaymentSuccess(client, repository.NewOrderRepository(db))

	app.Get("/test", func(ctx iris.Context) {
		ctx.HTML("server is running")
	})

	routes.OrderRoute(app, db)
	app.Listen(":8080")
}
