package dto

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type PlaceOrderFromCart struct {
	UserID int64 `json:"userid"`
	CartID int64 `json:""cartid`
}

type ItemOrderedResponse struct {
	OrderID    int64   `json:"orderid"`
	ProductID  int64   `json:"productid"`
	Quantity   int64   `json:"quantity"`
	CatagoryID int64   `json:"catagoryid"`
	BrandName  string  `json:"brandname"`
	TotalPrice float64 `json:"totalprice"`
}

func (args *PlaceOrderFromCart) Parse(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&args)
	if err != nil {
		return err
	}
	return nil
}

func (args *PlaceOrderFromCart) Validate() error {
	validate := validator.New()
	err := validate.Struct(args)
	if err != nil {
		return err
	}
	return nil
}
