package controller

import (
	"log"
	"order/model"
	pubsub "order/pub-sub"
	"order/repository"

	"github.com/google/uuid"

	"github.com/kataras/iris/v12"
)

type OrderController struct {
	Repo *repository.OrderRepository
}

func NewOrderController(repo *repository.OrderRepository) *OrderController {
	if repo == nil {
		log.Fatal("OrderRepository is null")
	}
	return &OrderController{Repo: repo}
}

func (oc *OrderController) RegisterRoutes(app *iris.Application) {
	order := app.Party("/orders")
	{
		order.Post("/", oc.CreateOrder)
		order.Get("/{id:string}", oc.GetOrder)
		order.Patch("/{id:string}/paid", oc.PaidOrder)
		order.Patch("/{id:string}/status", oc.UpdateOrderStatus)
	}
}

func (oc *OrderController) PaidOrder(ctx iris.Context) {
	id := ctx.Params().Get("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid order ID"})
		return
	}

	order, err := oc.Repo.Get(uuidID)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Order not found"})
		return
	}

	pubErr := pubsub.PublishOrderEvent(pubsub.OrderEvent{
		OrderID: order.ID,
		UserID:  order.UserID,
		Amount:  order.TotalAmount,
	})

	if pubErr != nil {
		log.Println("❌ Failed to publish Pub/Sub message:", pubErr)
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Order payment processed", "orderID": order.ID, "userID": order.UserID, "amount": order.TotalAmount})
}

// POST /orders
func (oc *OrderController) CreateOrder(ctx iris.Context) {
	var order model.Order
	if err := ctx.ReadJSON(&order); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid input"})
		return
	}

	newOrder, err := oc.Repo.Create(order)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to create order"})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(newOrder)
}

// GET /orders/:id
func (oc *OrderController) GetOrder(ctx iris.Context) {
	idParam := ctx.Params().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid order ID"})
		return
	}

	order, err := oc.Repo.Get(id)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Order not found"})
		return
	}

	ctx.JSON(order)
}

// PATCH /orders/:id/status
func (oc *OrderController) UpdateOrderStatus(ctx iris.Context) {
	idParam := ctx.Params().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid order ID"})
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := ctx.ReadJSON(&body); err != nil || body.Status == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid status"})
		return
	}

	if err := oc.Repo.UpdateStatus(id, body.Status); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to update status"})
		return
	}

	ctx.JSON(iris.Map{"message": "Order status updated"})
}
