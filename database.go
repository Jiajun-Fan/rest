package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func getDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./dict.db")
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	return &db
}

func createDB(init bool) *gorm.DB {
	db := getDB()
	if init {
		db.CreateTable(&Dict{})
		db.CreateTable(&UserWord{})
		db.CreateTable(&Word{})
		db.CreateTable(&Trans{})
		db.Model(&Word{}).AddUniqueIndex("word", "word")
	}
	return db
}
