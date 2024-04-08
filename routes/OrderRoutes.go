package routes

import (
	"apna-restaurant-2.0/controllers"
	"apna-restaurant-2.0/middleware"
	"github.com/gin-gonic/gin"
)

type OrderRoutes struct {
	orderController *controllers.OrderController
	engine          *gin.Engine
}

func NewOrderRoutes(orderController *controllers.OrderController, engine *gin.Engine) *OrderRoutes {
	return &OrderRoutes{
		orderController: orderController,
		engine:          engine,
	}
}

func (or *OrderRoutes) OrderRoute(rg *gin.RouterGroup) {
	router := rg.Group("")
	router.Use(middleware.Authenticate())
	router.POST("/table/new", or.orderController.AddTable)
	router.PATCH("/table/update", or.orderController.UpdateTable)
	router.GET("/table/all", or.orderController.GetAllTables)

	router.POST("/order/new", or.orderController.AddOrder)
	router.GET("/order/all", or.orderController.GetAllOrders)

	router.PATCH("/order/update", or.orderController.UpdateOrder)
	router.GET("/order/:id", or.orderController.GetOrderDetails)
	router.DELETE("/order/:id", or.orderController.CancelOrder)
	router.GET("/orders/:table_id", or.orderController.GetOrderDetailsForTable)
}
