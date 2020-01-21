package model

import (
	"time"

	alira "github.com/ivohutasoit/alira/model"
)

type Customer struct {
	alira.BaseModel
	Class     string `form:"-" json:"-" bson:"-" gorm:"default:SHOPOWNER"` // DISTRIBUTOR
	CreatedBy string `form:"-" json:"-" bson:"-"`
	UpdatedBy string `form:"-" json:"-" bson:"-"`
	Code      string `form:"code" json:"code" bson:"code" binding:"required"`
	Email     string
	Telephone string
	Fax       string
	Whatsapp  string
}

type CustomerDevice struct {
	alira.BaseModel
	CustomerID string
	DeviceName string
	IMEI       string
}

type Store struct {
	alira.BaseModel
	CustomerID string  `json:"customer_id" bson:"customer_id"`
	Class      string  `json:"class" bson:"class" gorm:"default:TRADITIONAL"`
	Address    string  `json:"address" bson:"address"`
	City       string  `json:"city" bson:"city"`
	State      string  `json:"state" bson:"state"`
	Country    string  `json:"country" bson:"country"`
	PostalCode string  `json:"postal_code" bson:"postal_code"`
	Telephone  string  `json:"telephone" bson:"telephone" gorm:"default:null"`
	Mobile     string  `json:"mobile" bson:"mobile" gorm:"default:null"`
	Website    string  `json:"website" bson:"website" gorm:"default:null"`
	Geocode    string  `json:"geocode" bson:"geocode" gorm:"default:null"`
	Longitude  float64 `json:"longitude" bson:"longitude" gorm:"default:null"`
	Latitude   float64 `json:"latitude" bson:"latitude" gorm:"default:null"`
}

type StoreUser struct {
	alira.BaseModel
	UserID string
	Role   string
	Pin    string
}

type StoreDevice struct {
	alira.BaseModel
	DeviceID string
	StoreID  string
}

type StoreProduct struct {
	alira.BaseModel
	StoreID   string
	ProductID string
	Image     string
}

type StoreProductInventory struct {
	alira.BaseModel
	StoreProductID string
	RackNo         string
	Available      int64
	Opname         int64
	Return         int64
	Sold           int64
}

type StoreProductPrice struct {
	alira.BaseModel
	StoreProductID string
	Unit           string
	BuyPrice       float64
	SellPrice      float64
	NotBefore      time.Time
	NotAfter       time.Time
}
