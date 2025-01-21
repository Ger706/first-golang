package model

type User struct {
	UserID   uint   `gorm:"primaryKey"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	Token    string `json:"token"`
}

func (User) TableName() string {
	return "users"
}
