package dto

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type UpdateBrand struct {
	BrandId   int64   `json:"brand_id"`
	BrandName string  `json:"brand_name"`
	Price     float64 `json:"price"`
}

func (args *UpdateBrand) Parse(r *http.Request) error {
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
	args.BrandId = int64(intID)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&args)
	if err != nil {
		return err
	}
	return nil
}

func (args *UpdateBrand) Validate() error {
	validate := validator.New()
	err := validate.Struct(args)
	if err != nil {
		return err
	}
	return nil
}
