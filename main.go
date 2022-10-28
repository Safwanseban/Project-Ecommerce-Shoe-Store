package main

import (
	c "github.com/Safwanseban/Project-Ecommerce/controllers"
	i "github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	i.ConnecttoDb()
	R.LoadHTMLGlob("templates/*.html")
	i.Getenv()

}

var R = gin.Default()

func main() {
	admin := R.Group("/admin")
	{
		admin.POST("/login", c.AdminLogin)
		admin.GET("/home", middlewares.AdminAuth(), c.AdminHome)
		// admin.GET("/logout", c.AdminLogout)
		admin.GET("/userdata", middlewares.AdminAuth(), c.Userdata)
		admin.PUT("/userdata/block/:id", middlewares.AdminAuth(), c.BlockUser)
		admin.PUT("/userdata/unblock/:id", middlewares.AdminAuth(), c.UnBlockUser)
		admin.POST("/addproduct", middlewares.AdminAuth(), c.ProductAdding)
		admin.PUT("/editproduct/:id", middlewares.AdminAuth(), c.Editproducts)
		admin.DELETE("/deleteproduct/:id", middlewares.AdminAuth(), c.DeleteProductById)
		admin.GET("/show-orders", middlewares.AdminAuth(), c.AdminShowOrders)
		admin.PUT("/cancel-orders", middlewares.AdminAuth(), c.AdminCancelOrders)
	}
	user := R.Group("/user").Use(middlewares.UserAuth())
	{
		user.GET("/", c.UserHome)

		user.GET("/addtocart", c.AddToCart)
		user.GET("/cart", c.Viewcart)
		user.GET("/checkout", c.Checkout)
		user.POST("/checkout/add-address", c.CheckOutAddAdress)
		// user.GET("/razorpay", c.RazorPay)

		user.GET("/profile", c.UserprofileGet)
		user.POST("/profile/edit", c.UserprofilePost)
		user.GET("/profile/address", c.ShowAddress)
		user.POST("/profile/add-address", c.AddAddress)
		user.PUT("/profile/cancel-order", c.Cancelorders)
		user.GET("/profile/view-order", c.ViewOrders)
		user.POST("/profile/change-password", c.ForgetPassword)
	}
	R.POST("/signup", c.Signup)
	R.POST("/login", c.LoginUser)
	R.POST("/login/otp", c.OtpLog)
	R.POST("/login/otpvalidate", c.ValidateOtp)
	R.GET("/products", c.ProductsView)
	R.GET("/products/:id", c.GetProductByID)
R.GET("/razorpay",c.RazorPay)
	R.Run()

}
