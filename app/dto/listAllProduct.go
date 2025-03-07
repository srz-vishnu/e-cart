package dto

type CatagoryListResponse struct {
	CatagoryID   int64  `json:"catagoryid"`
	CatagoryName string `json:"catagoryname"`
	Description  string `json:"description"`
}
