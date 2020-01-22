package controller

import "github.com/gin-gonic/gin"

import "net/http"

type CustomerController struct{}

func (ctrl *CustomerController) SearchHandler(c *gin.Context) {

}

func (ctrl *CustomerController) AddHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		return
	}

}
