package controller

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira-commerce/model"
	"github.com/ivohutasoit/alira-commerce/service"
)

type Product struct{}

// CreateHandler godoc
// @Summary Create new product
// @Description Create new product based on customer and or store
// @Tags Product
// @Accept json
// @Produce json
// @Param store body model.StoreProduct true "Store product"
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /product [post]
func (ctrl *Product) CreateHandler(c *gin.Context) {
	api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))

	var req model.StoreProduct
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

	ps := &service.Product{}
	data, err := ps.Create(c.GetString("user_id"), req)
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

	if api {
		c.JSON(http.StatusCreated, gin.H{
			"code":    http.StatusCreated,
			"status":  http.StatusText(http.StatusCreated),
			"message": data["message"].(string),
			"data":    data["product"].(*model.StoreProduct),
		})
		return
	}
}

// ListHandler godoc
// @Summary List of Product
// @Description List of Product based on authenticated user
// @Tags Product
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /product [get]
func (ctrl *Product) ListHandler(c *gin.Context) {

}

// DetailHandler godoc
// @Summary Product information
// @Description Detail of product based on id provided
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "product id or barcode"
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /product/{id} [get]
func (ctrl *Product) DetailHandler(c *gin.Context) {

}

// SearchHandler godoc
// @Summary Find orders
// @Description Find order based on customer or store
// @Tags Product
// @Accept json
// @Produce json
// @Param search body model.SearchProduct true "Store product search"
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /product/search [post]
func (ctrl *Product) SearchHandler(c *gin.Context) {

}
