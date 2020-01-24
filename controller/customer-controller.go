package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira/model/domain"
)

type CustomerController struct{}

func (ctrl *CustomerController) DetailHandler(c *gin.Context) {

}

func (ctrl *CustomerController) SearchHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		domain.Page["userid"] = "1"
		domain.Page["page"] = "customer"
		c.HTML(http.StatusOK, "customer-index.tmpl.html", domain.Page)
		return
	}
}

func (ctrl *CustomerController) AddHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		return
	}

}
