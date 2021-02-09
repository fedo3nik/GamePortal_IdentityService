package dtohttp

type UserDTO struct {
	ID           string `json:"id"`
	WarningCount uint   `json:"warningCount"`
	Nickname     string `json:"nickname"`
	Password     string `json:"password"`
	Email        string `json:"email"`
}
