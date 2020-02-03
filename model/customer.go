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

type Store struct {
	Owner   string `json:"owner" bson:"owner"`
	Code    string `json:"code" bson:"code" binding:"required,min=6"`
	Name    string `json:"name" bson:"name" binding:"required,min=3"`
	Address string `json:"address" bson:"address" binding:"required,min=10"`
}
