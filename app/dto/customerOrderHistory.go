package dto

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SearchByCustomerIdRequest struct {
	UserId int64 `json:"id"`
}

func (args *SearchByCustomerIdRequest) Parse(r *http.Request) error {
	strID := chi.URLParam(r, "id")
	if strID == "" {
		return fmt.Errorf("id parameter is missing or empty")
	}
	intID, err := strconv.Atoi(strID)
	if err != nil {
		return err
	}
	args.UserId = int64(intID)
	return nil
}