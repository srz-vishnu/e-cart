package dto

type UserCartItems struct {
	UserID    int64 `json:"userid"`
	ProductID int64 `json:"productid"`
	Quantity  int64 `json:"quantity"`
}
