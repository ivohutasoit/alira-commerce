package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira-commerce/controller"
	_ "github.com/ivohutasoit/alira-commerce/docs"
	"github.com/ivohutasoit/alira-commerce/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Web struct{}

func (route *Web) Initialize(r *gin.Engine) {
	authMiddleware := &middleware.Auth{}
	callbackMiddleware := &middleware.Callback{}
	web := r.Group("")
	web.Use(callbackMiddleware.WebLoginToken())
	web.Use(authMiddleware.SessionRequired())
	{
		index := &controller.Index{}
		webIndex := web.Group("/")
		webIndex.GET("", index.Get)
		webIndex.GET("/developer/*any", ginSwagger.CustomWrapHandler(&ginSwagger.Config{
			URL: "doc.json",
		}, swaggerFiles.Handler))

		auth := &controller.Auth{}
		authWeb := web.Group("/auth")
		{
			authWeb.GET("/logout", auth.LogoutHandler)
		}

		customer := &controller.Customer{}
		customerWeb := web.Group("/customer")
		{
			customerWeb.GET("", customer.SearchHandler)
			customerWeb.GET("/action", customer.ActionHandler)
			customerWeb.POST("/action", customer.ActionHandler)
		}
	}
}
