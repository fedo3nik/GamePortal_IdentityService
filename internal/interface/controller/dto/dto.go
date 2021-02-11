package controller

type UserDTO struct {
	ID           string `json:"id"`
	WarningCount uint   `json:"warningCount"`
	Nickname     string `json:"nickname"`
	Password     string `json:"password"`
	Email        string `json:"email"`
}

type SignUpResponse struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type SignInResponse struct {
	ID           string `json:"id"`
	Nickname     string `json:"nickname"`
	Email        string `json:"email"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AddWarnResponse struct {
	ID         string `json:"id"`
	SuccessAdd bool   `json:"success"`
}

type RemWarnResponse struct {
	ID         string `json:"id"`
	SuccessRem bool   `json:"success"`
}

type IsBannedResponse struct {
	ID       string `json:"id"`
	IsBanned bool   `json:"isBanned"`
}
