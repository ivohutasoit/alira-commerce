package route

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira-commerce/controller"
	"github.com/ivohutasoit/alira-commerce/middleware"
)

type Api struct{}

func (route *Api) Initialize(r *gin.Engine) {
	authMiddleware := &middleware.Auth{}
	api := r.Group(os.Getenv("URL_API"))
	api.Use(authMiddleware.TokenRequired())
	{
		auth := &controller.Auth{}
		authApi := api.Group("/auth")
		{
			authApi.POST("/login", auth.LoginHandler)
			authApi.POST("/token", auth.TokenHandler)
			authApi.POST("/logout", auth.LogoutHandler)
		}
		user := &controller.User{}
		userApi := api.Group("/user")
		{
			userApi.POST("/pin", user.PinHandler)
		}
	}
}
