package model

type StoreProduct struct {
	ID        string  `form:"id" json:"id" xml:"id"`
	Code      string  `form:"code" json:"code" xml:"code" binding:"required,min=3"`
	Store     string  `form:"store" json:"store" xml:"store" binding:"required"`
	Category  string  `form:"category" json:"category" xml:"category" binding:"required,min=3"`
	Barcode   string  `form:"barcode" json:"barcode" xml:"barcode" binding:"required"`
	Name      string  `form:"name" json:"name" xml:"name" binding:"required,min=3,max=20"`
	Currency  string  `form:"currency" json:"currency" xml:"currency" binding:"required,len=3"`
	Unit      string  `form:"unit" json:"unit" xml:"unit" binding:"required,gt=-1"`
	BuyPrice  float64 `form:"buy_price" json:"buy_price" xml:"buy_price" binding:"gt=-1"`
	SellPrice float64 `form:"sell_price" json:"sell_price" xml:"sell_price" binding:"required,gt=-1"`
	Quantity  int64   `form:"qty" json:"qty" xml:"qty" binding:"gt=-1"`
}
