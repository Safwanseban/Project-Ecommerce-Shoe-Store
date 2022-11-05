package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

func GenerateCoupon(c *gin.Context) {
	// Cname := c.PostForm("couponName")
	cCode := c.PostForm("couponcode")
	Cdiscount := c.PostForm("discount")
	discount, _ := strconv.Atoi(Cdiscount)
	cQuanityt := c.PostForm("quantity")
	quantity, _ := strconv.Atoi(cQuanityt)
	cValidity := c.PostForm("validity")
	validity, _ := strconv.Atoi(cValidity)

	timew := time.Now()
	fmt.Println(timew)
	expirationTime := time.Now().AddDate(0, 0, validity)
	expirationTime.Unix()
	coupons := models.Coupon{

		Coupon_Code: cCode,
		Discount:    uint(discount),
		Quantity:    uint(quantity),
		Validity:    expirationTime.Unix(),
	}
	initializers.DB.Create(&coupons)

	c.JSON(200, gin.H{
		"coupon":      coupons.Coupon_Code,
		"coupon-code": coupons.Coupon_Code,
		"msg":         "coupon added",
	})

	// if expirationTime.Unix()<time.Now().Local().Unix(){
	// 	fmt.Println("epired")
	// }else{
	// 	fmt.Println("not expired")
	// }
	// fmt.Println(expirationTime)

}

// func CouponGenerator(name string, discount int) string {

// }
