package model

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Nickname     string
	Email        string
	Phone        string
	Status       int32
	DeptID       int64
	CreatedAt    string
}
