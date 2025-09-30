package main

import (
	"fmt"
	"os"
)

// DB_PASSWORD=secret_pass go run env.go
func main() {
	dbPWD := os.Getenv("DB_PASSWORD")
	fmt.Println("password:", dbPWD)
}
