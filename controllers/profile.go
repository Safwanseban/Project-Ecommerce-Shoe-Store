package controllers

import (
	"fmt"
	"strconv"

	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

type Profile struct {
	Email      string
	First_Name string
	Last_name  string
	Phone      uint
	Country    string
	City       string
	Pincode    string
}

func UserprofileGet(c *gin.Context) {
	UserEmail := c.GetString("user")
	var profile Profile
	initializers.DB.Raw("select first_name,email,last_name,phone,city,pincode,country from users where email=?", UserEmail).Scan(&profile)
	c.JSON(200, gin.H{
		"profile": profile,
	})

}
func UserprofilePost(c *gin.Context) {
	UserEmail := c.GetString("user")
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	country := c.PostForm("country")
	city := c.PostForm("city")
	pincod := c.PostForm("pincode")
	pincode, _ := strconv.Atoi(pincod)
	// profile := models.User{
	// 	First_Name: firstname,
	// 	Last_Name:  lastname,
	// 	Country:    country,
	// 	City:       city,
	// 	Pincode:    uint(pincode),
	// }

	var user models.User
	record := initializers.DB.Raw("update users set first_name=?,last_name=?,country=?,city=?,pincode=? where email=?", firstname, lastname, country, city, pincode, UserEmail).Scan(&user)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"err": record.Error.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"msg": "updated successfully"})
}
func ShowAddress(c *gin.Context) {
	var user models.User
	UserEmail := c.GetString("user")
	initializers.DB.Raw("select id from users where email=?", UserEmail).Scan(&user)

	record := initializers.DB.Raw("select user_id,name,phone_number,pincode,house,area,landmark,city from addresses where user_id=?", user.ID).Scan(&Address)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"error": record.Error.Error(),
		})
	}
	c.JSON(200, gin.H{
		"address": Address,
	})

}
func AddAddress(c *gin.Context) { //add address through profile
	useremail := c.GetString("user")
	var user models.User
	initializers.DB.Raw("select id from users where email=?", useremail).Scan(&user)

	Name := c.PostForm("name")
	Phonenum := c.PostForm("phonenumber")
	phonenum, _ := strconv.Atoi(Phonenum)
	pincod := c.PostForm("pincode")
	pincode, _ := strconv.Atoi(pincod)
	area := c.PostForm("area")
	houseadd := c.PostForm("house")
	landmark := c.PostForm("landmark")
	city := c.PostForm("city")
	address := models.Address{
		UserId:       user.ID,
		Name:         Name,
		Phone_number: phonenum,
		Pincode:      pincode,
		Area:         area,
		House:        houseadd,
		Landmark:     landmark,
		City:         city,
	}
	record := initializers.DB.Create(&address)
	if record.Error != nil {
		c.JSON(404, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"msg": "address added"})
}

func WalletBalance(c *gin.Context) {
	useremail := c.GetString("user")
	var user models.User
	initializers.DB.Raw("select id,wallet_balance from users where email=?", useremail).Scan(&user)
	fmt.Println(user.Wallet_Balance)
	c.JSON(200, gin.H{
		"wallet-balance": user.Wallet_Balance,
	})

}
