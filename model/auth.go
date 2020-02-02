package model

type Login struct {
	UserID string `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
}

type Token struct {
	UserID string `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
	Token  string `form:"token" json:"token" xml:"token" binding:"required"`
}

type Pin struct {
	Code string `form:"code" json:"code" xml:"code" binding:"required"`
}
