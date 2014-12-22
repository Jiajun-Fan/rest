package main

import (
	"github.com/emicklei/go-restful"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
)

type DictService struct {
	db *gorm.DB
}

func (u DictService) Register() {
	ws := new(restful.WebService)
	ws.
		Path("/dict")
	//Path("/dict").
	//Consumes(restful.MIME_XML, restful.MIME_JSON).
	//Produces(restful.MIME_XML, restful.MIME_JSON)

	ws.Route(ws.GET("/{dict-name}").To(u.findDict).
		Doc("find a dictionary").
		Operation("findDict").
		Param(ws.PathParameter("dict-name", "identifier of the user").DataType("string")).
		Reads(Dict{}))

	ws.Route(ws.PUT("").To(u.createDict).
		Doc("create a dictionary").
		Operation("createDict").
		Reads(Dict{}))

	restful.Add(ws)
}

func (u DictService) findDict(request *restful.Request, response *restful.Response) {
	io.WriteString(response, "world")
}

func (u DictService) createDict(request *restful.Request, response *restful.Response) {
	dict := new(Dict)
	err := request.ReadEntity(dict)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	u.db.Create(dict)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(dict)
}
