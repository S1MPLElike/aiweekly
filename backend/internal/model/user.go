package model

import "time"

type User struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	Username     string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Phone        string `gorm:"type:varchar(11);uniqueIndex;not null" json:"phone"`

	// TokenVersion is used to invalidate previously issued refresh tokens.
	TokenVersion uint64 `gorm:"not null;default:0" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
