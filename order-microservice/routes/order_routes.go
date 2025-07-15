package routes

import (
	"order/controller"
	"order/repository"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

func OrderRoute(app *iris.Application, db *gorm.DB) {
	orderRepo := repository.NewOrderRepository(db)
	oc := controller.NewOrderController(orderRepo)
	oc.RegisterRoutes(app)
}
