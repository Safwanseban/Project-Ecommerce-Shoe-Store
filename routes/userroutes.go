package routes

import (
	c "github.com/Safwanseban/Project-Ecommerce/controllers"
	"github.com/Safwanseban/Project-Ecommerce/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(ctx *gin.Engine) {
	user := ctx.Group("/user")
	{
		user.GET("/", c.UserHome)
		user.POST("/add-to-wishlist", middlewares.UserAuth(), c.AddtoWishList)
		user.GET("/view-wishlist", middlewares.UserAuth(), c.ViewWishlist)
		user.DELETE("/remove-wishlist", middlewares.UserAuth(), c.RemoveFromWishlist)
		user.POST("/wishlist-to-cart", middlewares.UserAuth(), c.WishListToCart)
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
		user.GET("/profile/wallet", middlewares.UserAuth(), c.WalletBalance)
	}

	ctx.POST("/signup", c.Signup)
	ctx.POST("/login", c.LoginUser)
	ctx.POST("/login/otp", c.OtpLog)
	ctx.POST("/login/otpvalidate", c.CheckOtp)
	ctx.GET("/products", c.ProductsView)
	ctx.GET("/products/:id", c.GetProductByID)
	ctx.GET("/payment-success", middlewares.UserAuth(), c.RazorpaySuccess)
	ctx.GET("/success", c.Success)
}
