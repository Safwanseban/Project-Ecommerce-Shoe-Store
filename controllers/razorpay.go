package controllers

import (
	"fmt"
	"strconv"

	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

func RazorPay(c *gin.Context) {
	var user models.User
	// var cart models.Cart
	useremail := c.GetString("user")
	fmt.Println(useremail)
	initializers.DB.Raw("select id,phone from users where email=?", useremail).Scan(&user)
	var sumtotal string
	initializers.DB.Raw("select sum(total_price) from carts where user_id=?", user.ID).Scan(&sumtotal)
	fmt.Println(sumtotal)
	client := razorpay.NewClient("rzp_test_Nfnipdccvgb8fW", "UfwKXCGjiUrcfTEXpWlupcrN")

	data := map[string]interface{}{
		"amount":   sumtotal,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	fmt.Println(body)
	value := body["id"]

	if err != nil {
		c.JSON(404, gin.H{
			"msg": "error creating order",
		})
	}
	c.HTML(200, "app.html", gin.H{

		"UserID":       user.ID,
		"total":        sumtotal,
		"orderid":      value,
		"amount":       sumtotal,
		"Email":        useremail,
		"Phone_Number": user.Phone,
	})
	if err != nil {
		c.JSON(200, gin.H{
			"msg": value,
		})
	}
}

func RazorpaySuccess(c *gin.Context) {
	userid := c.Query("user_id")
	userID, _ := strconv.Atoi(userid)
	orderid := c.Query("order_id")
	paymentid := c.Query("payment_id")
	signature := c.Query("signature")
	id := c.Query("orderid")
	totalamount := c.Query("total")
	Rpay := models.RazorPay{
		UserID:          uint(userID),
		OrderId:         id,
		RazorPaymentId:  paymentid,
		Signature:       signature,
		RazorPayOrderID: orderid,
		AmountPaid:      totalamount,
	}
	initializers.DB.Create(&Rpay)
	var cart models.Cart
	initializers.DB.Raw("delete from carts where user_id=?", userID).Scan(&cart)
	fmt.Println(userID,orderid)
	OrderPlaced(userID,orderid)


}
func Success(c *gin.Context) {
	c.HTML(200, "succs.html", nil)

}
func OrderPlaced(Uid int,orderId string){
userid:=Uid
orderid:=orderId
var orders models.Orders
initializers.DB.Raw("update orders set order_status=?,payment_status=?,order_id=? where user_id=?","order completed","payment done",orderid,userid).Scan(&orders)
var ordereditems models.Orderd_Items

initializers.DB.Raw("update orderd_items set order_status=?,payment_status=? where user_id=?", "orderplaced", "Payment Completed",  userid).Scan(&ordereditems)

}

//rzp_live_UfVVryLMHd94Xr
//fm8HR53htJ2TZJfFJDbQHlU2
