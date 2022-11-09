package controllers

import (
	"strconv"

	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

func AddtoWishList(c *gin.Context) {
	var user models.User
	UserEmail := c.GetString("user")
	initializers.DB.Where("email=?", UserEmail).Find(&user)

	var WishList struct {
		Product_id uint
	}
	if err := c.BindJSON(&WishList); err != nil {
		c.JSON(404, gin.H{
			"err": err.Error(),
		})
	}
	var count int64
	initializers.DB.Model(&models.WishList{}).Where("product_id = ? and user_id=?", WishList.Product_id, user.ID).Count(&count)
	if count >= 1 {
		c.JSON(400, gin.H{
			"msg": "product already added to wishlist",
		})
		c.Abort()
		return
	}

	wishlist := models.WishList{
		UserID:     user.ID,
		Product_id: WishList.Product_id,
	}

	initializers.DB.Create(&wishlist)
	c.JSON(200, gin.H{
		"status": "true",
		"msg":    "product added to wishlist",
	})
}

var Wishlist []struct {
	Product_id   uint
	User_id      uint
	Product_name string
	Price        string
}

func ViewWishlist(c *gin.Context) {
	var user models.User
	UserEmail := c.GetString("user")
	initializers.DB.Where("email=?", UserEmail).Find(&user)
	initializers.DB.Raw("select users.id as user_id, products.product_id,products.product_name,products.price from wish_lists join products on wish_lists.product_id=products.product_id join users on wish_lists.user_id=users.id where user_id=? ", user.ID).Scan(&Wishlist)
	c.JSON(200, gin.H{
		"Wishlist": Wishlist,
	})
}
func RemoveFromWishlist(c *gin.Context) {

	productid := c.Query("productID")
	var user models.User
	UserEmail := c.GetString("user")
	initializers.DB.Where("email=?", UserEmail).Find(&user)
	var wishlist models.WishList
	initializers.DB.Raw("delete from wish_lists where product_id=? and user_id=?", productid, user.ID).Scan(&wishlist)
	c.JSON(200, gin.H{
		"msg": "product deleted from wishlist",
	})
}
func WishListToCart(c *gin.Context) {
	productid := c.Query("productID")
	productID, _ := strconv.Atoi(productid)
	var user models.User
	UserEmail := c.GetString("user")
	initializers.DB.Where("email=?", UserEmail).Find(&user)
	var count int64
	initializers.DB.Model(&models.WishList{}).Where("product_id = ? and user_id=?", productid, user.ID).Count(&count)
	if count <= 0 {
		c.JSON(400, gin.H{
			"msg": "item not available in wishlist",
		})
		c.Abort()
		return
	}
	var total int64
	var price uint
	initializers.DB.Raw("select price from products where product_id=?", productid).Scan(&price)
	initializers.DB.Model(&models.Cart{}).Where("product_id = ? and user_id=?", productid, user.ID).Count(&total)
	if total >= 1 {
		c.JSON(400, gin.H{"msg": "item is already in cart"})

	} else {
		cart := models.Cart{
			ProductID:   uint(productID),
			Quantity:    1,
			UserId:      user.ID,
			Total_Price: price,
		}
		initializers.DB.Create(&cart)
		c.JSON(200, gin.H{
			"msg": "item added to cart",
		})
	}

}
