package models

// Model representing a User Entity
type UserModel struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
