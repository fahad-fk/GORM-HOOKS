package main

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Error checking intentionally skipped at some places
func TestCleanup(t *testing.T) {
	db, _ := gorm.Open("sqlite3", "/home/fahad/Emp.db")
	defer db.Close()
	// Setup the cleaner
	cleaner := DeleteCreatedEntities(db)
	defer cleaner()

	// Create
	db.Create(&Fruit{Name: "Banana", Price: 100})
	db.Create(&Fruit{Name: "Amrood", Price: 2000})
	db.Create(&Fruit{Name: "Orange", Price: 400})

	// Perform some tests here
	// ...
	// ...
}

func DeleteCreatedEntities(db *gorm.DB) func() {
	type entity struct {
		table   string
		keyname string
		key     interface{}
	}
	var entries []entity
	hookName := "cleanupHook"

	// Setup the onCreate Hook
	db.Callback().Create().After("gorm:create").Register(hookName, func(scope *gorm.Scope) {
		fmt.Printf("Inserted entities of %s with %s=%v\n", scope.TableName(), scope.PrimaryKey(), scope.PrimaryKeyValue())
		entries = append(entries, entity{table: scope.TableName(), keyname: scope.PrimaryKey(), key: scope.PrimaryKeyValue()})
	})
	return func() {
		// Remove the hook once we're done
		defer db.Callback().Create().Remove(hookName)
		// Find out if the current db object is already a transaction
		_, inTransaction := db.CommonDB().(*sql.Tx)
		tx := db
		if !inTransaction {
			tx = db.Begin()
		}
		// Loop from the end. It is important that we delete the entries in the
		// reverse order of their insertion
		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf("Deleting entities from '%s' table with key %v\n", entry.table, entry.key)
			tx.Table(entry.table).Where(entry.keyname+" = ?", entry.key).Delete("")
		}

		if !inTransaction {
			tx.Commit()
		}
	}
}
