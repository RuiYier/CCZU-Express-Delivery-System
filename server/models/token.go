package models

import "time"

type UserToken struct {
	TokenId      int64     `gorm:"primaryKey;autoIncrement" json:"token_id"`
	UserId       int64     `gorm:"not null;index" json:"user_id"`
	AccessToken  string    `gorm:"not null" json:"access_token"`
	RefreshToken string    `gorm:"not null" json:"refresh_token"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}
