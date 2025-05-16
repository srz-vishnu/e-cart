package dto

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type PlaceOrderFromCart struct {
	//UserID int64 `json:"userid"`
	CartID int64 `json:""cartid"`
}

// type ItemOrderedResponse struct {
// 	OrderID    int64               `json:"order_id"`
// 	TotalPrice float64             `json:"total_price"`
// 	Items      []OrderItemResponse `json:"items"`
// }

// type OrderItemResponse struct {
// 	ProductID  int64   `json:"product_id"`
// 	Quantity   int64   `json:"quantity"`
// 	CategoryID int64   `json:"category_id"`
// 	BrandName  string  `json:"brand_name"`
// 	Price      float64 `json:"price"`
// }

type UserDetailsResponse struct {
	Username    string `json:"username"`
	Address     string `json:"address"`
	Pincode     int64  `json:"pincode"`
	PhoneNumber int64  `json:"phone_number"`
	Email       string `json:"email"`
}

type OrderItemResponse struct {
	ProductID  int64   `json:"product_id"`
	Quantity   int64   `json:"quantity"`
	CategoryID int64   `json:"category_id"`
	BrandName  string  `json:"brand_name"`
	Price      float64 `json:"price"`
}

type ItemOrderedResponse struct {
	OrderID     int64               `json:"order_id"`
	TotalPrice  float64             `json:"total_price"`
	UserDetails UserDetailsResponse `json:"user_details"`
	Items       []OrderItemResponse `json:"items"`
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
