package controller

import "github.com/gin-gonic/gin"

type Store struct{}

// ListHandler godoc
// @Summary Search Store
// @Description Find store based on parameter
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
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /store [post]
func (ctrl *Store) SearchHandler(c *gin.Context) {

}
