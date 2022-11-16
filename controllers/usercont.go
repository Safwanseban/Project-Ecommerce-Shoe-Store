package controllers

import (
	"fmt"
	"net/http"

	"github.com/Safwanseban/Project-Ecommerce/auth"
	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(404, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}
	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(404, gin.H{"err": err.Error()})
		c.Abort()
		return
	}
	record := initializers.DB.Create(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": record.Error.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"email": user.Email, "msg": "Go to LoginPage"})

}

type UserLogin struct {
	Email        string
	Password     string
	Block_status bool
}

func LoginUser(c *gin.Context) {
	var ulogin UserLogin
	var user models.User
	if err := c.ShouldBindJSON(&ulogin); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	record := initializers.DB.Raw("select * from users where email=?", ulogin.Email).Scan(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	if user.Block_status {
		c.JSON(404, gin.H{"msg": "user has been blocked By admin"})
		c.Abort()
		return
	}
	credentialcheck := user.CheckPassword(ulogin.Password)
	if credentialcheck != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}
	tokenString, err := auth.GenerateJWT(user.Email)
	fmt.Println(tokenString)
	token:=tokenString["access_token"]	
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAuth", token, 3600*24*30, "", "", false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"email": ulogin.Email, "password": ulogin.Password, "token": tokenString,})
}
func UserHome(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "welcome User Home"})

}

func ForgetPassword(c *gin.Context) {
	var user models.User
	useremail := c.GetString("user")
	newpassword := c.PostForm("password")
	initializers.DB.Raw("select password,id from users where email=?", useremail).Scan(&user)
	// user.Password=newpassword
	if err := user.HashPassword(newpassword); err != nil {
		c.JSON(404, gin.H{"err": err.Error()})
		c.Abort()
		return
	}
	fmt.Println(user.HashPassword(newpassword))

	initializers.DB.Raw("update users set password=? where id=?", user.Password, user.ID).Scan(&user)
	fmt.Println(useremail)

	fmt.Println(newpassword)
	c.JSON(200, gin.H{
		"pass":  newpassword,
		"email": useremail,
	})

}
