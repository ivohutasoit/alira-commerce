package controller

import "github.com/gin-gonic/gin"

type Product struct{}

// CreateHandler godoc
// @Summary Create new product
// @Description Create new product based on customer and or store
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /product [post]
func (ctrl *Product) CreateHandler(c *gin.Context) {

}

// DetailHandler godoc
// @Summary Product information
// @Description Detail of product based on id provided
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
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /order/search [post]
func (ctrl *Product) SearchHandler(c *gin.Context) {

}
