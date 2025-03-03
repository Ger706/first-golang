package model

import "time"

type User struct {
	UserID    uint       `gorm:"primaryKey" json:"user_id"`
	Username  string     `gorm:"column:username" json:"username"`
	Email     string     `gorm:"column:email" json:"email"`
	Password  string     `gorm:"column:password" json:"password"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	Token     string     `gorm:"-" json:"token"`
}

func (User) TableName() string {
	return "users"
}
