package internal

import (
	"e-cart/app/dto"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo interface {
	SaveUserDetails(args *dto.UserDetailSaveRequest) (int64, error)
	GetUserByUsername(username string) (*Userdetail, error)
	ChangePassword(userID int64, hashedPwd string) error
	GetUserDetailByID(userId int64) (*Userdetail, error)
	SaveToken(userId int64, token string, expiry time.Time) error
	UpdateUserDetails(args *dto.UpdateUserDetailRequest, UserId int64) error
	IsUserActive(userID int64) (bool, error)
	GetProductDetails(productID, categoryID int64) (*Brand, error)
	CheckProductInCart(userID, productID int64) (*Cart, error)
	FetchCartItems(userID, cartID int64) ([]Cart, error)
	AddOrUpdateCart(userID int64, product *Brand, quantity int64, totalAmount float64) error
	GetCartWithProductDetails(userID int64, productID int64) (*Cart, error)
	UpdateCartOrderStatus(userID, orderID, cartID int64) error
	UpdateStockCount(orderItems []OrderItem) ([]Brand, error)
	ViewCart(userID int64) ([]Cart, error)
	ClearCart(userID int64) error
	CreateOrder(userID int64, totalAmount float64, cartItems []Cart) (*Order, []OrderItem, error)
	GetUserByID(userID int64) (*Userdetail, error)
	GetOrderHistoryByUserID(userID int64) ([]Order, error)
	AddOrUpdateFavorite(userID int64, args dto.UserFavoriteBrandRequest) error
	GetFavoriteBrandIDs(userID int64) ([]int64, error)
	GetBrandsByIDs(brandIDs []int64) ([]Brand, error)
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
	IsAdmin     bool      `gorm:"column:isadmin;default:false;not null"` // default false for user, true is used when  admin logins
}

type ActiveToken struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	UserID    int64     `gorm:"column:user_id;unique;not null"`
	Token     string    `gorm:"column:token;not null"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Userdetail) TableName() string {
	return "userdetails"
}

type Cart struct {
	ID          int64   `gorm:"primaryKey"`
	UserID      int64   `gorm:"column:user_id;not null"`
	ProductID   int64   `gorm:"column:product_id;not null"` // Foreign key to Brand.ID
	Quantity    int64   `gorm:"column:quantity;not null;default:1"`
	Price       float64 `gorm:"column:price;not null"`
	TotalAmount float64 `gorm:"column:totalamount;not null; default:0"`
	OrderStatus bool    `gorm:"column:orderstatus;not null;default:true"`
	OrderDetail int64   `gorm:"column:orderdetail;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Brand       Brand `gorm:"foreignKey:ProductID"` // Relationship to Brand
}
type Order struct {
	ID        int64       `gorm:"primaryKey"`
	UserID    int64       `gorm:"index;not null"` // Foreign key to Userdetail
	Total     float64     `gorm:"not null"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
	User      Userdetail  `gorm:"foreignKey:UserID;references:ID"` // Relation to Userdetail table
	Items     []OrderItem `gorm:"foreignKey:OrderID"`              // One-to-many relation with OrderItem
	UpdatedAt time.Time   `gorm:"column:updated_at;autoUpdateTime"`
}

type OrderItem struct {
	ID        int64     `gorm:"primaryKey"`
	OrderID   int64     `gorm:"index;not null"` // Foreign key to Order
	ProductID int64     `gorm:"not null"`       // Foreign key to Product (Brand)
	Quantity  int64     `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	Order     Order     `gorm:"foreignKey:OrderID;references:ID"`   // Relation to Order
	Product   Brand     `gorm:"foreignKey:ProductID;references:ID"` // Relation to Brand, orderid is foreign key to order table, a table le primary id anne ivide reference id ayite irikane
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type UserFavoriteBrand struct {
	ID       int64      `gorm:"primaryKey"`
	UserID   int64      `gorm:"column:user_id;not null"`
	User     Userdetail `gorm:"foreignKey:UserID"`
	BrandID  int64      `gorm:"column:brand_id;not null"`
	Brand    Brand      `gorm:"foreignKey:BrandID"`
	Favorite bool       `gorm:"column:favorite;default:false;not null"`
}

func (r *UserRepoImpl) SaveUserDetails(args *dto.UserDetailSaveRequest) (int64, error) {

	user := Userdetail{
		//ID:          args.UserID,
		Address:     args.Address,
		Mail:        args.Mail,
		Username:    args.UserName,
		Password:    args.Password,
		Pincode:     args.Pincode,
		Phonenumber: args.Phone,
		IsAdmin:     args.IsAdmin,
	}
	//GORM's Create method to insert the new user
	if err := r.db.Table("userdetails").Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *UserRepoImpl) GetUserByUsername(username string) (*Userdetail, error) {
	var user Userdetail
	if err := r.db.Table("userdetails").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepoImpl) SaveToken(userId int64, token string, expiry time.Time) error {
	saveToken := ActiveToken{
		UserID:    userId,
		Token:     token,
		ExpiresAt: expiry,
	}
	//GORM's Create method to insert the token
	// if err := r.db.Table("active_tokens").Create(&saveToken).Error; err != nil {
	// 	return err
	// }
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}}, // unique constraint
		DoUpdates: clause.AssignmentColumns([]string{"token", "expires_at", "updated_at"}),
	}).Create(&saveToken).Error
}

func (r *UserRepoImpl) GetUserDetailByID(userId int64) (*Userdetail, error) {
	var user Userdetail
	if err := r.db.Table("userdetails").Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepoImpl) UpdateUserDetails(args *dto.UpdateUserDetailRequest, UserId int64) error {

	updates := map[string]interface{}{
		"username": args.UserName,
		"mail":     args.Mail,
		"address":  args.Address,
		//"password":   args.Password,  commented password bcoz password updation should be done using an seperate api
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
	log.Println("user details updated successfully..")
	return nil
}

func (r *UserRepoImpl) ChangePassword(userID int64, hashedPwd string) error {
	// Use GORM's Update to modify only the password field
	if err := r.db.Table("userdetails").Where("id = ?", userID).Update("password", hashedPwd).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepoImpl) IsUserActive(userID int64) (bool, error) {
	var user Userdetail

	// Fetch the user details by userID
	if err := r.db.Table("userdetails").Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, fmt.Errorf("user not found")
		}
		return false, err
	}

	return user.Status, nil
}

func (r *UserRepoImpl) GetProductDetails(brandID, categoryID int64) (*Brand, error) {
	var product Brand

	if err := r.db.Table("brands").Where("id = ? AND category_id = ? ", categoryID, brandID).First(&product).Error; err != nil {
		// If product not found, then GORM error
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product with id %d not found", brandID)
		}
		return nil, err
	}

	return &product, nil
}

func (r *UserRepoImpl) CheckProductInCart(userID, productID int64) (*Cart, error) {
	var cart Cart

	// Check if the product exists in the user's cart
	if err := r.db.Table("carts").Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error; err != nil {
		// If not found, return nil
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (r *UserRepoImpl) AddOrUpdateCart(userID int64, product *Brand, quantity int64, totalAmount float64) error {
	var existingCart Cart

	// Check if the product already exists in the user's cart
	if err := r.db.Table("carts").Where("user_id = ? AND product_id = ?", userID, product.ID).First(&existingCart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If cart item does not exist, create a new cart item
			newCart := Cart{
				UserID:      userID,
				ProductID:   product.ID,
				Quantity:    quantity,
				Price:       product.Price,
				TotalAmount: totalAmount,
				//Brand:     Brand,
			}
			if err := r.db.Table("carts").Create(&newCart).Error; err != nil {
				return err
			}
			// Return the newly created cart item
			return nil
		} else {
			return err
		}
	}

	// If product already exists in the cart, update the quantity
	existingCart.Quantity += quantity
	if err := r.db.Table("carts").Save(&existingCart).Error; err != nil {
		return err
	}

	// Return the updated cart item
	return nil
}

func (r *UserRepoImpl) GetCartWithProductDetails(userID int64, productID int64) (*Cart, error) {
	var cart Cart

	// Use Preload to load the related Brand (product) details
	if err := r.db.Preload("Brand").Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error; err != nil {
		return nil, err
	}

	return &cart, nil
}

// FetchCartItems retrieves the user's cart items.
func (r *UserRepoImpl) FetchCartItems(userID, cartID int64) ([]Cart, error) {
	var cartItems []Cart

	if err := r.db.Preload("Brand").Where("user_id = ? AND id = ?", userID, cartID).Find(&cartItems).Error; err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, fmt.Errorf("no items in the cart to place an order for the given userID and cartID")
	}

	return cartItems, nil
}

// UpdateCartOrder updates the order status to false and adds the order ID to orderdetail
func (r *UserRepoImpl) UpdateCartOrderStatus(userID, orderID, cartID int64) error {
	// Create a map for fields to update
	updates := map[string]interface{}{
		"OrderStatus": false,
		"OrderDetail": orderID,
	}

	// Updating
	result := r.db.Model(&Cart{}).Where("user_id = ? AND id = ?", userID, cartID).Updates(updates)

	// Check for errors
	if result.Error != nil {
		return result.Error
	}

	// Check if rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

func (r *UserRepoImpl) UpdateStockCount(orderItems []OrderItem) ([]Brand, error) {
	var updatedBrands []Brand
	var productIDs []int64

	// Collect all product IDs first
	for _, item := range orderItems {
		productIDs = append(productIDs, item.ProductID)
	}

	// Preload all brands in a single query
	var brands []Brand
	if err := r.db.Where("id IN ?", productIDs).Find(&brands).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch product details: %w", err)
	}

	// Create a map for quick lookup
	brandMap := make(map[int64]Brand)
	for _, brand := range brands {
		brandMap[brand.ID] = brand
	}

	// Update stock counts
	for _, item := range orderItems {
		brand, exists := brandMap[item.ProductID]
		if !exists {
			return nil, fmt.Errorf("product ID %d not found", item.ProductID)
		}

		if brand.StockCount < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product ID %d (current: %d, required: %d)",
				item.ProductID, brand.StockCount, item.Quantity)
		}

		brand.StockCount -= item.Quantity
		if err := r.db.Save(&brand).Error; err != nil {
			return nil, fmt.Errorf("failed to update stock for product ID %d: %w", item.ProductID, err)
		}

		updatedBrands = append(updatedBrands, brand)
	}

	return updatedBrands, nil
}

// ClearCart deletes all items from the cart for the given user.
func (r *UserRepoImpl) ViewCart(userID int64) ([]Cart, error) {
	var cartItems []Cart
	if err := r.db.Preload("Brand").Where("user_id = ? AND orderstatus = ?", userID, true).Find(&cartItems).Error; err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (r *UserRepoImpl) ClearCart(userID int64) error {
	var cartItems Cart

	result := r.db.Preload("Brand").Where("user_id = ? AND orderstatus = ?", userID, true).Delete(&cartItems)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no cart items found for this user")
	}
	return nil
}

func (r *UserRepoImpl) GetUserByID(userID int64) (*Userdetail, error) {
	var user Userdetail
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepoImpl) CreateOrder(userID int64, totalAmount float64, cartItems []Cart) (*Order, []OrderItem, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the order (using tx)
	newOrder := &Order{
		UserID: userID,
		Total:  totalAmount,
	}

	if err := tx.Create(newOrder).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	// Create order items (using tx)
	var createdItems []OrderItem
	for _, item := range cartItems {
		// Get brand details (using tx)
		var brand Brand
		if err := tx.Where("id = ?", item.ProductID).First(&brand).Error; err != nil {
			tx.Rollback()
			return nil, nil, fmt.Errorf("failed to get brand details: %w", err)
		}

		orderItem := OrderItem{
			OrderID:   newOrder.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Product:   brand,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, nil, err
		}

		createdItems = append(createdItems, orderItem)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, err
	}

	return newOrder, createdItems, nil
}

func (r *UserRepoImpl) GetOrderHistoryByUserID(userID int64) ([]Order, error) {
	var orders []Order

	err := r.db.
		Preload("Items").Preload("Items.Product").Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no orders found for user with ID %d", userID)
		}
		return nil, err
	}

	return orders, nil
}

func (r *UserRepoImpl) AddOrUpdateFavorite(userID int64, args dto.UserFavoriteBrandRequest) error {
	var fav UserFavoriteBrand

	// Check if the favorite entry already exists
	err := r.db.Table("user_favorite_brands").Where("user_id = ? AND brand_id = ?", userID, args.BrandID).First(&fav).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Not found: create a new record
		newFav := UserFavoriteBrand{
			UserID:   userID,
			BrandID:  args.BrandID,
			Favorite: args.Favorite,
		}

		if err := r.db.Table("user_favorite_brands").Create(&newFav).Error; err != nil {
			return err
		}
		return nil
	} else if err != nil {
		// Any other DB error
		return err
	}

	// Record exists: update the favorite field
	fav.Favorite = args.Favorite
	if err := r.db.Table("user_favorite_brands").Save(&fav).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepoImpl) GetFavoriteBrandIDs(userID int64) ([]int64, error) {
	var favs []UserFavoriteBrand
	err := r.db.Where("user_id = ? AND favorite = true", userID).Find(&favs).Error
	if err != nil {
		return nil, err
	}

	brandIDs := make([]int64, len(favs))
	for i, fav := range favs {
		brandIDs[i] = fav.BrandID
	}
	return brandIDs, nil
}

func (r *UserRepoImpl) GetBrandsByIDs(brandIDs []int64) ([]Brand, error) {
	if len(brandIDs) == 0 {
		return []Brand{}, nil
	}

	var brands []Brand
	err := r.db.Where("id IN ?", brandIDs).Find(&brands).Error
	if err != nil {
		return nil, err
	}
	return brands, nil
}
