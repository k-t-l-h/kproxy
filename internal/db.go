package internal

import "net/http"

const (
	GetAllQuery = ""
	GetOneQuery = ""
)

func (p *Proxy) getRequestsList() ([]http.Request, error) {
	return nil, nil
}

func (p *Proxy) getRequest(id int) (http.Request, error) {
	return http.Request{}, nil
}
