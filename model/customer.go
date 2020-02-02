package model


type CustomerProfile struct {
	ID        string `json:"id" bson:"id"`
	UserID    string `json:"user_id" bson:"user_id"`
	Code      string
	FirstName string
	LastName  string
	Name      string `json:"name" bson:"name"`
	Username  string `json:"username" bson:"username"`
	Email     string `json:"email" bson:"email"`
	Mobile    string `json:"mobile" bson:"mobile"`
	Status    string
	Payment   bool
}

type CustomerStore struct {
	ID string
	Name string
	Address string
	Status string
}