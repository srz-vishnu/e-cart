package dto

type GetUserDetailsResponse  struct {
	UserName string `json:"username"`
	Mail     string `json:"mail"`
	Address  string `json:"address"`
	Pincode  int64  `json:"pincode"`
	Phone    int64  `json:"phonenumber"`
}

// func (args *UpdateUserDetailRequest) Parse(r *http.Request) error {
// 	strID := chi.URLParam(r, "userid")
// 	intID, err := strconv.Atoi(strID)
// 	if err != nil {
// 		return err
// 	}
// 	args.UserID = int64(intID)

// 	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
// 		return err
// 	}

// 	return nil
// }
