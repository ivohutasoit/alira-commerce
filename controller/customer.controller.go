package controller

import (
	"net/http"
	"os"
	"strings"

	"github.com/ivohutasoit/alira-commerce/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira"
	"github.com/ivohutasoit/alira-commerce/service"
	"github.com/ivohutasoit/alira/util"
)

type Customer struct{}

func (ctrl *Customer) DetailHandler(c *gin.Context) {
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))
	var id, accessToken string

	if api {
		id = c.Param("id")
		tokens := strings.Split(c.Request.Header.Get("Authorization"), " ")
		accessToken = tokens[1]
	} else {
		alira.ViewData["view"] = "customer"
		id = c.Query("id")
		session := sessions.Default(c)
		accessToken = session.Get("access_token").(string)
	}
	cs := &service.Customer{}
	data, err := cs.Get(accessToken, id)
	if err != nil {
		if api {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"error":  err.Error(),
			})
		} else {
			alira.ViewData["error"] = err.Error()
			c.HTML(http.StatusOK, "customer-index.tmpl.html", alira.ViewData)
		}
		return
	}

	if api {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"status":  http.StatusText(http.StatusOK),
			"message": "Customer detail",
			"data": map[string]interface{}{
				"customer": data["customer"].(*model.CustomerProfile),
			},
		})
		return
	}
	alira.ViewData["customer"] = data["customer"].(*model.CustomerProfile)
	alira.ViewData["stores"] = data["stores"]
	c.HTML(http.StatusOK, "customer-detail.tmpl.html", alira.ViewData)
}

func (ctrl *Customer) SearchHandler(c *gin.Context) {
	alira.ViewData["view"] = "customer"
	page := c.DefaultQuery("page", "1")
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))
	session := sessions.Default(c)

	cs := &service.Customer{}
	if c.Request.Method == http.MethodGet {
		alira.ViewData["message"] = nil
		data, err := cs.Search(session.Get("access_token"), page)
		if err != nil {
			alira.ViewData["data"] = nil
			alira.ViewData["error"] = err.Error()
		} else {
			alira.ViewData["error"] = nil
			paginator := data["paginator"].(*util.Paginator)
			alira.ViewData["data"] = paginator
		}
		c.HTML(http.StatusOK, "customer-index.tmpl.html", alira.ViewData)
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
			alira.ViewData["error"] = err.Error()
			c.HTML(http.StatusBadRequest, "customer-index.tmpl.html", alira.ViewData)
			return
		}
	}
}

func (ctrl *Customer) ActionHandler(c *gin.Context) {
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

func (ctrl *Customer) CreateHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		alira.ViewData["view"] = "customer"
		c.HTML(http.StatusOK, "customer-create.tmpl.html", alira.ViewData)
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
			alira.ViewData["view"] = "customer"
			alira.ViewData["message"] = nil
			alira.ViewData["error"] = err.Error()
			c.HTML(http.StatusBadRequest, "customer-create.tmpl.html", alira.ViewData)
			return
		}
	}

	cs := &service.Customer{}
	session := sessions.Default(c)
	data, err := cs.Create(req.Code, req.Username, req.Email, req.Mobile,
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
			alira.ViewData["view"] = "customer"
			alira.ViewData["message"] = nil
			alira.ViewData["error"] = err.Error()
			c.HTML(http.StatusBadRequest, "customer-create.tmpl.html", alira.ViewData)
			return
		}
	}

	alira.ViewData["view"] = "customer"
	alira.ViewData["message"] = data["message"].(string)
	alira.ViewData["error"] = nil
	c.HTML(http.StatusCreated, "customer-create.tmpl.html", alira.ViewData)
}
