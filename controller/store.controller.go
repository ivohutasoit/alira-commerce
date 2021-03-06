package controller

import "github.com/gin-gonic/gin"

type Store struct{}

// CreateHandler godoc
// @Summary Create new store
// @Description Create new store with using customer id
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /store [post]
func (ctrl *Store) CreateHandler(c *gin.Context) {

}

// DetailHandler godoc
// @Summary Store information
// @Description Detail of store based on id provided
// @Accept json
// @Produce json
// @Param id path string true "store id"
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /store/{id} [get]
func (ctrl *Store) DetailHandler(c *gin.Context) {

}

// SearchHandler godoc
// @Summary Search Store
// @Description Find store based on parameter
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [accessing token]"
// @Success 200 {string} string "{"code": 200, "status": "OK", "message": "Success", "data": "data"}"
// @Failure 400 {string} string "{"code": 400, "status": "Bad request", "error": "Error"}"
// @Router /store/search [post]
func (ctrl *Store) SearchHandler(c *gin.Context) {

}
