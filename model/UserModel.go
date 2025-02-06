package model

type User struct {
	UserID     uint   `gorm:"primaryKey" json:"user_id"`
	Username   string `gorm:"column:username" json:"username"`
	Password   string `gorm:"column:password" json:"password"`
	Is_Expert  int    `gorm:"column:is_expert" json:"is_expert"`
	Is_Premium int    `gorm:"column:is_premium" json:"is_premium"`
	Token      string `json:"token"`
}

func (User) TableName() string {
	return "users"
}
