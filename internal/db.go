package internal

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
)

const (
	GetAllQuery = ""
	GetOneQuery = "SELECT request FROM requests WHERE id=$1;"
	SaveOneQuery = "INSERT INTO requests(request) VALUES ($1);"
)

func (p *Proxy) getRequestsList() ([]http.Request, error) {
	return nil, nil
}

func (p *Proxy) getRequest(id int) (*http.Request, error) {
	var message string
	row := p.Pool.QueryRow(context.Background(), GetOneQuery, id)
	row.Scan(&message)
	r := bytes.NewReader([]byte(message))
	reader := bufio.NewReader(r)
	req, _ := http.ReadRequest(reader)
	return req, nil
}

func (p *Proxy) writeRequest(message string)  {
	_, err := p.Pool.Exec(context.Background(), SaveOneQuery, message)
	if err != nil {
		return
	}
}