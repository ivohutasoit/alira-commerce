package controller

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ivohutasoit/alira/util"

	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira"
)

type Auth struct{}

type Login struct {
	UserID string `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
}

// LoginHandler godoc
// @Summary Log in user
// @Description Handler user authentication
// @Accept json
// @Produce json
// @Param login body Login true "Login request"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /auth/login [post]
func (ctrl *Auth) LoginHandler(c *gin.Context) {
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))

	if c.Request.Method == http.MethodGet {
		url, err := util.GenerateUrl(c.Request.TLS, c.Request.Host, "/", true)
		if err != nil {
			fmt.Println(err)
		}
		redirect := fmt.Sprintf("%s?redirect=%s", os.Getenv("URL_LOGIN"), url)
		c.Redirect(http.StatusMovedPermanently, redirect)
		return
	}

	var req Login
	if api {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &alira.Response{
				Code:   http.StatusBadRequest,
				Status: http.StatusText(http.StatusBadRequest),
				Error:  err.Error(),
			})
			return
		}
	}
	if api {
		c.JSON(http.StatusOK, &alira.Response{
			Code:    http.StatusOK,
			Status:  http.StatusText(http.StatusBadRequest),
			Message: "Login successful",
		})
		return
	}
}

// LogoutHandler godoc
// @Summary Log out authenticated user
// @Description Handler log out authenticated user
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /auth/logout [post]
func (ctlr *Auth) LogoutHandler(c *gin.Context) {
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))

	if api {
		c.JSON(http.StatusOK, &alira.Response{})
	}
}
