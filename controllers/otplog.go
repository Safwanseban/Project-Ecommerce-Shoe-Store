package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/Safwanseban/Project-Ecommerce/auth"
	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func init() {
	initializers.Getenv()
	initializers.ConnecttoDb()
}
func OtpLog(c *gin.Context) {

	Mob := c.Query("number")
	var (
		accountSid string
		authToken  string
		fromPhone  string
		client     *twilio.RestClient
	)

	result := ChekNumber(Mob)
	fmt.Println(result)

	if !result {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "Mobile number doesnt exist! Please SignUp",
		})
		return
	}

	// Get Twillio credentials from .env file
	accountSid = os.Getenv("ACCOUNT_SID")
	authToken = os.Getenv("AUTH_TOKEN")
	fromPhone = os.Getenv("FROM_PHONE")

	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	Mobile := "+91" + Mob
	// Creatin 4 digit OTP
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(9999-1000) + 1000
	otp := strconv.Itoa(value)
	otpDb := models.Otp{Mobile: Mob, Otp: otp}
	initializers.DB.Create(&otpDb)

	params := openapi.CreateMessageParams{}
	params.SetTo(Mobile)
	params.SetFrom(fromPhone)
	params.SetBody("Your OTP - " + otp)

	_, err := client.Api.CreateMessage(&params)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "error sending OTP",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": "OTP Sent Succesfully",
	})

}

// Checking if the number belongs to any user
func ChekNumber(str string) bool {

	mobilenumber := str
	var checkOtp models.User
	initializers.DB.Raw("SELECT phone FROM users WHERE phone=?", mobilenumber).Scan(&checkOtp)
	return checkOtp.Phone == mobilenumber

}

func ValidateOtp(c *gin.Context) {

	sotp := c.Query("otps")
	var userNum string

	initializers.DB.Raw("SELECT mobile FROM otps WHERE otp=?", sotp).Scan(&userNum)

	var user models.User
	//database.DB.Where("users.phone_number = ?", userNum).Find(&user)
	initializers.DB.First(&user, "phone = ?", userNum)

	// Look up request user
	var otp models.Otp
	initializers.DB.First(&otp, "otp = ?", sotp)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid OTP",
		})

		return
	}


	tokenstring, err := auth.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}
	// Sent it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAuth",tokenstring,3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "ok",
		"data":    tokenstring,
	})

	initializers.DB.Raw("DELETE FROM otps WHERE mobile=?", userNum).Scan(&otp)

}
