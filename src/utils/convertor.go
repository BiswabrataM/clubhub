package utils

import (
	"fmt"
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConvertStrToInt(stringVal string) int {
	fmt.Println("--------------->>>>", stringVal)

	intVal, err := strconv.Atoi(stringVal)
	fmt.Println("--------------->>>>", err, intVal)
	if err != nil {
		return 0
	}
	return intVal
}
