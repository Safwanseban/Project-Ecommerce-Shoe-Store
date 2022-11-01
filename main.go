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
	user := R.Group("/user")
	{
		user.GET("/", c.UserHome)

		user.POST("/addtocart", middlewares.UserAuth(), c.AddToCart)
		user.GET("/cart", middlewares.UserAuth(), c.Viewcart)
		user.GET("/checkout", middlewares.UserAuth(), c.Checkout)
		user.POST("/checkout/add-address", middlewares.UserAuth(), c.CheckOutAddAdress)
		//raxorpay related
		user.GET("/razorpay", middlewares.UserAuth(), c.RazorPay)

		user.GET("/profile", middlewares.UserAuth(), c.UserprofileGet)
		user.POST("/profile/edit", middlewares.UserAuth(), c.UserprofilePost)
		user.GET("/profile/address", middlewares.UserAuth(), c.ShowAddress)
		user.POST("/profile/add-address", middlewares.UserAuth(), c.AddAddress)
		user.PUT("/profile/cancel-order", middlewares.UserAuth(), c.Cancelorders)
		user.GET("/profile/view-order", middlewares.UserAuth(), c.ViewOrders)
		user.POST("/profile/change-password", middlewares.UserAuth(), c.ForgetPassword)
	}
	R.POST("/signup", c.Signup)
	R.POST("/login", c.LoginUser)
	R.POST("/login/otp", c.OtpLog)
	R.POST("/login/otpvalidate", c.CheckOtp)
	R.GET("/products", c.ProductsView)
	R.GET("/products/:id", c.GetProductByID)
	R.GET("/payment-success", middlewares.UserAuth(), c.RazorpaySuccess)
	R.GET("/success", c.Success)
	R.Run()

}
