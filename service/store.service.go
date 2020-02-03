package service

import (
	"errors"
	"strings"

	"github.com/ivohutasoit/alira"
	"github.com/ivohutasoit/alira/database/commerce"
)

type Store struct{}

func (s *Store) Create(args ...interface{}) (map[interface{}]interface{}, error) {
	if len(args) < 4 {
		return nil, errors.New("not enough paramater")
	}
	var owner, code, name, address string
	for i, v := range args {
		param, ok := v.(string)
		if !ok {
			return nil, errors.New("plain text is not string type")
		}
		switch i {
		case 0:
			owner = strings.TrimSpace(param)
		case 1:
			code = strings.TrimSpace(strings.ToUpper(param))
		case 2:
			name = strings.TrimSpace(strings.Title(param))
		case 3:
			address = strings.TrimSpace(strings.Title(param))
		}
	}
	customer := &commerce.Customer{}
	alira.GetConnection().Where("id = ? AND status = 'ACTIVE'",
		owner).First(&customer)
	if customer.Model.ID == "" {
		return nil, errors.New("invalid owner")
	}

	store := &commerce.Store{
		CustomerID: customer.Model.ID,
		Code:       code,
		Name:       name,
		Address:    address,
	}
	alira.GetConnection().Create(&store)

	return map[interface{}]interface{}{
		"message": "Store has been created successful",
		"store":   store,
	}, nil
}

func (s *Store) Search(args ...interface{}) (map[interface{}]interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("not enough paramater")
	}
	findBy, ok := args[0].(string)
	if !ok {
		return nil, errors.New("plain text is not string type")
	}
	var stores []commerce.Store
	switch findBy {
	case "customer":
		alira.GetConnection().Where("customer_id = ?", args[1].(string)).Find(&stores)
	default:
		alira.GetConnection().Find(&stores)
	}

	return map[interface{}]interface{}{
		"stores": stores,
	}, nil
}
