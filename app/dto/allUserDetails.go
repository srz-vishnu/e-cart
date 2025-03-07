package dto

type AllUserDetails struct {
	//UserID   int64  `json:"userid"`
	UserName string `json:"username" validate:"required"`
	Mail     string `json:"mail" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Phone    int64  `json:"phonenumber" validate:"required"`
}
