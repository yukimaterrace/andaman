package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPStatusError is http status error
type HTTPStatusError struct {
	Code    int
	Message string
}

func (err HTTPStatusError) Error() string {
	return fmt.Sprintf("code:%v message:%s", err.Code, err.Message)
}

type client struct {
	host         string
	commonHeader http.Header
	client       *http.Client
}

func newClient(host string, commonHeader http.Header) *client {
	return &client{
		host:         host,
		commonHeader: commonHeader,
		client:       &http.Client{},
	}
}

func (client *client) request(method string, path string, query url.Values, requestBody interface{}, response interface{}) error {
	rawQuery := ""
	if query != nil {
		rawQuery = query.Encode()
	}

	url := url.URL{
		Scheme:   "https",
		Host:     client.host,
		Path:     path,
		RawQuery: rawQuery,
	}

	var body *bytes.Buffer
	if requestBody != nil {
		buf, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(buf)
	}

	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url.String(), body)
	} else {
		req, err = http.NewRequest(method, url.String(), nil)
	}

	if err != nil {
		return err
	}

	req.Header = client.commonHeader

	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	code := resp.StatusCode
	if code < 200 || code >= 300 {
		return &HTTPStatusError{
			Code:    code,
			Message: string(content),
		}
	}

	if err := json.Unmarshal(content, response); err != nil {
		return err
	}

	return nil
}

func (client *client) get(path string, query url.Values, response interface{}) error {
	return client.request("GET", path, query, nil, response)
}

func (client *client) post(path string, requestBody interface{}, response interface{}) error {
	return client.request("POST", path, nil, requestBody, response)
}

func (client *client) put(path string, requestBody interface{}, response interface{}) error {
	return client.request("PUT", path, nil, requestBody, response)
}
