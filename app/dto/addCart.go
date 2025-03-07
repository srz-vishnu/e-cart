package dto

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type AddItemToCart struct {
	UserID    int64   `json:"userid"`
	ProductID int64   `json:"productid"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
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
