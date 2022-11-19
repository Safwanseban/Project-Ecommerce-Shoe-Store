package initializers

import (
	"fmt"
	"os"

	"github.com/Safwanseban/Project-Ecommerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	Getenv()
}
func ConnecttoDb() {
	var err error
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbname, password)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("error connecting to database")

	}
	DB.AutoMigrate(
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
