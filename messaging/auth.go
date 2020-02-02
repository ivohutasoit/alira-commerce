package messaging

type LoggedUser struct {
	UserID       string `json:"user_id"`
	Purpose      string `json:"purpose"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
