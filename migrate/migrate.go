package main

import (
	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
)

func init() {
	initializers.ConnecttoDb()
	initializers.Getenv()
}
func main() {
	initializers.DB.AutoMigrate(
		&models.User{},

		&models.Address{},
		&models.Product{},
		&models.Cart{},
		&models.Brand{},
		&models.Otp{},
		&models.Orders{},
		&models.Orderd_Items{},
		&models.Cartsinfo{},
		&models.ShoeSize{},
		&models.Catogory{},
		&models.RazorPay{},
		&models.Coupon{},
		&models.Applied_Coupons{},
		&models.WishList{},
	)
}
