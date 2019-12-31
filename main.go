package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Fruit struct {
	gorm.Model
	Name  string
	Price uint
}

func main() {
	db, err := gorm.Open("sqlite3", "/home/fahad/Emp.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Fruit{})

	// Create
	db.Create(&Fruit{Name: "Grapes", Price: 750})

	// Read
	//var product Product
	//db.First(&product, 1)                   // find product with id 1
	//db.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	//db.Model(&product).Update("Price", 3000)

	// Delete - delete product
	//db.Delete(&Fruit)
}
