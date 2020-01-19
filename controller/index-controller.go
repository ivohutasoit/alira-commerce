package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexController struct{}

func (ctrl *IndexController) Get(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl.html", nil)
}

func (ctrl *IndexController) Post(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.tmpl.html", gin.H{
		"userid": "1",
	})
}
