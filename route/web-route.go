package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira-commerce/controller"
	"github.com/ivohutasoit/alira-commerce/middleware"
)

type WebRoute struct{}

func (route *WebRoute) Initialize(r *gin.Engine) {
	authMiddleware := &middleware.Auth{}
	web := r.Group("")
	web.Use(authMiddleware.SessionRequired())
	{
		index := &controller.Index{}
		webIndex := web.Group("/")
		webIndex.GET("", index.Get)

		customer := &controller.Customer{}
		customerWeb := web.Group("/customer")
		{
			customerWeb.GET("", customer.SearchHandler)
			customerWeb.GET("/action", customer.ActionHandler)
			customerWeb.POST("/action", customer.ActionHandler)
		}
	}
}
