package domain

import "time"

type Demo struct {
	ID        int64  `gorm:"primaryKey"`
	Username  string `gorm:"column:username;unique;not null"`
	Password  string `gorm:"column:password;not null"`
	Address   string `gorm:"column:address;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}