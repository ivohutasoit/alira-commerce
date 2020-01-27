package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira/model/domain"
)

type IndexController struct{}

func (ctrl *IndexController) Get(c *gin.Context) {
	domain.Page["page"] = ""
	c.HTML(http.StatusOK, "dashboard.tmpl.html", domain.Page)
}
