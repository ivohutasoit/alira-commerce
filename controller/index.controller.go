package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira"
)

type Index struct{}

func (ctrl *Index) Get(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.tmpl.html", alira.ViewData)
}
