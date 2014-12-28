package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"net/http"
)

type DictService struct {
	db *gorm.DB
}

func err500(err error, response *restful.Response) {
	response.AddHeader("Content-Type", "text/plain")
	response.WriteErrorString(http.StatusInternalServerError, err.Error())
}

func err400(err error, response *restful.Response) {
	response.AddHeader("Content-Type", "text/plain")
	response.WriteErrorString(http.StatusBadRequest, err.Error())
}

func (u DictService) Register() {
	ws := new(restful.WebService)
	ws.
		Path("/dict").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/list").To(u.listDict).
		Doc("list dictionaries").
		Operation("listDict").
		Writes(Dict{}))

	ws.Route(ws.GET("/find/{dict-name}").To(u.findDict).
		Doc("find a dictionary").
		Operation("findDict").
		Param(ws.PathParameter("name", "dictionary name").DataType("string")).
		Reads(Dict{}))

	ws.Route(ws.PUT("/create").To(u.createDict).
		Doc("create a dictionary").
		Operation("createDict").
		Reads(Dict{}))

	ws.Route(ws.POST("/trans").To(u.translate).
		Doc("translate and add to dictionary").
		Operation("translate").
		Writes(Word{}))

	restful.Add(ws)
}

func (u DictService) translate(request *restful.Request, response *restful.Response) {
	word := new(UserWord)
	db := getDB()
	err_input := request.ReadEntity(&word)
	if err_input != nil {
		err500(err_input, response)
		return
	}

	dict := new(Dict)
	dict.Id = -1
	if word.DictId != 0 {
		db.Where("id = ?", word.DictId).Find(dict)
	}
	if dict.Id == -1 {
		err400(errors.New("bad dictionary id"), response)
		return
	}
	fmt.Printf("%+v", dict)

	url := "http://openapi.baidu.com/public/2.0/bmt/translate?client_id=b59swMowBKPkg98uiQnKqsAi&from=auto&to=auto&q=" + word.Word
	fmt.Printf("%s\n", url)
	r, err_baidu := http.Get(url)
	if err_baidu != nil {
		err500(err_baidu, response)
		return
	}
	defer r.Body.Close()
	body, err_content := ioutil.ReadAll(r.Body)
	if err_content != nil {
		err500(err_content, response)
		return
	}

	var transResult TransResult
	err_json := json.Unmarshal(body, &transResult)
	if err_json != nil {
		err500(err_json, response)
		return
	}
	fmt.Printf("%+v\n", transResult)
	response.WriteEntity(word)
}

func (u DictService) findDict(request *restful.Request, response *restful.Response) {
	io.WriteString(response, "world")
}

func (u DictService) listDict(request *restful.Request, response *restful.Response) {
	dicts := make([]Dict, 0, 25)
	u.db.Limit(25).Find(&dicts)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(dicts)
}

func (u DictService) createDict(request *restful.Request, response *restful.Response) {
	dict := new(Dict)
	err := request.ReadEntity(&dict)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	u.db.Create(dict)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(dict)
}

type TransService struct {
	db *gorm.DB
}
