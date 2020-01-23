package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerController struct{}

func (ctrl *CustomerController) DetailHandler(c *gin.Context) {

}

func (ctrl *CustomerController) SearchHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "customer-index.tmpl.html", gin.H{
			"userid": "1",
			"page":   "customer",
		})
		return
	}
}

func (ctrl *CustomerController) AddHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		return
	}

}
