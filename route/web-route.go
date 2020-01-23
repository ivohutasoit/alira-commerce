package route

import (
	"github.com/ivohutasoit/alira-commerce/controller"

	"github.com/gin-gonic/gin"
)

type WebRoute struct{}

func (route *WebRoute) Initialize(r *gin.Engine) {
	web := r.Group("")
	{
		index := &controller.IndexController{}
		webIndex := web.Group("/")
		webIndex.GET("", index.Get)
		webIndex.POST("", index.Post)

		customer := &controller.CustomerController{}
		customerWeb := web.Group("/customer")
		{
			customerWeb.GET("", customer.SearchHandler)
		}
	}
}
