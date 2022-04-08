package internal

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"
)

func (p *Proxy) GetRequest(message []byte) (*http.Request, error) {
	r := bytes.NewReader(message)
	reader := bufio.NewReader(r)
	req, err := http.ReadRequest(reader)
	return req, err
}

func (p *Proxy) PutRequest(request *http.Request) ([]byte, error) {
	b, err := httputil.DumpRequest(request, true)
	return b, err
}

func (p *Proxy) PutResponse(response *http.Response) ([]byte, error) {
	b, err := httputil.DumpResponse(response, true)
	return b, err
}
