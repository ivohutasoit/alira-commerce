package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/ivohutasoit/alira"
	"github.com/ivohutasoit/alira-commerce/route"
	"github.com/ivohutasoit/alira/database/commerce"
	"github.com/joho/godotenv"
)

func init() {
	alira.GetConnection().Debug().AutoMigrate(&commerce.Customer{},
		&commerce.CustomerUser{},
		&commerce.Store{},
		&commerce.Product{},
		&commerce.StoreProduct{},
		&commerce.StoreProductPrice{})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("eror loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("$ORT must be set")
	}
	router := gin.New()
	router.Use(gin.Logger())

	store := cookie.NewStore([]byte(os.Getenv("SECRET_KEY")))
	router.Use(sessions.Sessions("ALIRASESSION", store))
	router.LoadHTMLGlob("views/*/*.tmpl.html")
	router.Static("/static", "static")

	web := &route.WebRoute{}
	web.Initialize(router)

	router.Run(":" + port)
}
