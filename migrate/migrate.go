package main

import (
	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
)


func init()  {
	initializers.ConnecttoDb()
}
func main()  {
	initializers.DB.AutoMigrate(
		&models.User{},
		&models.Cart{},
		&models.Address{},
		&models.Product{},
		&models.Brand{},
		&models.Otp{},
		&models.Orders{},
		&models.Orderd_Items{},
		models.Cartsinfo{},
	models.ShoeSize{},models.Catogory{})
}