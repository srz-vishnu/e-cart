package internal

import (
	"e-cart/app/dto"
	"time"

	"gorm.io/gorm"
)

type ProductRepo interface {
	CreateAndUpsertProductDetail(args *dto.CreateCategoryDetailRequest) (int64, error)
	GetAllProducts() ([]Category, error)
	GetCategoryByID(categoryID int64) (*Category, error)
	GetAllBrands() ([]Brand, error)
}

type ProductRepoImpl struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepo {
	return &ProductRepoImpl{
		db: db,
	}
}

// Category represents a single category that can have multiple brands
type Category struct {
	ID           int64      `gorm:"primaryKey"`
	Categoryname string     `gorm:"column:categoryname;unique;not null"`
	Description  string     `gorm:"column:description;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	IsDeleted    bool       `gorm:"column:is_deleted;default:false"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
	Brands       []Brand    `gorm:"foreignKey:CategoryID"` // The "one" side (has many brands)
}

// Brand represents a brand that belongs to a category
type Brand struct {
	ID          int64      `gorm:"primaryKey"`
	CategoryID  int64      `gorm:"column:category_id;not null"` // Foreign key to Category
	BrandName   string     `gorm:"column:brandname;not null"`
	Price       float64    `gorm:"column:price;not null"`
	StockCount  int64      `gorm:"column:stockcount;not null"`
	ReleaseDate time.Time  `gorm:"column:release_date;not null"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	IsDeleted   bool       `gorm:"column:is_deleted;default:false"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

// To add new product in to the list
// func (r *ProductRepoImpl) CreateNewProductDetail(args *dto.CreateCategoryDetailRequest) (int64, error) {
// 	// Create the product first
// 	catagory := Category{
// 		ID:           args.CategoryID,
// 		Categoryname: args.CategoryName,
// 		Description:  args.Description,
// 	}

// 	// Insert the product into the productdetails table
// 	if err := r.db.Table("categories").Create(&catagory).Error; err != nil {
// 		return 0, err
// 	}

// 	// Create categories for this product
// 	for _, cat := range args.Brands {
// 		brand := Brand{
// 			CategoryID: catagory.ID,
// 			BrandName:  cat.BrandName,
// 			Price:      cat.Price,
// 			StockCount: cat.StockCount,
// 		}
// 		// Insert category for this product into the categories table
// 		if err := r.db.Table("brands").Create(&brand).Error; err != nil {
// 			return 0, err
// 		}
// 	}

// 	return catagory.ID, nil
// }

func (r *ProductRepoImpl) CreateAndUpsertProductDetail(args *dto.CreateCategoryDetailRequest) (int64, error) {
	var category Category

	// Check if the category exists by CategoryName
	if err := r.db.Table("categories").Where("categoryname = ?", args.CategoryName).First(&category).Error; err != nil {
		// If category not found, create a new category
		if err == gorm.ErrRecordNotFound {
			category = Category{
				ID:           args.CategoryID,
				Categoryname: args.CategoryName,
				Description:  args.Description,
			}
			// Insert the new category into the categories table
			if err := r.db.Table("categories").Create(&category).Error; err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	// Add or update brands for the category
	for _, cat := range args.Brands {
		var existingBrand Brand

		// Check if the brand already exists under the category
		if err := r.db.Table("brands").Where("category_id = ? AND brandname = ?", category.ID, cat.BrandName).First(&existingBrand).Error; err != nil {
			// If the brand does not exist, insert it as a new brand
			if err == gorm.ErrRecordNotFound {
				brand := Brand{
					CategoryID: category.ID,
					BrandName:  cat.BrandName,
					Price:      cat.Price,
					StockCount: cat.StockCount,
				}
				// Insert the new brand into the brands table
				if err := r.db.Table("brands").Create(&brand).Error; err != nil {
					return 0, err
				}
			} else {
				return 0, err
			}
		} else {
			// If  brand exists, update the stock count by adding the new stock count to the existing one
			existingBrand.StockCount += cat.StockCount

			// Option
			existingBrand.Price = cat.Price

			// Update the brand's stock count and price in the database
			if err := r.db.Table("brands").Save(&existingBrand).Error; err != nil {
				return 0, err
			}
		}
	}

	return category.ID, nil
}

func (r *ProductRepoImpl) GetAllProducts() ([]Category, error) {

	// Query the products only
	var products []Category
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepoImpl) GetCategoryByID(categoryID int64) (*Category, error) {
	var category Category
	// preload nokanam
	if err := r.db.Preload("Brands").First(&category, categoryID).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *ProductRepoImpl) GetAllBrands() ([]Brand, error) {
	var brand []Brand
	if err := r.db.Find(&brand).Error; err != nil {
		return nil, err
	}
	return brand, nil
}
