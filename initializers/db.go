package initializers

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnecttoDb() {
	var err error
	dsn := "host=localhost user=safwan password=Safwan@123 dbname=ecommerce port=5432 "
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("error connecting to database")

	}
}
