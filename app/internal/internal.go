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
	GetUserByUsername(username string) (*Userdetail, error)
	UpdateUserDetails(args *dto.UpdateUserDetailRequest, UserId int64) error
	//	AddItemToCart(args *dto.AddItemToCart) error
	IsUserActive(userID int64) (bool, error)
	GetProductDetails(productID, categoryID, quantity int64) (*Brand, error)
	CheckProductInCart(userID, productID int64) (*Cart, error)
	FetchCartItems(userID, cartID int64) ([]Cart, error)
	//CartWithProductDetails(UserID, ProductID int64) (*Cart, error)
	AddOrUpdateCart(userID int64, product *Brand, quantity int64) error
	GetCartWithProductDetails(userID int64, productID int64) (*Cart, error)
	CreateOrder(userID int64, totalAmount float64, cartItems []Cart) (*Order, error)
	CreateOrderItems(orderID int64, cartItems []Cart) ([]OrderItem, error)
	UpdateCartOrderStatus(userID, orderID, cartID int64) error
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
	ID          int64      `gorm:"primaryKey"`
	UserID      int64      `gorm:"column:user_id;not null"`    // Foreign key to User
	ProductID   int64      `gorm:"column:product_id;not null"` // Foreign key to Brand
	Quantity    int64      `gorm:"column:quantity;not null;default:1"`
	Price       float64    `gorm:"column:price;not null"`
	OrderStatus bool       `gorm:"column:orderstatus;not null;default:true"` // true if its active , once place cheytha false avanam
	OrderDetail int64      `gorm:"column:orderdetail;"`                      // orderstatus false akumbol, orderid add cheyanam
	User        Userdetail `gorm:"foreignKey:UserID;references:ID"`          // Relation to Userdetail table
	BrandID     int64      `gorm:"column:id"`                                // Foreign key to Brand
	Brand       Brand      `gorm:"foreignKey:ProductID"`                     // Assuming ProductID references Brand.ID

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID        int64       `gorm:"primaryKey"`
	UserID    int64       `gorm:"index;not null"` // Foreign key to Userdetail
	Total     float64     `gorm:"not null"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
	User      Userdetail  `gorm:"foreignKey:UserID;references:ID"` // Relation to Userdetail table
	Items     []OrderItem `gorm:"foreignKey:OrderID"`              // One-to-many relation with OrderItem
}

type OrderItem struct {
	ID        int64   `gorm:"primaryKey"`
	OrderID   int64   `gorm:"index;not null"` // Foreign key to Order
	ProductID int64   `gorm:"not null"`       // Foreign key to Product (Brand)
	Quantity  int64   `gorm:"not null"`
	Price     float64 `gorm:"not null"`
	Order     Order   `gorm:"foreignKey:OrderID;references:ID"`   // Relation to Order
	Product   Brand   `gorm:"foreignKey:ProductID;references:ID"` // Relation to Brand, orderid is foreign key to order table, a table le primary id anne ivide reference id ayite irikane
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

func (r *UserRepoImpl) GetUserByUsername(username string) (*Userdetail, error) {
	var user Userdetail
	if err := r.db.Table("userdetails").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
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

// func (r *UserRepoImpl) AddItemToCart(args *dto.AddItemToCart) error {

// 	item := Cart{
// 		UserID:    args.UserID,
// 		ProductID: args.ProductID,
// 		Quantity:  args.Quantity,
// 		//Price:     args.Price,
// 	}

// 	if err := r.db.Table("carts").Create(&item).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

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

func (r *UserRepoImpl) GetProductDetails(brandID, categoryID, quantity int64) (*Brand, error) {
	var product Brand

	if err := r.db.Table("brands").Where("id = ? AND category_id = ? ", categoryID, brandID).First(&product).Error; err != nil {
		// If product not found, then GORM error
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product with id %d not found", brandID)
		}
		return nil, err
	}

	// Check if the requested quantity is available in total stock
	if quantity > product.StockCount {
		return nil, fmt.Errorf("requested quantity %d exceeds available stock %d for product ID %d with Brand ID %d", quantity, product.StockCount, categoryID, brandID)
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

func (r *UserRepoImpl) AddOrUpdateCart(userID int64, product *Brand, quantity int64) error {
	var existingCart Cart

	// Check if the product already exists in the user's cart
	if err := r.db.Table("carts").Where("user_id = ? AND product_id = ?", userID, product.ID).First(&existingCart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If cart item does not exist, create a new cart item
			newCart := Cart{
				UserID:    userID,
				ProductID: product.ID,
				Quantity:  quantity,
				Price:     product.Price,
				// Product: Brand{
				// 	CategoryID: product.CategoryID,
				// 	BrandName:  product.BrandName,
				// },
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

	// if err := r.db.Table("carts").Preload("brands").Where("user_id = ? AND id = ?", userID, cartID).Find(&cartItems).Error; err != nil {
	// 	return nil, err
	// }

	if err := r.db.Preload("Brand").Where("user_id = ? AND id = ?", userID, cartID).Find(&cartItems).Error; err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, fmt.Errorf("no items in the cart to place an order for the given userID and cartID")
	}

	return cartItems, nil
}

// func (r *UserRepoImpl) CartWithProductDetails(UserID, ProductID int64) (*Cart, error) {
// 	var cart Cart

// 	// Use Preload to load the related product (Brand) details
// 	if err := r.db.Preload("Product").Where("user_id = ? AND product_id = ?", UserID, ProductID).First(&cart).Error; err != nil {
// 		return nil, err
// 	}
// 	return &cart, nil
// }

func (r *UserRepoImpl) CreateOrder(userID int64, totalAmount float64, cartItems []Cart) (*Order, error) {
	// Create a new Order
	newOrder := &Order{
		UserID: userID,
		Total:  totalAmount,
	}

	// Create the order in the database
	if err := r.db.Create(newOrder).Error; err != nil {
		return nil, err
	}

	// Create OrderItems for each cart item
	for _, item := range cartItems {
		orderItem := OrderItem{
			OrderID:   newOrder.ID,
			ProductID: item.ProductID, // Just store the reference
			Quantity:  item.Quantity,
			Price:     item.Price,
			// Remove the Product field - we don't need to store brand info here
		}

		// Save the OrderItem
		if err := r.db.Create(&orderItem).Error; err != nil {
			return nil, err
		}
	}

	return newOrder, nil
}

// func (r *UserRepoImpl) CreateOrder(userID int64, totalAmount float64, cartItems []Cart) (*Order, error) {
// 	// Create a new Order
// 	newOrder := &Order{
// 		UserID: userID,
// 		Total:  totalAmount,
// 	}

// 	// Create the order in the database
// 	if err := r.db.Create(newOrder).Error; err != nil {
// 		return nil, err
// 	}

// 	// Create OrderItems for each cart item
// 	for _, item := range cartItems {
// 		orderItem := OrderItem{
// 			OrderID:   newOrder.ID,
// 			ProductID: item.ProductID,
// 			Quantity:  item.Quantity,
// 			Price:     item.Price,
// 			Product:   Brand{BrandName: item.Brand.BrandName},
// 		}

// 		// Save the OrderItem
// 		if err := r.db.Create(&orderItem).Error; err != nil {
// 			return nil, err
// 		}
// 	}

// 	return newOrder, nil
// }

// CreateOrderItems adds each cart item to the order.
func (r *UserRepoImpl) CreateOrderItems(orderID int64, cartItems []Cart) ([]OrderItem, error) {

	var createdItems []OrderItem

	for _, item := range cartItems {
		orderItem := OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Product: Brand{
				BrandName:  item.Brand.BrandName,
				CategoryID: item.Brand.CategoryID,
			},
		}

		if err := r.db.Table("order_items").Create(&orderItem).Error; err != nil {
			return nil, err
		}
		createdItems = append(createdItems, orderItem)
	}

	return createdItems, nil
}

// ClearCart deletes all items from the cart for the given user.
func (r *UserRepoImpl) ClearCart(userID int64) error {
	if err := r.db.Table("carts").Where("user_id = ?", userID).Delete(&Cart{}).Error; err != nil {
		return err
	}
	return nil
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
