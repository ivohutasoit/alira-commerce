package controller

import (
	"net/http"
	"os"
	"strings"

	"github.com/ivohutasoit/alira-commerce/model"
	"github.com/ivohutasoit/alira-commerce/service"

	"github.com/gin-gonic/gin"
)

type User struct{}

// PinHandler godoc
// @Summary Change user pin
// @Description Update authenticated user pin
// @Tags User
// @Accept json
// @Produce json
// @Param pin body model.Pin true "User pin"
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /user/pin [post]
func (ctrl *User) PinHandler(c *gin.Context) {
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))

	var req model.Pin
	if api {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"error":  err.Error(),
			})
		}
	}
	us := &service.User{}
	if api {
		tokens := strings.Split(c.Request.Header.Get("Authorization"), " ")
		data, err := us.ChangeUserPin(tokens[1], req.Code)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"status":  http.StatusText(http.StatusOK),
			"message": data["message"].(string),
		})
	}

}
