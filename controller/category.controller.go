package controller

import "github.com/gin-gonic/gin"

type Category struct{}

// ListHandler godoc
// @Summary List of Store
// @Description List of product category based on customer
// @Tags Product
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /category [get]
func (ctrl *Category) ListHandler(c *gin.Context) {

}
