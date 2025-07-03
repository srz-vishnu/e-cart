package dto

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CatagoryDetailsByIdRequest struct {
	CatagoryId int64 `json:"id"`
}

type CategoryDetailsResponse struct {
	CategoryID   int64                 `json:"categoryid"`
	CategoryName string                `json:"categoryname"`
	Description  string                `json:"description"`
	Brands       []BrandDetailRequests `json:"brands"`
}

type BrandDetailRequests struct {
	BrandName  string  `json:"brandname"`
	Price      float64 `json:"price"`
	StockCount int64   `json:"stockcount"`
	Model      string  `json:"model"`
	ImageLink  string  `json:"imagelink"`
}

func (args *CatagoryDetailsByIdRequest) Parse(r *http.Request) error {
	strID := chi.URLParam(r, "id")
	if strID == "" {
		return fmt.Errorf("id parameter is missing or empty")
	}
	intID, err := strconv.Atoi(strID)
	if err != nil {
		return err
	}
	args.CatagoryId = int64(intID)
	return nil
}
