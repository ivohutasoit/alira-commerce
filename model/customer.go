package model

type CustomerRecord struct {
	ID       string
	Code     string
	Name     string
	Username string
	Email    string
	Mobile   string
	Status   string
	Payment  bool
}

type CustomerProfile struct {
	ID        string
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
