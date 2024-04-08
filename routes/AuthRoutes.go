package routes

import (
	"apna-restaurant-2.0/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	authController *controllers.AuthController
	engine         *gin.Engine
}

func NewAuthRoutes(authController *controllers.AuthController, engine *gin.Engine) *AuthRoutes {
	return &AuthRoutes{
		authController: authController,
		engine:         engine,
	}
}

func (rc *AuthRoutes) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/signup", rc.authController.SignUpUser)
	router.POST("/signin", rc.authController.SignInUser)
}