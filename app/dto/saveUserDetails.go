package dto

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
}

type SaveUserResponse struct {
	UserId int64 `json:"userid"`
}

func (args *UserDetailSaveRequest) Parse(r *http.Request) error {
	// Extract the 'id' from the URL
	strID := chi.URLParam(r, "id")
	if strID == "" {
		return fmt.Errorf("id parameter is missing")
	}

	// Convert the string ID to an integer (or another type depending on your ID type)
	intID, err := strconv.Atoi(strID)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	// Store the parsed ID into your struct if needed
	args.UserID = int64(intID)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&args)
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

// func (args *UserDetailSaveRequest) Parse(r *http.Request) error {
// 	// Extract the 'id' from the URL
// 	strID := chi.URLParam(r, "id")
// 	if strID == "" {
// 		return fmt.Errorf("id parameter is missing")
// 	}

// 	// Convert the string ID to an integer (or another type depending on your ID type)
// 	intID, err := strconv.Atoi(strID)
// 	if err != nil {
// 		return fmt.Errorf("invalid id: %v", err)
// 	}

// 	// Store the parsed ID into your struct if needed
// 	args.ID = int64(intID)  // Assuming args.ID is an int64, modify as necessary

// 	// Decode the request body into the struct
// 	decoder := json.NewDecoder(r.Body)
// 	err = decoder.Decode(&args)
// 	if err != nil {
// 		return fmt.Errorf("error decoding request body: %v", err)
// 	}

// 	return nil
// }
