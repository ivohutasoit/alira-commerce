package service

import (
	"errors"

	"github.com/ivohutasoit/alira"
	"github.com/ivohutasoit/alira/database/commerce"
)

type Store struct{}

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
