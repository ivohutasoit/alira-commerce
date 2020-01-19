package main

import (
	"fmt"
	"os"

	"github.com/ivohutasoit/alira-commerce/route"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("$PORT must be set")
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("views/*/*.tmpl.html")
	router.Static("/static", "static")

	web := &route.WebRoute{}
	web.Initialize(router)

	router.Run(":" + port)
}
