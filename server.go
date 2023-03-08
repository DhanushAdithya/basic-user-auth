package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func LoadEnv() error {
	prod := os.Getenv("PROD")

	if prod != "true" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	err := LoadEnv()
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	app.Use(cors.New())

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(":" + port)
	fmt.Println("Server started on port " + port)
}
