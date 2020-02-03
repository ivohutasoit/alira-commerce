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
		&commerce.StoreUser{},
		&commerce.Product{},
		&commerce.StoreProduct{},
		&commerce.StoreProductPrice{})
}

// @title Alira Commerce API
// @version alpha
// @description Documentation of Alira commerce provides capability to manage customer store, inventory and sales order
// @termsOfService https://commerce.alira.com/terms/

// @contact.name Alira Support
// @contact.url https://www.commerce.alira.com/support
// @contact.email hello@alira.com

// @license.name GNU General Public License v3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html

// @host aliracommerce.herokuapp.com
// @BasePath /api/alpha
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

	web := &route.Web{}
	web.Initialize(router)

	api := &route.Api{}
	api.Initialize(router)

	router.Run(":" + port)
}
