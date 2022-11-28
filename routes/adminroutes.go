package routes

import (
	c "github.com/Safwanseban/Project-Ecommerce/controllers"
	"github.com/Safwanseban/Project-Ecommerce/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRooutes(ctx *gin.Engine) {
	admin := ctx.Group("/admin")
	{
		admin.POST("/signup", c.AdminSignup)
		admin.POST("/login", c.AdminLogin)
		admin.GET("/home", middlewares.AdminAuth(), c.AdminHome)
		// admin.GET("/logout", c.AdminLogout)
		admin.GET("/userdata", middlewares.AdminAuth(), c.Userdata)
		admin.PUT("/userdata/block/:id", middlewares.AdminAuth(), c.BlockUser)
		admin.PUT("/userdata/unblock/:id", middlewares.AdminAuth(), c.UnBlockUser)
		admin.GET("/getcatogories", middlewares.AdminAuth(), c.ListingAllCat)
		admin.PUT("/apply-brand-discount", middlewares.AdminAuth(), c.ApplyBrandDiscount)

		admin.POST("/addproduct", middlewares.AdminAuth(), c.ProductAdding)
		admin.PUT("/editproduct/:id", middlewares.AdminAuth(), c.Editproducts)
		admin.DELETE("/deleteproduct/:id", middlewares.AdminAuth(), c.DeleteProductById)
		admin.GET("/show-orders", middlewares.AdminAuth(), c.AdminShowOrders)
		admin.PUT("/cancel-orders", middlewares.AdminAuth(), c.AdminCancelOrders)
		admin.PUT("/change-orderstatus", middlewares.AdminAuth(), c.AdminChangeOrderStatus)
		admin.POST("/generate-coupon", middlewares.AdminAuth(), c.GenerateCoupon)
	}
}
