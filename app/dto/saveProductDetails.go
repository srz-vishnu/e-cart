package dto

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

type CreateCategoryDetailRequest struct {
	CategoryID   int64                `json:"categoryid"`
	CategoryName string               `json:"categoryname" validate:"required"`
	Description  string               `json:"description"`
	Brands       []BrandDetailRequest `json:"brands" validate:"required"`
}

type BrandDetailRequest struct {
	BrandName  string  `json:"brandname" validate:"required"`
	Price      float64 `json:"price" validate:"required"`
	StockCount int64   `json:"stockcount" validate:"required"`
	ImageLink  string  `json:"imagelink" validate:"required"`
}

// type CreateProductResponds struct {
// 	ProductID int64 `json:"productid"`
// }

type CreateProductResponds struct {
	ProductID   int64           `json:"product_id"`
	Category    string          `json:"category_name"`
	Description string          `json:"description"`
	Brands      []BrandResponse `json:"brands"`
}

type BrandResponse struct {
	BrandName  string  `json:"brand_name"`
	Price      float64 `json:"price"`
	StockCount int64   `json:"stock_count"`
	ImageLink  string  `json:"imagelink" validate:"required"`
}

func (args *CreateCategoryDetailRequest) Parse(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&args)
	if err != nil {
		return err
	}

	// here we convert the incoming Category_name to Upper case format
	args.CategoryName = strings.ToUpper(args.CategoryName)

	return nil
}

func (args *CreateCategoryDetailRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(args)
	if err != nil {
		return err
	}
	return nil
}
