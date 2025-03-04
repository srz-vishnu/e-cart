package internal

import (
	"e-cart/app/dto"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type UserRepo interface {
	SaveUserDetails(args *dto.UserDetailSaveRequest) (int64, error)
	UpdateUserDetails(args *dto.UpdateUserDetailRequest, UserId int64) error
}

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &UserRepoImpl{
		db: db,
	}
}

type Userdetail struct {
	ID          int64     `gorm:"primaryKey"`
	Username    string    `gorm:"column:username;unique;not null"`
	Password    string    `gorm:"column:password;not null"`
	Address     string    `gorm:"column:address;not null"`
	Pincode     int64     `gorm:"column:pincode;not null"`
	Phonenumber int64     `gorm:"column:phone_number; not null"`
	Mail        string    `gorm:"column:mail;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
	// IsDeleted   bool       `gorm:"column:is_deleted;default:false"`
	// DeletedAt   *time.Time `gorm:"column:deleted_at"`
	// DeletedBy   *int64     `gorm:"column:deleted_by"`
}

func (r *UserRepoImpl) SaveUserDetails(args *dto.UserDetailSaveRequest) (int64, error) {

	user := Userdetail{
		ID:       args.UserID,
		Address:  args.Address,
		Mail:     args.Mail,
		Username: args.UserName,
		Password: args.Password,
		Pincode:  args.Pincode,
	}
	//GORM's Create method to insert the new user
	if err := r.db.Table("userdetails").Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *UserRepoImpl) UpdateUserDetails(args *dto.UpdateUserDetailRequest, UserId int64) error {

	updates := map[string]interface{}{
		"username":   args.UserName,
		"mail":       args.Mail,
		"address":    args.Address,
		"password":   args.Password,
		"pincode":    args.Pincode,
		"updated_at": time.Now(),
	}

	result := r.db.Table("userdetails").Where("id=?", UserId).Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	// Check if any rows were updated
	if result.RowsAffected == 0 {
		return fmt.Errorf("no active user found with ID %d to update", UserId)
	}

	log.SetFlags(0)
	log.Println("user password updated successfully..")
	return nil
}
