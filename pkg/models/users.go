package models

type User struct {
	ID       string
	Username string `json:"username"`
	Password string `json:"-"`
}
