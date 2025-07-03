package dto

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BrandFullDetailByIdRequest struct {
	BrandId int64 `json:"brandid"`
}

type BrandFullDetailByIdResponse struct {
	BrandId          int64    `json:"brandid"`
	BrandName        string   `json:"brandname"`
	Price            float64  `json:"price"`
	StockCount       int64    `json:"stockcount"`
	ImageLink        string   `json:"imagelink"`
	GalleryLinks     []string `json:"gallerylinks"`
	BrandDescription string   `json:"branddescription"`
	CategoryID       int64    `json:"categoryid"`
	CategoryName     string   `json:"categoryname"`
}

func (args *BrandFullDetailByIdRequest) Parse(r *http.Request) error {
	strID := chi.URLParam(r, "id")
	if strID == "" {
		return fmt.Errorf("id parameter is missing or empty")
	}
	intID, err := strconv.Atoi(strID)
	if err != nil {
		return err
	}
	args.BrandId = int64(intID)
	return nil
}
