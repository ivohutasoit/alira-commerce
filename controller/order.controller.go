package controller

import "github.com/gin-gonic/gin"

type Order struct{}

// CreateHandler godoc
// @Summary Create order
// @Description Create order by store
// @Tags Sales Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /order [post]
func (ctrl *Order) CreateHandler(c *gin.Context) {

}

// DetailHandler godoc
// @Summary Order detail
// @Description Detail of order based on id or reference number provided
// @Tags Sales Order
// @Accept json
// @Produce json
// @Param id path string true "order id or reference number"
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /order/{id} [get]
func (ctrl *Order) DetailHandler(c *gin.Context) {

}

// SearchHandler godoc
// @Summary Find orders
// @Description Find order based on customer or store
// @Tags Sales Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /order/search [post]
func (ctrl *Order) SearchHandler(c *gin.Context) {

}
