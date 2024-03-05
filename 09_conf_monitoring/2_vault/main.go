package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
)

// token=123 go run 2_vault/main.go
func main() {
	client, err := api.NewClient(&api.Config{
		Address: "http://127.0.0.1:8200",
	})

	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("token")
	client.SetToken(token)
	secretValues, err := client.Logical().Read("secret/data/pg_credentials")
	if err != nil {
		log.Fatal(err)
	}
	if secretValues == nil {
		log.Fatal("Secret not found")
	}
	fmt.Printf("data: %s\n", secretValues.Data["data"])
}
