package dto

type BrandDetailResponse struct {
	BrandName  string  `json:"brandname" validate:"required"`
	Price      float64 `json:"price" validate:"required"`
	StockCount int64   `json:"stockcount" validate:"required"`
}
