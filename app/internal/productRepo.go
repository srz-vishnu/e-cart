package internal

import (
	"e-cart/app/dto"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ProductRepo interface {
	CreateAndUpsertProductDetail(args *dto.CreateCategoryDetailRequest) (*Category, error)
	GetAllProducts() ([]Category, error)
	GetCategoryByID(categoryID int64) (*Category, error)
	GetCategoryByName(categoryName string) (*Category, error)
	GetAllBrands() ([]Brand, error)
	UpdateCategory(categoryID int64, newCategoryName string) error
	UpdateBrand(brandID int64, newBrandName string, newPrice float64) error
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
	Category    Category   `gorm:"foreignKey:CategoryID"`       // Add this to establish relation
	BrandName   string     `gorm:"column:brandname;not null"`
	Price       float64    `gorm:"column:price;not null"`
	StockCount  int64      `gorm:"column:stockcount;not null"`
	ImageLink   string     `gorm:"column:image_link"`
	ReleaseDate time.Time  `gorm:"column:release_date;not null"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	IsDeleted   bool       `gorm:"column:is_deleted;default:false"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

// To add or update product in to the list
func (r *ProductRepoImpl) CreateAndUpsertProductDetail(args *dto.CreateCategoryDetailRequest) (*Category, error) {
	var category Category

	normalizedCategoryName := strings.ToLower(args.CategoryName)

	// Check if category exists (case-insensitive)
	err := r.db.Table("categories").
		Where("LOWER(categoryname) = ?", normalizedCategoryName).
		First(&category).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			category = Category{
				ID:           args.CategoryID,
				Categoryname: strings.Title(normalizedCategoryName), // Normalize case
				Description:  args.Description,
			}
			if err := r.db.Table("categories").Create(&category).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		// Check ID match
		if category.ID != args.CategoryID {
			return nil, fmt.Errorf("category name '%s' already exists with ID %d, but request provided ID %d",
				category.Categoryname, category.ID, args.CategoryID)
		}
	}

	// Upsert brands
	for _, b := range args.Brands {
		var existingBrand Brand
		normalizedBrandName := strings.ToLower(b.BrandName)

		err := r.db.Table("brands").
			Where("category_id = ? AND LOWER(brandname) = ?", category.ID, normalizedBrandName).
			First(&existingBrand).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				newBrand := Brand{
					CategoryID: category.ID,
					BrandName:  strings.ToUpper(b.BrandName), // Store with consistent casing
					Price:      b.Price,
					StockCount: b.StockCount,
					ImageLink:  b.ImageLink,
				}
				if err := r.db.Table("brands").Create(&newBrand).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			// Update stock and price
			existingBrand.StockCount += b.StockCount
			existingBrand.Price = b.Price
			existingBrand.ImageLink = b.ImageLink
			if err := r.db.Table("brands").Save(&existingBrand).Error; err != nil {
				return nil, err
			}
		}
	}

	// Load updated brands to return complete category info
	var updatedBrands []Brand
	if err := r.db.Table("brands").Where("category_id = ?", category.ID).Find(&updatedBrands).Error; err != nil {
		return nil, err
	}
	category.Brands = updatedBrands

	return &category, nil
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
	if err := r.db.Preload("Brands").First(&category, categoryID).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *ProductRepoImpl) GetCategoryByName(categoryName string) (*Category, error) {
	var category Category
	if err := r.db.Preload("Brands").Where("categoryname = ?", categoryName).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *ProductRepoImpl) GetAllBrands() ([]Brand, error) {
	var brand []Brand

	if err := r.db.Preload("Category").Find(&brand).Error; err != nil {
		return nil, err
	}

	return brand, nil
}

// UpdateCategory updates the name of a category by its ID using GORM
func (r *ProductRepoImpl) UpdateCategory(categoryID int64, newCategoryName string) error {
	category := &Category{}
	if err := r.db.First(category, categoryID).Error; err != nil {
		return err
	}

	category.Categoryname = newCategoryName
	if err := r.db.Save(category).Error; err != nil {
		return err
	}
	return nil
}

// UpdateBrand updates the name of a brand by its ID using GORM
func (r *ProductRepoImpl) UpdateBrand(brandID int64, newBrandName string, newPrice float64) error {
	brand := &Brand{}
	if err := r.db.First(brand, brandID).Error; err != nil {
		return err
	}

	brand.BrandName = newBrandName
	if err := r.db.Save(brand).Error; err != nil {
		return err
	}
	return nil
}
