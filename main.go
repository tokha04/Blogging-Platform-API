package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tokha04/blogging-platform-api/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("could not load .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.Default()
	routes.Routes(router)

	log.Fatal(router.Run(":" + port))
}
