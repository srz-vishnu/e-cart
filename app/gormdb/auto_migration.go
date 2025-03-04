package gormdb

import (
	"log"

	"e-cart/app/internal"

	"gorm.io/gorm"
)

func Automigration(db *gorm.DB) error {
	if err := db.AutoMigrate(&internal.Userdetail{}); err != nil {
		log.Fatalf("Migration error for user:%v", err)
	}
	if err := db.AutoMigrate(&internal.Productdetail{}); err != nil {
		log.Fatalf("migration failed for Product : %v", err)
	}
	log.Println("Migration success")
	return nil
}
