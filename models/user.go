package models

import "time"

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex" json:"user_name"`
	Password  string    `json:"password"`
	Sex       string    `json:"sex"`
	BornDate  time.Time `json:"born_date"`
	Identify  string    `json:"ide"`
	AvatarURL string    `json:"avatar_url"`
}
