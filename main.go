package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

func route(ws *restful.WebService) {
	ws.Route(ws.GET("/users").To(users))
}

func main() {
	ws := new(restful.WebService)
	route(ws)
	restful.Add(ws)
	http.ListenAndServe(":8332", nil)
}
