package dto

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SearchByCatagoryIdRequest struct {
	CatagoryId int64 `json:"id"`
}

func (args *SearchByCatagoryIdRequest) Parse(r *http.Request) error {
	strID := chi.URLParam(r, "id")
	log.Printf("Extracted ID from URL issssssssssssssssss: '%s'\n", strID)
	if strID == "" {
		return fmt.Errorf("id parameter is missing or empty")
	}
	intID, err := strconv.Atoi(strID)
	if err != nil {
		return err
	}
	args.CatagoryId = int64(intID)
	return nil
}
