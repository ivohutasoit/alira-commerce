package controller

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Auth struct{}

func (ctrl *Auth) LoginHandler(c *gin.Context) {
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))

	if c.Request.Method == http.MethodGet {
		return
	}

	type Request struct {
	}

	if api {

	}
}
