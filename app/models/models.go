package models // or package product if you placed it in internal/product

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// StringArray is a custom type to handle []string as JSON in DB
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray: value is not []byte")
	}
	return json.Unmarshal(bytes, s)
}
