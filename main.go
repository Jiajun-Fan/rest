package main

import (
	"log"
	"net/http"
)

const debug debugging = true

type debugging bool

func (d debugging) Printf(fmt string, args ...interface{}) {
	if d {
		log.Printf(fmt, args...)
	}
}

func main() {

	db := createDB(true)
	//db := createDB(false)

	u := DictService{db}
	u.Register()

	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
