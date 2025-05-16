package dto

type FavoriteBrandResponse struct {
	BrandID   int64   `json:"brand_id"`
	BrandName string  `json:"brand_name"`
	Price     float64 `json:"price"`
	Stock     int64   `json:"stock"`
	ImageLink string  `json:"Image_link"`
}
