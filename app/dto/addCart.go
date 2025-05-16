package dto

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type AddItemToCart struct {
	//UserID     int64 `json:"userid"`
	CategoryID int64 `json:"category_id"`
	Quantity   int64 `json:"quantity"`
	BrandId    int64 `json:"brandid"`
}

type CartItemResponse struct {
	UserID     int64   `json:"userid"`
	ProductID  int64   `json:"productid"`
	Quantity   int64   `json:"quantity"`
	BrandName  string  `json:"brandname"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"totalprice"`
}

func (args *AddItemToCart) Parse(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&args)
	if err != nil {
		return err
	}
	return nil
}

func (args *AddItemToCart) Validate() error {
	validate := validator.New()
	err := validate.Struct(args)
	if err != nil {
		return err
	}
	return nil
}
