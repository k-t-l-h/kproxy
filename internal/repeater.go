package internal

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (p *Proxy) SaveRequest() {
}

func (p *Proxy) GetList(w http.ResponseWriter, r *http.Request) {
	//just return all
	requests, err := p.getRequestsList()
	if err != nil {
	}

	data, err := json.Marshal(requests)
	if err != nil {
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (p *Proxy) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids, found := vars["id"]
	if !found {
		return
	}
	id, err := strconv.Atoi(ids)
	request, err := p.getRequest(id)
	if err != nil {
	}
	data, err := json.Marshal(request)
	if err != nil {
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (p *Proxy) RepeateOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids, found := vars["id"]
	if !found {
		return
	}
	id, err := strconv.Atoi(ids)
	request, err := p.getRequest(id)
	if err != nil {
	}

	client := http.DefaultClient
	do, err := client.Do(request)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(do.Body)
	if err != nil {
		return
	}
	w.WriteHeader(do.StatusCode)
	w.Write(body)
}

func (p *Proxy) RepeatSQLInj(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids, found := vars["id"]
	if !found {
		return
	}
	id, err := strconv.Atoi(ids)
	request, err := p.getRequest(id)
	if err != nil {
	}
	data, err := json.Marshal(request)
	if err != nil {
	}

	// default status code and body len
	// code := request.Response.StatusCode
	// def := request.Response.Body

	//add \' or \" to headers
	for _, strings := range request.Header {
		for i := 0; i < len(strings); i++ {
		}
	}

	//add \' or \" to form values
	if request.Form != nil {
		for _, strings := range request.Form {
			for i := 0; i < len(strings); i++ {
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
