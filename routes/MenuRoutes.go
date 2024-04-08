package routes

import (
	"apna-restaurant-2.0/controllers"
	"apna-restaurant-2.0/middleware"
	"github.com/gin-gonic/gin"
)

type MenuRoutes struct {
	menuController *controllers.MenuController
	engine         *gin.Engine
}

func NewMenuRoutes(menuController *controllers.MenuController, engine *gin.Engine) *MenuRoutes {
	return &MenuRoutes{
		menuController: menuController,
		engine:         engine,
	}
}

func (mc *MenuRoutes) MenuRoute(rg *gin.RouterGroup) {
	router := rg.Group("/menu")
	router.Use(middleware.Authenticate())
	router.POST("/new", mc.menuController.AddMenu)
	router.GET("/all", mc.menuController.GetAllMenus)
	router.GET("/:id", mc.menuController.GetMenuByID)
	router.PATCH("/update", mc.menuController.UpdateMenu)

	router.POST("/new-item", mc.menuController.AddMenuItem)
	router.GET("/all-menuitems", mc.menuController.GetAllMenuItems)
	router.GET("/menuitem/:id", mc.menuController.GetMenuitemByID)
	router.PATCH("/update-menuitem", mc.menuController.UpdateMenuitem)
}
