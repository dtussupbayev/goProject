package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID                int    `json:"id" db:"id"`
	Email             string `json:"email" db:"email"`
	Password          string `json:"password,omitempty" db:"-"`
	EncryptedPassword string `json:"-" db:"password"`
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}
