package models

import "time"

type LoginRequest struct {
	Email    string `json:"Email" form:"Email" binding:"required"`
	Password string `json:"Password" form:"Password" binding:"required"`
}
type AuthorizedUserInfo struct {
	Id int `json:"id"`
}
type Tokens struct {
	AccessToken  string
	RefreshToken string
}
type Session struct {
	RefreshToken string    `json:"refreshToken" db:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expiresAt"`
}
