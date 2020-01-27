package controller

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"

	"github.com/ivohutasoit/alira/util"

	"github.com/ivohutasoit/alira-commerce/service"

	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira/model/domain"
)

type CustomerController struct{}

func (ctrl *CustomerController) DetailHandler(c *gin.Context) {
	domain.Page["view"] = "customer"
	code := c.Query("code")
	if c.Request.Method == http.MethodGet {
		domain.Page["message"] = nil
		domain.Page["error"] = nil
		domain.Page["code"] = code

		c.HTML(http.StatusOK, "customer-detail.tmpl.html", domain.Page)
		return
	}
}

func (ctrl *CustomerController) SearchHandler(c *gin.Context) {
	domain.Page["view"] = "customer"
	page := c.DefaultQuery("page", "1")
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))
	session := sessions.Default(c)

	customerService := &service.CustomerService{}
	if c.Request.Method == http.MethodGet {
		domain.Page["message"] = nil
		data, err := customerService.Search(session.Get("access_token"), page)
		if err != nil {
			domain.Page["data"] = nil
			domain.Page["error"] = err.Error()
		} else {
			domain.Page["error"] = nil
			paginator := data["paginator"].(*util.Paginator)
			domain.Page["data"] = paginator
		}
		c.HTML(http.StatusOK, "customer-index.tmpl.html", domain.Page)
		return
	}

	type Request struct {
		Code     string `form:"code" json:"code" xml:"code"`
		Fullname string `form:"fullname" json:"fullname" xml:"fullname"`
		Status   string `form:"status" json:"status" xml:"status"`
	}

	var req Request
	if api {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"error":  err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&req); err != nil {
			domain.Page["error"] = err.Error()
			c.HTML(http.StatusBadRequest, "customer-index.tmpl.html", domain.Page)
			return
		}
	}
}

func (ctrl *CustomerController) ActionHandler(c *gin.Context) {
	name := c.DefaultQuery("name", "create")
	switch name {
	case "detail":
		ctrl.DetailHandler(c)
		return
	default:
		ctrl.CreateHandler(c)
		return
	}
}

func (ctrl *CustomerController) CreateHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		domain.Page["view"] = "customer"
		c.HTML(http.StatusOK, "customer-create.tmpl.html", domain.Page)
		return
	}

	type Request struct {
		Code      string `form:"code" json:"code" xml:"code" binding:"required"`
		Username  string `form:"username" json:"username" xml:"username" binding:"required"`
		Email     string `form:"email" json:"email" xml:"email" binding:"required"`
		Mobile    string `form:"mobile" json:"mobile" xml:"mobile" binding:"required"`
		FirstName string `form:"first_name" json:"first_name" xml:"first_name" binding:"required"`
		LastName  string `form:"last_name" json:"last_name" xml:"last_name" binding:"required"`
		Payment   int    `form:"payment" json:"payment" xml:"payment"`
	}

	var req Request
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))
	if api {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"error":  err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&req); err != nil {
			domain.Page["view"] = "customer"
			domain.Page["message"] = nil
			domain.Page["error"] = err.Error()
			c.HTML(http.StatusBadRequest, "customer-create.tmpl.html", domain.Page)
			return
		}
	}

	customerService := &service.CustomerService{}
	session := sessions.Default(c)
	data, err := customerService.Create(req.Code, req.Username, req.Email, req.Mobile,
		req.FirstName, req.LastName, req.Payment, session.Get("access_token"))
	if err != nil {
		if api {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"error":  err.Error(),
			})
			return
		} else {
			domain.Page["view"] = "customer"
			domain.Page["message"] = nil
			domain.Page["error"] = err.Error()
			c.HTML(http.StatusBadRequest, "customer-create.tmpl.html", domain.Page)
			return
		}
	}

	domain.Page["view"] = "customer"
	domain.Page["message"] = data["message"].(string)
	domain.Page["error"] = nil
	c.HTML(http.StatusCreated, "customer-create.tmpl.html", domain.Page)
}
