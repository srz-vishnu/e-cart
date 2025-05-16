package dto

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BlockUserRequest struct {
	UserID int64 `json:"userid"`
}

func (args *BlockUserRequest) Parse(r *http.Request) error {
	strID := chi.URLParam(r, "userid")
	intID, err := strconv.Atoi(strID)
	if err != nil {
		return fmt.Errorf("invalid user ID to block")
	}
	args.UserID = int64(intID)

	return nil
}
