package dist

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	ErrHttpStatusCode = errors.New("http return bad status code.")
)

type HttpClient struct {
	c	*http.Client
	sc int 	// last status code
}

func NewHttpClient() *HttpClient {
	var transport = &http.Transport{
		TLSHandshakeTimeout: 5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	var hc = &HttpClient{
		c: &http.Client{
			Timeout: 30 * time.Second,
			Transport: transport,
		},
	}
	return hc
}

// post json and return json
func (self *HttpClient) PostJson(url string, json []byte) ([]byte, error) {
	resp, err := self.c.Post(url, "application/json", bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	self.sc = resp.StatusCode

	if self.sc == http.StatusOK || self.sc == http.StatusCreated{
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}
	return nil, ErrHttpStatusCode
}

func (self *HttpClient)LastStatusCode() int {
	return self.sc
}