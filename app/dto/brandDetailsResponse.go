package dto

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type SearchProductByNameRequest struct {
	CategoryName string
}

type CategoryDetailResponse struct {
	CategoryID   int64                `json:"categoryid"`
	CategoryName string               `json:"categoryname"`
	Description  string               `json:"description"`
	Brands       []BrandDetailRequest `json:"brands"`
}

type BrandDetailResponses struct {
	BrandName  string  `json:"brandname"`
	Price      float64 `json:"price"`
	StockCount int64   `json:"stockcount"`
}

func (args *SearchProductByNameRequest) Parse(r *http.Request) error {
	categoryName := chi.URLParam(r, "categoryname")

	if categoryName == "" {
		return fmt.Errorf("name parameter is missing or empty")
	}
	args.CategoryName = strings.ToUpper(categoryName)

	return nil
}
