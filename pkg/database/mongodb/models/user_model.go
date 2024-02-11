package models

type UserModel struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessLevel string `json:"access_level"`
	Password    string `json:"-"`
}
