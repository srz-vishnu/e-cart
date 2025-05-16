package dto

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type UserFavoriteBrandRequest struct {
	BrandID  int64 `json:"brandid" validate:"required"`
	Favorite bool  `json:"favourite" validate:"required"`
}

func (args *UserFavoriteBrandRequest) Parse(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&args)
	if err != nil {
		return err
	}
	return nil
}

func (args *UserFavoriteBrandRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(args)
	if err != nil {
		return err
	}
	return nil
}
