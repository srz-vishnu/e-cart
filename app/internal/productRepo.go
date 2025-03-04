package internal

import (
	"e-cart/app/dto"
	"time"

	"gorm.io/gorm"
)

type ProductRepo interface {
	CreateProductDetail(args *dto.CreateProductDetailRequest) (int64, error)
}

type ProductRepoImpl struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepo {
	return &ProductRepoImpl{
		db: db,
	}
}

type Productdetail struct {
	ID           int64      `gorm:"primaryKey"`
	Productname  string     `gorm:"column:productname;unique;not null"`
	Description  string     `gorm:"column:description;not null"`
	Productprice float64    `gorm:"column:productPrice;not null"`
	Stockcount   int64      `gorm:"column:stockcount;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	IsDeleted    bool       `gorm:"column:is_deleted;default:false"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (r *ProductRepoImpl) CreateProductDetail(args *dto.CreateProductDetailRequest) (int64, error) {
	product := Productdetail{
		ID:           args.ProductID,
		Productname:  args.ProductName,
		Description:  args.Description,
		Productprice: args.ProductPrice,
		Stockcount:   args.StockCount,
	}
	//GORM's Create method to insert the new user
	if err := r.db.Table("productdetails").Create(&product).Error; err != nil {
		return 0, err
	}
	return product.ID, nil
}
