package main

import (
	docs "github.com/Safwanseban/Project-Ecommerce/docs"
	i "github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/routes"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func init() {
	i.ConnecttoDb()
	R.LoadHTMLGlob("templates/*.html")

	i.Getenv()

}

var R = gin.Default()

func main() {
	docs.SwaggerInfo.Title = "E-Commerce API"
	docs.SwaggerInfo.Description = "An e-commerce API which is purely written in GO."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:3000"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}


	routes.AdminRooutes(R)
	routes.UserRoutes(R)
	R.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	R.Run()

}
