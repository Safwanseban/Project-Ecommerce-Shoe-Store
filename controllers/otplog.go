package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Safwanseban/Project-Ecommerce/auth"
	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"

	"github.com/twilio/twilio-go"
	// openapi "github.com/twilio/twilio-go/rest/api/v2010"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

func init() {
	initializers.Getenv()
	initializers.ConnecttoDb()
}

var (
	accountSid string
	authToken  string
	fromPhone  string

	client *twilio.RestClient
)

func OtpLog(c *gin.Context) {
	accountSid = os.Getenv("ACCOUNT_SID")
	authToken = os.Getenv("AUTH_TOKEN")
	fromPhone = os.Getenv("FROM_PHONE")
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	Mob := c.Query("number")

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

	mobile := "+91" + Mob
	// Creatin 4 digit OTP

	params := &verify.CreateVerificationParams{}
	params.SetTo(mobile)
	params.SetChannel("sms")
	resp, err := client.VerifyV2.CreateVerification(fromPhone, params)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{
			"status":  false,
			"message": "error sending OTP",
		})
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
		c.JSON(200, gin.H{
			"status":  true,
			"message": "OTP Sent Succesfully",
		})
	}

}

// Checking if the number belongs to any user
func ChekNumber(str string) bool {

	mobilenumber := str
	var checkOtp models.User
	initializers.DB.Raw("SELECT phone FROM users WHERE phone=?", mobilenumber).Scan(&checkOtp)
	return checkOtp.Phone == mobilenumber

}
func CheckOtp(c *gin.Context) {
	accountSid = os.Getenv("ACCOUNT_SID")
	authToken = os.Getenv("AUTH_TOKEN")
	fromPhone = os.Getenv("FROM_PHONE")
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	Mob := c.Query("number")
	code := c.Query("otps")

	ChekNumber(Mob)
	var user models.User
	initializers.DB.First(&user, "phone = ?", Mob)

	mobile := "+91" + Mob
	fromPhone = os.Getenv("FROM_PHONE")
	fmt.Println(mobile)
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(mobile)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(fromPhone, params)

	if err != nil {
		fmt.Println(err.Error())
	} else if *resp.Status == "approved" {
		tokenstring,  err := auth.GenerateJWT(user.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create token",
			})

			return
		}
		// Sent it back
		fmt.Println(tokenstring)
		token:=tokenstring["access_token"]	
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("UserAuth", token, 3600*24*30, "", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"message":   "ok",
			"data":      tokenstring,

		})
	} else {
		
		c.JSON(404, gin.H{
			"msg": "otp is invalid",
		})
	}
}

// func ValidateOtp(c *gin.Context) {

// 	sotp := c.Query("otps")
// 	var userNum string

// 	initializers.DB.Raw("SELECT mobile FROM otps WHERE otp=?", sotp).Scan(&userNum)

// 	var user models.User
// 	//database.DB.Where("users.phone_number = ?", userNum).Find(&user)
// 	initializers.DB.First(&user, "phone = ?", userNum)

// 	// Look up request user
// 	var otp models.Otp
// 	initializers.DB.First(&otp, "otp = ?", sotp)

// 	if user.ID == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  false,
// 			"message": "Invalid OTP",
// 		})

// 		return
// 	}

// 	tokenstring, ex, err := auth.GenerateJWT(user.Email)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to create token",
// 		})

// 		return
// 	}
// 	// Sent it back
// 	c.SetSameSite(http.SameSiteLaxMode)
// 	c.SetCookie("UserAuth", tokenstring, 3600*24*30, "", "", false, true)

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":    true,
// 		"message":   "ok",
// 		"data":      tokenstring,
// 		"expiresAt": ex,
// 	})

// 	initializers.DB.Raw("DELETE FROM otps WHERE mobile=?", userNum).Scan(&otp)

// }
