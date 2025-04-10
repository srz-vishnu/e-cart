package dto

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type UserDetailSaveRequest struct {
	UserID   int64  `json:"userid"`
	UserName string `json:"username" validate:"required"`
	Mail     string `json:"mail" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Pincode  int64  `json:"pincode" validate:"required"`
	Phone    int64  `json:"phonenumber" validate:"required"`
	Password string `json:"password" validate:"required"`
	IsAdmin  bool   `json:"isadmin"`
}

type SaveUserResponse struct {
	UserId int64 `json:"userid"`
}

func (args *UserDetailSaveRequest) Parse(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&args)
	if err != nil {
		return err
	}
	return nil
}

func (args *UserDetailSaveRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(args)
	if err != nil {
		return err
	}
	return nil
}
