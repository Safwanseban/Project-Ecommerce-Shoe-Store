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

func ReturnOrder(c *gin.Context) {
	var order models.Orderd_Items
	var user models.User
	userEmail := c.GetString("user")

	initializers.DB.Raw("select id from users where email=?", userEmail).Scan(&user)
	var Order_return struct {
		Order_ID string
	}
	if err := c.BindJSON(&Order_return); err != nil {
		c.JSON(200, gin.H{"err": err.Error()})
	}
	initializers.DB.Where("orders_id=?", Order_return.Order_ID).Find(&order)
	if order.Order_Status == "returned" {
		c.JSON(400, gin.H{
			"msg": "order already returened",
		})
		c.Abort()
		return
	}

	var balance int
	initializers.DB.Raw("SELECT wallet_balance FROM users WHERE id = ?", user.ID).Scan(&balance)
	//starting transcations
	tx := initializers.DB.Begin()
	record1 := tx.Model(&models.Orderd_Items{}).Where("orders_id=?", Order_return.Order_ID).Update("order_status", "returned")
	if record1.Error != nil {
		tx.Rollback()
		c.JSON(404, gin.H{
			"err": record1.Error.Error(),
		})
	}
	newBalance := balance + int(order.Total_amount)
	record2 := tx.Model(&user).Where("id = ?", user.ID).Update("wallet_balance", newBalance)
	if record2.Error != nil {
		tx.Rollback()
		c.JSON(404, gin.H{
			"err": record1.Error.Error(),
		})
	}
	tx.Commit()

	c.JSON(200, gin.H{
		"msg": "order returned",
	})

}

func Cancelorders(c *gin.Context) {
	var user models.User
	userEmail := c.GetString("user")
	orderid := c.Query("orderID")
	update_status := "order cancelled"
	initializers.DB.Raw("select id from users where email=?", userEmail).Scan(&user) //getting user id
	var orders models.Orderd_Items
	// initializers.DB.First(&orders).Where("orders_id=?", orderid)
	initializers.DB.Where("orders_id=?", orderid).Find(&orders)
	fmt.Println(orders.Order_Status, orders.OrdersID)
	fmt.Println(orderid)
	if orders.Order_Status == update_status {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "Order already Cancelled",
		})
		return
	}
	initializers.DB.Raw("update orderd_items set order_status=? where orders_id=?", "order cancelled", orderid).Scan(&orders)
	var price int
	initializers.DB.Raw("SELECT total_amount FROM orderd_items WHERE orders_id = ?", orderid).Scan(&price)
	var balance int
	initializers.DB.Raw("SELECT wallet_balance FROM users WHERE id = ?", user.ID).Scan(&balance)
	newBalance := price + balance
	initializers.DB.Raw("update users set wallet_balance=? where id=?", newBalance, user.ID).Scan(&user)
	c.JSON(200, gin.H{
		"msg": "order canccelled",
	})

}
