package controllers

import (
	"fmt"

	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

func RazorPay(c *gin.Context) {
	var user models.User
	var cart models.Cart
	useremail := c.GetString("user")
	fmt.Println(useremail)
	initializers.DB.Raw("select id from users where email=?", useremail).Scan(&user)
	initializers.DB.Raw("select sum(total_price from carts where user_id=?", user.ID).Scan(&cart)
	fmt.Println(cart.Total_Price)
	client := razorpay.NewClient("rzp_live_UfVVryLMHd94Xr", "rzp_live_UfVVryLMHd94Xr")

	data := map[string]interface{}{
		"amount":   cart.Total_Price,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	fmt.Println(body)
	if err != nil {
		c.JSON(404, gin.H{
			"msg": "error creating order",
		})
	}
	c.HTML(200,"app.html",gin.H{
		"total":cart.Total_Price,
	})
}

//rzp_live_UfVVryLMHd94Xr
//rzp_live_UfVVryLMHd94Xr
