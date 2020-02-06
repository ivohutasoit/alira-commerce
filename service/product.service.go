package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ivohutasoit/alira"
	"github.com/ivohutasoit/alira-commerce/model"
	"github.com/ivohutasoit/alira/database/commerce"
)

type Product struct{}

func (s *Product) Create(args ...interface{}) (map[interface{}]interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("not enough parameter")
	}
	userid, ok := args[0].(string)
	if !ok {
		return nil, errors.New("plain text is not type string")
	}
	req, ok := args[1].(*model.StoreProduct)
	if !ok {
		return nil, errors.New("plain object is not type StoreProduct")
	}
	custUser := &commerce.CustomerUser{}
	alira.GetConnection().Where("user_id = ? AND active = ?", userid, true).First(&custUser)
	if custUser.Model.ID == "" {
		return nil, errors.New("invalid user")
	}

	store := &commerce.Store{}
	alira.GetConnection().Where("id = ? OR code = ?",
		req.Store, strings.ToUpper(req.Store)).First(&store)
	if store.Model.ID == "" {
		return nil, errors.New("invalid store")
	}

	if store.CustomerID != custUser.CustomerID {
		return nil, errors.New("permission denied")
	} else if custUser.Role != "OWNER" {
		storeUser := &commerce.StoreUser{}
		alira.GetConnection().Where("store_id = ? AND customer_user_id = ?",
			store.Model.ID, custUser.Model.ID).First(&storeUser)
		if storeUser.Model.ID == "" {
			return nil, errors.New("permission denied")
		}
	}

	product := &commerce.Product{}
	alira.GetConnection().Where("barcode = ?", req.Barcode).First(&product)
	if product.Model.ID == "" {
		product = &commerce.Product{
			Barcode: req.Barcode,
			Name:    strings.Title(strings.TrimSpace(req.Name)),
		}
		alira.GetConnection().Create(product)
	}

	category := &commerce.StoreProductCategory{}
	alira.GetConnection().Where("store_id = ? AND (code = ? OR name = ?)",
		store.Model.ID,
		strings.ToUpper(strings.TrimSpace(req.Category)),
		strings.Title(strings.TrimSpace(req.Category))).First(&category)
	if category.Model.ID == "" {
		category = &commerce.StoreProductCategory{
			StoreID: store.Model.ID,
			Code:    strings.ToUpper(strings.TrimSpace(req.Category)),
			Name:    strings.Title(strings.TrimSpace(req.Category)),
		}
		alira.GetConnection().Create(&category)
	}

	storeProduct := &commerce.StoreProduct{
		StoreID:    store.Model.ID,
		ProductID:  product.Model.ID,
		CategoryID: category.Model.ID,
		Name:       strings.Title(strings.TrimSpace(req.Name)),
		Available:  req.Quantity,
	}
	alira.GetConnection().Create(storeProduct)

	price := &commerce.StoreProductPrice{
		StoreProductID: storeProduct.Model.ID,
		Quantity:       1,
		Unit:           req.Unit,
		Currency:       req.Currency,
		BuyPrice:       req.BuyPrice,
		SellPrice:      req.SellPrice,
		NotBefore:      time.Now(),
	}
	alira.GetConnection().Create(price)

	req.ID = storeProduct.Model.ID

	return map[interface{}]interface{}{
		"message": "Product has been added",
		"product": req,
	}, nil
}

func (s *Product) Get(args ...interface{}) (map[interface{}]interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("not enough parameter")
	}
	id, ok := args[0].(string)
	if !ok {
		return nil, errors.New("plain text is not type string")
	}
	storeProduct := &commerce.StoreProduct{}
	alira.GetConnection().Where("id = ?", id).First(&storeProduct)
	if storeProduct.Model.ID == "" {
		return nil, errors.New("invalid product")
	}
	product := &commerce.Product{}
	alira.GetConnection().Where("id = ?", storeProduct.ProductID).First(&product)
	if storeProduct.Model.ID == "" {
		return nil, errors.New("unknown product")
	}
	detail := &messaging.StoreProduct{
		ID:      storeProduct.Model.ID,
		Code:    storeProduct.Code,
		Barcode: product.Barcode,
		Name:    storeProduct.Name,
	}

	store := &commerce.Store{}
	alira.GetConnection().Where("id = ?", storeProduct.StoreID).First(&store)
	category := &commerce.StoreProductCategory{}
	alira.GetConnection().Where("id = ?", storeProduct.CategoryID).First(&category)
	price := &commerce.StoreProductPrice{}
	alira.GetConnection().Where("store_product_id = ?", storeProduct.ID).First(&price)

	if store.ID != "" {
		detail.Store = fmt.Sprintf("%s - %s", store.Code, store.Name)
	}

	if category.ID != "" {
		detail.Category = fmt.Sprintf("%s - %s", category.Code, category.Name)
	}

	if price.ID != "" {
		detail.Currency = price.Currency
		detail.Unit = price.Unit
		detail.BuyPrice = price.BuyPrice
		detail.SellPrice = price.SellPrice
		detail.Quantity = price.Quantity
	}

	return map[interface{}]interface{}{
		"product": detail,
	}, nil
}
