package controller

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira-commerce/model"
	"github.com/ivohutasoit/alira-commerce/service"
	"github.com/ivohutasoit/alira/database/commerce"
)

type Store struct{}

// CreateHandler godoc
// @Summary Create new store
// @Description Create new store with
// @Tags Store
// @Accept json
// @Produce json
// @Param store body model.Store true "Store info"
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /store [post]
func (ctrl *Store) CreateHandler(c *gin.Context) {
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))

	var req model.Store
	if api {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"error":  err.Error(),
			})
			return
		}
	}
	ss := &service.Store{}
	data, err := ss.Create(c.GetString("customer_id"), req.Code, req.Name, req.Address)
	if err != nil {
		if api {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": http.StatusText(http.StatusBadRequest),
				"error":  err.Error(),
			})
			return
		}
	}

	store := data["store"].(*commerce.Store)
	if api {
		c.JSON(http.StatusCreated, gin.H{
			"code":    http.StatusCreated,
			"status":  http.StatusText(http.StatusCreated),
			"message": data["message"].(string),
			"data": map[string]interface{}{
				"id":   store.Model.ID,
				"code": store.Code,
			},
		})
		return
	}
}

// DetailHandler godoc
// @Summary Store information
// @Description Detail of store based on id provided
// @Tags Store
// @Accept json
// @Produce json
// @Param id path string true "store id"
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /store/{id} [get]
func (ctrl *Store) DetailHandler(c *gin.Context) {

}

// ListHandler godoc
// @Summary List of Store
// @Description List of store based on customer
// @Tags Store
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /store [get]
func (ctrl *Store) ListHandler(c *gin.Context) {

}

// SearchHandler godoc
// @Summary Search Store
// @Description Find store based on parameter
// @Tags Store
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /store/search [post]
func (ctrl *Store) SearchHandler(c *gin.Context) {

}
