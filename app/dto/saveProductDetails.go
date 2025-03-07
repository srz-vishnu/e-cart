package dto

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

// type CreateProductDetailRequest struct {
// 	ProductID    int64   `json:"productid"`
// 	ProductName  string  `json:"productname" validate:"required"`
// 	Description  string  `json:"description"`
// 	ProductPrice float64 `json:"productprice" validae:"required"`
// 	StockCount   int64   `json:"stockcount" validate:"required"`
// }

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
}

type CreateProductResponds struct {
	ProductID int64 `json:"productid"`
}

func (args *CreateCategoryDetailRequest) Parse(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&args)
	if err != nil {
		return err
	}
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
