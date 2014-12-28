package main

import (
	"log"
	"net/http"
)

func main() {

	db := createDB(true)

	u := DictService{db}
	u.Register()

	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
