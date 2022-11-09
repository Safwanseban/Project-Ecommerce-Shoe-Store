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

var orderscan struct {
	Total_Amount   uint
	Count          uint
	Order_Status   string
	Payment_Status string
	Payment_Method string
}

// func Cancelorders(c *gin.Context) {
// 	var user models.User
// 	userEmail := c.GetString("user")
// 	orderid := c.Query("orderID")
// 	update_status := "order cancelled"
// 	initializers.DB.Raw("select id from users where email=?", userEmail).Scan(&user) //getting user id

// 	var order models.Orderd_Items

// 	var total_amount uint
// 	initializers.DB.Raw("select total_amount,count(*) as count from orderd_items where orders_id=? and user_id=? and order_status=? group by total_amount", orderid, user.ID, update_status).Scan(&orderscan)
// 	if orderscan.Count >= 1 {
// 		c.JSON(300, gin.H{
// 			"msg": "order already cancelled",
// 		})
// 		c.Abort()
// 		return
// 	}
// 	initializers.DB.Raw("select payment_status,order_status,payment_method from orderd_items where orders_id=? and user_id=?", orderid, user.ID).Scan(&orderscan)
// 	fmt.Println(orderscan.Count, orderscan.Total_Amount, orderscan.Payment_Method)
// 	fmt.Println(total_amount)

// 	fmt.Println(user.ID, userEmail)

// 	record := initializers.DB.Raw("update orderd_items set order_status=? where orders_id=?", "order cancelled", orderid).Scan(&order)

// 	var balance uint
// 	if orderscan.Order_Status == update_status && orderscan.Payment_Method == "Razorpay" && order.Payment_Status == "Payment Completed" {
// 		fmt.Println("hai")
// 		initializers.DB.Raw("select wallet_balance from users where id=?", user.ID).Scan(&balance)

// 		newBalance := balance + orderscan.Total_Amount
// 		initializers.DB.Raw("update users set wallet_balance =? where id=?", newBalance, user.ID).Scan(&user)
// 	}

// 	fmt.Println(user.Wallet_Balance)
// 	if record.Error == nil {
// 		c.JSON(200, gin.H{
// 			"msg": "order cancelled successfully",
// 		})
// 	}
// }

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
