package dto

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type CreateProductDetailRequest struct {
	ProductID    int64   `json:"productid"`
	ProductName  string  `json:"productname" validate:"required"`
	Description  string  `json:"description"`
	ProductPrice float64 `json:"productprice" validae:"required"`
	StockCount   int64   `json:"stockcount" validate:"required"`
}

type CreateProductResponds struct {
	ProductID int64 `json:"productid"`
}

func (args *CreateProductDetailRequest) Parse(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&args)
	if err != nil {
		return err
	}
	return nil
}

func (args *CreateProductDetailRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(args)
	if err != nil {
		return err
	}
	return nil
}
