package initializers

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Getenv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading env file")

	}

}
