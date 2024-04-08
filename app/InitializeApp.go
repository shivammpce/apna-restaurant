package app

import (
	"log"

	"apna-restaurant-2.0/controllers"
	"apna-restaurant-2.0/db/config"
	repo "apna-restaurant-2.0/db/sqlc"
	"apna-restaurant-2.0/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	router *gin.Engine
	db     *repo.Queries
	// AuthController *controllers.AuthController
	// AuthRoutes     *routes.AuthRoutes
)

func init() {
	// db := config.ConnectToDB()
	// router = gin.Default()
	// queries := repo.New(db)
	// AuthController = controllers.NewAuthController(queries)
	// AuthRoutes = routes.NewAuthRoutes(AuthController, router)
}

func StartApplication() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading env data")
	}
	db := config.ConnectToDB()
	router = gin.Default()
	queries := repo.New(db)
	AuthController := controllers.NewAuthController(queries)
	AuthRoutes := routes.NewAuthRoutes(AuthController, router)
	MenuController := controllers.NewMenuController(queries)
	MenuRoutes := routes.NewMenuRoutes(MenuController, router)
	OrderController := controllers.NewOrderController(queries)
	OrderRoutes := routes.NewOrderRoutes(OrderController, router)
	mapUrls(AuthRoutes, MenuRoutes, OrderRoutes)
	router.Run(":8080")
}
