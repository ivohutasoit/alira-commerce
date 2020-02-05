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

type StoreOrder struct {
	RecipeNo    string         `form:"recipe_no" json:"recipe_no" xml:"recipe_no"`
	Teller      string         `form:"teller" json:"teller" xml:"teller" binding:"required"`
	Store       string         `form:"store" json:"store" xml:"store" binding:"required"`
	Date        string         `form:"date" json:"date" xml:"date" binding:"required"`
	Time        string         `form:"time" json:"time" xml:"time" binding:"required"`
	Currency    string         `form:"currency" json:"currency" xml:"currency" binding:"required"`
	PaymentMode string         `form:"payment_mode" json:"payment_mode" xml:"payment_mode"`
	PaymentNo   string         `form:"payment_no" json:"payment_no" xml:"payment_no"`
	Subtotal    float64        `form:"subtotal" json:"subtotal" xml:"subtotal"`
	TaxPercent  float64        `form:"tax_percent" json:"tax_percent" xml:"tax_percent"`
	TaxAmount   float64        `form:"tax_amount" json:"tax_amount" xml:"tax_amount"`
	Total       float64        `form:"total" json:"total" xml:"total" binding:"required"`
	Rounding    float64        `form:"rounding" json:"rounding" xml:"rounding" binding:"rounding"`
	Status      string         `form:"status" json:"status" xml:"status" binding:"required"`
	Products    []OrderProduct `form:"products" json:"products" xml:"products" binding:"required"`
}

type OrderProduct struct {
	ID       string  `form:"id" json:"id" xml:"id" binding:"required"`
	Code     string  `form:"code" json:"code" xml:"code"`
	Name     string  `form:"name" json:"name" xml:"name"`
	Quantity int64   `form:"qty" json:"qty" xml:"qty" binding:"required"`
	Currency string  `form:"currency" json:"currency" xml:"currency" binding:"required"`
	Price    float64 `form:"price" json:"price" xml:"price" binding:"required"`
}

type SearchProduct struct {
	Keyword     string         `form:"keyword" json:"keyword" xml:"keyword" binding:"required"`
	Page        int64          `form:"page" json:"page" xml:"page"`
	Limit       int64          `form:"limit" json:"limit" xml:"limit"`
	TotalPage   int64          `form:"total_page" json:"total_page" xml:"total_page" binding:"-"`
	TotalRecord int64          `form:"total_record" json:"total_record" xml:"total_record" binding:"-"`
	Products    []StoreProduct `form:"products" json:"products" xml:"products" `
}
