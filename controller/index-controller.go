package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexController struct{}

func (ctrl *IndexController) Get(c *gin.Context) {
	alira.ViewData["page"] = ""
	c.HTML(http.StatusOK, "dashboard.tmpl.html", alira.ViewData)
}
