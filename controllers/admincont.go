package controllers

import (
	"fmt"
	"net/http"

	"github.com/Safwanseban/Project-Ecommerce/auth"
	"github.com/Safwanseban/Project-Ecommerce/initializers"

	i "github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// ShowAccount godoc
//
//	@Summary		Show an account
//	@Description	get string by ID
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Account ID"
//	@Success		200	{object}	model.Account
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/accounts/{id} [get]
type AdminLogins struct {
	Email    string
	Password string
}

var UserDb = map[string]string{
	"email":    "safwan@gmail.com",
	"password": "safwan",
}

// AdminSignup godoc
//	
//	@Summary		Show an account
//	@Description	get string by ID
//	@Tags			admin/signup
//	@Accept			json
//	@Produce		json
//
//	@Success		200	{object}	models.Admin
//	@Router			/admin/signup [post]
func AdminSignup(c *gin.Context) {
	var admin models.Admin
	var count uint
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(404, gin.H{
			"err": err.Error(),
		})
		c.Abort()
		return
	}

	i.DB.Raw("select count(*) from admins where email=?", admin.Email).Scan(&count)
	if count > 0 {
		c.JSON(400, gin.H{
			"status": "false",
			"msg":    "an admin with same email already exists",
		})
		c.Abort()
		return
	}
	if err := admin.HashPassword(admin.Password); err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
	}
	record := i.DB.Create(&admin)
	if record.Error != nil {
		c.JSON(400, gin.H{
			"err": record.Error.Error(),
		})
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "new admin created",
	})
}
func AdminLogin(c *gin.Context) { // admin login page post
	var u AdminLogins
	var admin models.Admin
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})

		c.Abort()
		return
	}
	record := initializers.DB.Raw("select * from admins where email=?", u.Email).Scan(&admin)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	credentialcheck := admin.CheckPassword(u.Password)
	if credentialcheck != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}

	tokenstring, err := auth.GenerateJWT(u.Email) //generating a jwt

	token := tokenstring["access_token"]
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Adminjwt", token, 3600*24*30, "", "", false, true)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      true,
		"message":     "ok",
		"tokenstring": tokenstring,
	})

}

func AdminHome(c *gin.Context) {
	fmt.Println("hai")
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

type Userdet []struct {
	ID           uint
	First_Name   string
	Last_Name    string
	Email        string
	Phone        string
	Block_status bool
	Country      string
	City         string
	Pincode      uint
}

func Userdata(c *gin.Context) {

	var user Userdet
	i.DB.Raw("SELECT id,first_name,last_name,email,phone,block_status,country,city,pincode FROM users ORDER BY id ASC").Scan(&user)
	if search := c.Query("search"); search != "" {
		i.DB.Raw("SELECT id,first_name,last_name,email,phone,block_status,country,city,pincode FROM users where first_name like ? ORDER BY id ASC ", search).Scan(&user)
	}

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
	record := i.DB.Raw("select user_id,product_id,product_name,price,applied_coupons,orders_id,order_status,payment_status,payment_method,total_amount from orderd_items ").Scan(&ordered_items)
	if search := c.Query("search"); search != "" {
		record := i.DB.Raw("select user_id,product_id,product_name,price,applied_coupons,orders_id,order_status,payment_status,payment_method,total_amount from orderd_items where  (product_name ilike ? or payment_method ilike ?)", "%"+search+"%", "%"+search+"%").Scan(&ordered_items)
		if record.Error != nil {
			c.JSON(404, gin.H{"err": record.Error.Error()})
			c.Abort()
			return
		}
	}

	if record.Error != nil {
		c.JSON(404, gin.H{"err": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"orderes": ordered_items})
}

func AdminChangeOrderStatus(c *gin.Context) {

	var Order_status struct {
		Order_Id     string
		Order_Status string
	}
	if err := c.BindJSON(&Order_status); err != nil {
		c.JSON(404, gin.H{
			"err": err.Error(),
		})
	}
	record := initializers.DB.Model(&models.Orderd_Items{}).Where("orders_id = ?", Order_status.Order_Id).Update("order_status", Order_status.Order_Status)
	if record.Error == nil {
		c.JSON(200, gin.H{
			"order_status": Order_status.Order_Status,
			"msg":          "order status have been successfully changed",
		})
	}

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
