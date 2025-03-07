package internal

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AdminRepo interface {
	BlockUser(userId int64) error
	UnBlockUser(userId int64) error
	GetAllUsers() ([]Userdetail, error)
}

type AdminRepoImpl struct {
	db *gorm.DB
}

func NewAdminRepo(db *gorm.DB) AdminRepo {
	return &AdminRepoImpl{
		db: db,
	}
}

func (r *AdminRepoImpl) BlockUser(userId int64) error {

	updates := map[string]interface{}{
		"status":     false,
		"updated_at": time.Now(),
	}

	result := r.db.Table("userdetails").Where("id = ?", userId).Updates(updates)

	// Check for errors during the update operation
	if result.Error != nil {
		return result.Error
	}

	// If no rows were affected, error indicating the user was not found
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d to block", userId)
	}

	return nil
}

func (r *AdminRepoImpl) UnBlockUser(userId int64) error {

	updates := map[string]interface{}{
		"status":     true,
		"updated_at": time.Now(),
	}
	result := r.db.Model("userdetails").Where("id = ?", userId).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	// If no rows were affected, error indicating the user was not found
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with id %d to unblock", userId)
	}
	return nil
}

func (r *AdminRepoImpl) GetAllUsers() ([]Userdetail, error) {

	var details []Userdetail

	result := r.db.Where("status = ?", true).Find(&details)

	if result.Error != nil {
		return nil, result.Error
	}

	return details, nil
}
