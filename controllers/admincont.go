package controllers

import (
	"net/http"

	"github.com/Safwanseban/Project-Ecommerce/auth"
	i "github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

type AdminLogins struct {
	Email    string
	Password string
}

var UserDb = map[string]string{
	"email":    "safwan@gmail.com",
	"password": "safwan",
}

func AdminLogin(c *gin.Context) { // admin login page post
	var u AdminLogins

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})

		c.Abort()
		return
	}

	if UserDb["email"] == u.Email && UserDb["password"] == u.Password && c.Request.Method == "POST" {

		// session, err := i.Store.Get(c.Request, "admin")
		// session.Values["emails"] = u.Email
		// fmt.Println(session.Values["emails"])
		// session.Save(c.Request, c.Writer)
		// fmt.Println(err)
		tokenstring, err := auth.GenerateJWT(u.Email)
		c.SetCookie("Adminjwt", tokenstring, 3600*24*30, "", "", false, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "ok",
			"data":    tokenstring,
		})

	} else {
		c.JSON(404, gin.H{"msg": "error login"})
	}

}

func AdminHome(c *gin.Context) {

	// ok := middlewares.IsAdminloggeedin(c)
	// if !ok {
	// 	c.JSON(404, gin.H{
	// 		"msg": "session is not there",
	// 	})
	// } else {
	c.JSON(202, gin.H{
		"msg": "Welcome to Admin Panel",
	})

}
func AdminLogout(c *gin.Context) { // adminLogout page

	// session, err := i.Store.Get(c.Request, "admin")
	// session.Options.MaxAge = -1
	// fmt.Println(session.Values["emails"])
	// session.Save(c.Request, c.Writer)
	// fmt.Println(err)

}

func Userdata(c *gin.Context) {
	var user []models.User
	i.DB.Raw("SELECT * FROM users ORDER BY id ASC").Scan(&user)
	c.JSON(200, gin.H{"user": user})
}
func BlockUser(c *gin.Context) {
	params := c.Param("id")
	var user models.User
	i.DB.Raw("UPDATE users SET block_status=true where id=?", params).Scan(&user)
	c.JSON(http.StatusOK, gin.H{"msg": "Blocked succesfully"})
}
func UnBlockUser(c *gin.Context) {
	params := c.Param("id")
	var user models.User
	i.DB.Raw("UPDATE users SET block_status=false where id=?", params).Scan(&user)
	c.JSON(http.StatusOK, gin.H{"msg": "Unblocked succesfully"})
}
func AdminShowOrders(c *gin.Context) {
	var ordered_items Orderd_Items
	record := i.DB.Raw("select user_id,product_id,product_name,price,orders_id,order_status,payment_status,payment_method,total_amount from orderd_items ").Scan(&ordered_items)
	if record.Error != nil {
		c.JSON(404, gin.H{"err": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"orderes": ordered_items})
}
func AdminCancelOrders(c *gin.Context) {

	var ordered_items Orderd_Items
	orderid := c.Query("orderID")
	record := i.DB.Raw("select user_id,product_id,product_name,price,orders_id,order_status,payment_status,payment_method,total_amount from orderd_items").Scan(&ordered_items)
	if record.Error != nil {
		c.JSON(404, gin.H{"err": record.Error.Error()})
		c.Abort()
		return
	}

	i.DB.Raw("update orderd_items set order_status=? where orders_id=?", "order cancelled", orderid).Scan(&ordered_items)
	c.JSON(200, gin.H{"orderes": ordered_items})
}
