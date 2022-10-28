package controllers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

type Orderd_Items []struct {
	UserId         uint
	Product_id     uint
	OrdersID       string
	Product_Name   string
	Price          string
	Order_Status   string
	Payment_Status string
	PaymentMethod  string
	Total_amount   uint
}

func CreateOrderId() string {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(9999999999-1000000000) + 1000000000
	id := strconv.Itoa(value)
	orderID := "OID" + id
	return orderID
}

func ViewOrders(c *gin.Context) {
	fmt.Println("hai")
	var user models.User
	var ordered_items Orderd_Items
	userEmail := c.GetString("user")
	initializers.DB.Raw("select id from users where email=?", userEmail).Scan(&user)
	record := initializers.DB.Raw("select user_id,product_id,product_name,price,orders_id,order_status,payment_status,payment_method,total_amount from orderd_items where user_id =?", user.ID).Scan(&ordered_items)
	if record.Error != nil {
		c.JSON(404, gin.H{"err": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"orderes": ordered_items})

}

func Cancelorders(c *gin.Context) {
	var user models.User
	var ordered_items Orderd_Items
	userEmail := c.GetString("user")
	orderid := c.Query("orderID")
	initializers.DB.Raw("select id from users where email=?", userEmail).Scan(&user)
	record := initializers.DB.Raw("select user_id,product_id,product_name,price,orders_id,order_status,payment_status,payment_method,total_amount from orderd_items where user_id =?", user.ID).Scan(&ordered_items)
	if record.Error != nil {
		c.JSON(404, gin.H{"err": record.Error.Error()})
		c.Abort()
		return
	}
	


	initializers.DB.Raw("update orderd_items set order_status=? where orders_id=?", "order cancelled", orderid).Scan(&ordered_items)
	c.JSON(200, gin.H{"orderes": ordered_items})
}
