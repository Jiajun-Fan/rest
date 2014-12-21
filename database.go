package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func createDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./aalist.db")
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	return &db
}
