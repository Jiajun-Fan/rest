package main

import (
	"encoding/json"
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type HttpException struct {
	err         error
	resp        *restful.Response
	status_code int
}

func try(e *HttpException) {
	if e.err != nil {
		panic(e)
	}
}

func except() {
	e := recover()
	if e == nil {
		return
	}
	if exp, ok := e.(*HttpException); ok {
		exp.resp.AddHeader("Content-Type", "text/plain")
		exp.resp.WriteErrorString(exp.status_code, exp.err.Error())
	} else {
		panic(e)
	}
}

type DictService struct {
	db *gorm.DB
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
		Writes(RealWord{}))

	restful.Add(ws)
}

func (u DictService) translate(request *restful.Request, response *restful.Response) {
	rword := new(RealWord)

	uword := new(UserWord)
	err := request.ReadEntity(&uword)

	defer except()

	try(&HttpException{err, response, http.StatusBadRequest})

	if uword.DictId == 0 || uword.Word == "" {
		try(&HttpException{errors.New("bad request"), response, http.StatusBadRequest})
	}

	dict := new(Dict)
	u.db.Where("id = ?", uword.DictId).Find(dict)
	if u.db.NewRecord(dict) {
		try(&HttpException{errors.New("bad dictionary id"), response, http.StatusBadRequest})
	}

	u.db.Where("word = ? AND dict_id = ?", uword.Word, uword.DictId).Find(uword)
	if !u.db.NewRecord(uword) {
		u.db.Where("id = ?", uword.WordId).Find(rword)
		if u.db.NewRecord(rword) {
			try(&HttpException{errors.New("Internal error"), response, http.StatusInternalServerError})
		} else {
			u.db.Model(rword).Related(&rword.Trans)
		}
		response.WriteEntity(rword)
		return
	}

	u.db.Where("word = ?", uword.Word).Find(rword)

	if u.db.NewRecord(rword) {

		rword.Word = uword.Word

		url := "http://openapi.baidu.com/public/2.0/translate/dict/simple?client_id=b59swMowBKPkg98uiQnKqsAi&from=auto&to=auto&q=" + rword.Word
		r, err := http.Get(url)
		try(&HttpException{err, response, http.StatusInternalServerError})

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		try(&HttpException{err, response, http.StatusInternalServerError})

		var transResult TransResult
		err = json.Unmarshal(body, &transResult)
		try(&HttpException{err, response, http.StatusInternalServerError})

		if transResult.Error != 0 {
			try(&HttpException{errors.New("can't translate"), response, http.StatusInternalServerError})
		}

		for i := range transResult.Data.Symbols {
			symbol := transResult.Data.Symbols[i]
			for j := range symbol.Parts {
				part := symbol.Parts[j]
				for k := range part.Means {
					var t Trans
					t.Trans = part.Means[k]
					rword.Trans = append(rword.Trans, t)
				}
			}
		}

		u.db.Create(rword)
	}

	uword.WordId = rword.Id
	u.db.Create(uword)

	response.WriteEntity(rword)
}

func (u DictService) findDict(request *restful.Request, response *restful.Response) {
	dict := new(Dict)
	err := request.ReadEntity(&dict)
	defer except()
	try(&HttpException{err, response, http.StatusBadRequest})
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
