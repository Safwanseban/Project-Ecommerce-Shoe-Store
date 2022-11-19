package main

import (
	i "github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	i.ConnecttoDb()
	R.LoadHTMLGlob("templates/*.html")
	R.LoadHTMLFiles("")
	i.Getenv()

}

var R = gin.Default()

func main() {

	routes.AdminRooutes(R)
	routes.UserRoutes(R)
	R.Run()

}
