package controller

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira"
	"github.com/ivohutasoit/alira-commerce/messaging"
	"github.com/ivohutasoit/alira-commerce/model"
	"github.com/ivohutasoit/alira-commerce/service"
	"github.com/ivohutasoit/alira/util"
)

type Auth struct{}

// LoginHandler godoc
// @Summary Log in user
// @Description Handler user authentication
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body model.Login true "Login Request"
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

	var req model.Login
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
	as := &service.Auth{}
	data, err := as.AuthenticateUser(req.UserID)
	if err != nil {
		if api {
			c.AbortWithStatusJSON(http.StatusBadRequest, &alira.Response{
				Code:   http.StatusBadRequest,
				Status: http.StatusText(http.StatusBadRequest),
				Error:  err.Error(),
			})
			return
		}
	}
	loggedUser := data["user"].(*messaging.LoggedUser)
	if api {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"status":  http.StatusText(http.StatusOK),
			"message": data["message"].(string),
			"data":    loggedUser,
		})
		return
	}
}

// TokenHandler godoc
// @Summary Verify token
// @Description Authentication token verification handler
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token body model.Token true "Authentication Token"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /auth/token [post]
func (ctrl *Auth) TokenHandler(c *gin.Context) {
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))

	var req model.Token
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
	as := &service.Auth{}
	data, err := as.VerifyToken(req.UserID, req.Token)
	if err != nil {
		if api {
			c.AbortWithStatusJSON(http.StatusBadRequest, &alira.Response{
				Code:   http.StatusBadRequest,
				Status: http.StatusText(http.StatusBadRequest),
				Error:  err.Error(),
			})
			return
		}
	}
	loggedUser := data["user"].(*messaging.LoggedUser)
	if api {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"status":  http.StatusText(http.StatusOK),
			"message": data["message"].(string),
			"data":    loggedUser,
		})
		return
	}
}

// LogoutHandler godoc
// @Summary Log out authenticated user
// @Description Handler log out authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /auth/logout [post]
func (ctlr *Auth) LogoutHandler(c *gin.Context) {
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))
	as := &service.Auth{}
	var accessToken string
	session := sessions.Default(c)
	if api {
		tokens := strings.Split(c.Request.Header.Get("Authorization"), " ")
		accessToken = tokens[1]
	} else {
		accessToken = session.Get("access_token").(string)
	}

	data, err := as.RemoveSessionToken(accessToken)
	if err != nil {
		if api {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"error":  err.Error(),
			})
		}
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	if api {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"status":  http.StatusText(http.StatusOK),
			"message": data["message"].(string),
		})
		return
	}

	session.Clear()
	session.Save()

	c.Redirect(http.StatusMovedPermanently, "/")
}
