package dto

type AllUserDetails struct {
	//UserID   int64  `json:"userid"`
	UserName string `json:"username"`
	Mail     string `json:"mail"`
	Address  string `json:"address"`
	Phone    int64  `json:"phonenumber"`
}
