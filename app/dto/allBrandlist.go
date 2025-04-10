package dto

type BrandDetailResponse struct {
	BrandName    string  `json:"brandname"`
	BrandId      int64   `json:"brandid"`
	Price        float64 `json:"price" `
	StockCount   int64   `json:"stockcount"`
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"categoryname"`
}
