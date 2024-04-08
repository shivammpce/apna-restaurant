package app

import "apna-restaurant-2.0/routes"

func mapUrls(authRoute *routes.AuthRoutes, menuRoute *routes.MenuRoutes, orderRoute *routes.OrderRoutes) {
	authRoute.AuthRoute(&router.RouterGroup)
	menuRoute.MenuRoute(&router.RouterGroup)
	orderRoute.OrderRoute(&router.RouterGroup)
}
