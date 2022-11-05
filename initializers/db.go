package initializers

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	Getenv()
}
func ConnecttoDb() {
	var err error
	dsn := os.Getenv("DB_DETAILS")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("error connecting to database")

	}
}
