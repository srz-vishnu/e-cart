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
	if err := db.AutoMigrate(&internal.Category{}); err != nil {
		log.Fatalf("migration failed for Product : %v", err)
	}
	if err := db.AutoMigrate(&internal.Brand{}); err != nil {
		log.Fatalf("migration failed for Product : %v", err)
	}
	if err := db.AutoMigrate(&internal.Cart{}); err != nil {
		log.Fatalf("migration failed for Cart : %v", err)
	}
	if err := db.AutoMigrate(&internal.Order{}); err != nil {
		log.Fatalf("migration failed for order : %v", err)
	}
	if err := db.AutoMigrate(&internal.OrderItem{}); err != nil {
		log.Fatalf("migration failed for order item : %v", err)
	}
	if err := db.AutoMigrate(&internal.UserFavoriteBrand{}); err != nil {
		log.Fatalf("migration failed for favorite brand : %v", err)
	}
	if err := db.AutoMigrate(&internal.ActiveToken{}); err != nil {
		log.Fatalf("migration failed for active token : %v", err)
	}
	log.Println("Migration success")
	return nil
}
