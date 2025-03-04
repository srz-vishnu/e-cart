package domain

import "time"

type User struct {
	ID          int64     `gorm:"primaryKey"`
	Username    string    `gorm:"column:username;unique;not null"`
	Password    string    `gorm:"column:password;not null"`
	Address     string    `gorm:"column:address;not null"`
	Pincode     int64     `gorm:"column:pincode;not null"`
	Phonenumber int64     `gorm:"column:phone_number; not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
	// IsDeleted   bool       `gorm:"column:is_deleted;default:false"`
	// DeletedAt   *time.Time `gorm:"column:deleted_at"`
	// DeletedBy   *int64     `gorm:"column:deleted_by"`
}
