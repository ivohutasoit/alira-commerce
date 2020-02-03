package messaging

type LoggedUser struct {
	UserID       string `json:"user_id"`
	UsePin       bool   `json:"use_pin"`
	Purpose      string `json:"purpose"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
