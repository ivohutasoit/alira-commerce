package service

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
	default:
		
	}

	return nil, nil
}
