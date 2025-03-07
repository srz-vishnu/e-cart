package dto

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type UpdateUserDetailRequest struct {
	UserID   int64  `json:"userid"`
	UserName string `json:"username" validate:"required"`
	Mail     string `json:"mail" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Pincode  int64  `json:"pincode" validate:"required"`
	Phone    int64  `json:"phonenumber" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (args *UpdateUserDetailRequest) Parse(r *http.Request) error {
	strID := chi.URLParam(r, "userid")
	intID, err := strconv.Atoi(strID)
	if err != nil {
		return err
	}
	args.UserID = int64(intID)

	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		return err
	}

	return nil
}

func (args *UpdateUserDetailRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(args)
	if err != nil {
		return err
	}
	return nil
}
