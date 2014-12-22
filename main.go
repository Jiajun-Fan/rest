package main

import (
	"log"
	"net/http"
)

func main() {

	//db := createDB()
	//db.CreateTable(&Dict{})

	u := DictService{createDB()}
	u.Register()

	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
