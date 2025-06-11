package dto

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentpassword" validate:"required"`
	NewPassword     string `json:"newpassword" validate:"required"`
	ConfirmPassword string `json:"confirmpassword" validate:"required"`
}

func (args *ChangePasswordRequest) Parse(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&args)
	if err != nil {
		return err
	}
	return nil
}

func (args *ChangePasswordRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(args)
	if err != nil {
		return err
	}
	return nil
}
