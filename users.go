package main

import (
	"github.com/emicklei/go-restful"
	"io"
)

func users(req *restful.Request, resp *restful.Response) {
	db := getDB()
	user := new(User)
	db.First(user)
	io.WriteString(resp, user.NickName)
}
