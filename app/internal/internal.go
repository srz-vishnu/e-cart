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
	AddItemToCart(args *dto.AddItemToCart) error
	//UserCart()
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
	Status      bool      `gorm:"column:status;default:true;not null"` // Boolean field, default true
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
	// IsDeleted   bool       `gorm:"column:is_deleted;default:false"`
	// DeletedAt   *time.Time `gorm:"column:deleted_at"`
	// DeletedBy   *int64     `gorm:"column:deleted_by"`
}

type Cart struct {
	ID        int64      `gorm:"primaryKey"`
	UserID    int64      `gorm:"column:user_id;not null"`    // Foreign key to User
	ProductID int64      `gorm:"column:product_id;not null"` // Foreign key to Brand
	Quantity  int64      `gorm:"column:quantity;not null;default:1"`
	Price     float64    `gorm:"column:price;not null"`
	User      Userdetail `gorm:"foreignKey:UserID;references:ID"`    // Relation to Userdetail table ()
	Product   Brand      `gorm:"foreignKey:ProductID;references:ID"` // Relation to Brand (Product) table (ProductID in the currenet table is a foreign key to the table brand, also a table le id anne ID)
	CreatedAt time.Time
	UpdatedAt time.Time
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

func (r *UserRepoImpl) AddItemToCart(args *dto.AddItemToCart) error {

	item := Cart{
		UserID:    args.UserID,
		ProductID: args.ProductID,
		Quantity:  args.Quantity,
		Price:     args.Price,
	}

	if err := r.db.Table("carts").Create(&item).Error; err != nil {
		return err
	}

	return nil
}

// func (r *UserRepoImpl) AddToCart(userID, productID int64, quantity int64) error {
// 	var cart Cart
// 	var product Brand

// 	// Fetch product details to get the price
// 	if err := r.db.Table("brands").Where("id = ?", productID).First(&product).Error; err != nil {
// 		return err
// 	}

// 	// Check if the product already exists in the user's cart
// 	if err := r.db.Table("carts").Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error; err != nil {
// 		// If not found, add a new item to the cart
// 		if err == gorm.ErrRecordNotFound {
// 			cart = Cart{
// 				UserID:    userID,
// 				ProductID: productID,
// 				Quantity:  quantity,
// 				Price:     product.Price, // Set the price from the product
// 			}
// 			if err := r.db.Table("carts").Create(&cart).Error; err != nil {
// 				return err
// 			}
// 		} else {
// 			return err
// 		}
// 	} else {
// 		// If found, update the quantity and price remains the same
// 		cart.Quantity += quantity
// 		if err := r.db.Table("carts").Save(&cart).Error; err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
