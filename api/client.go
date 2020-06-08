package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPStatusError is an error for http status
type HTTPStatusError struct {
	code    int
	message string
}

func (err HTTPStatusError) Error() string {
	return fmt.Sprintf("code:%v message:%s", err.code, err.message)
}

// Client is an API client
type Client struct {
	host         string
	commonHeader http.Header
	client       *http.Client
}

// NewClient is a constructor of API client
func NewClient(host string, commonHeader http.Header) *Client {
	return &Client{
		host:         host,
		commonHeader: commonHeader,
		client:       &http.Client{},
	}
}

func (client *Client) request(method string, path string, query url.Values, bodyMap map[string]interface{}, response interface{}) error {
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
	if bodyMap != nil {
		buf, err := json.Marshal(bodyMap)
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
			code:    code,
			message: string(content),
		}
	}

	if err := json.Unmarshal(content, response); err != nil {
		return err
	}

	return nil
}

// Get is a method for get request
func (client *Client) Get(path string, query url.Values, response interface{}) error {
	return client.request("GET", path, query, nil, response)
}

// Post is a method for post request
func (client *Client) Post(path string, bodyMap map[string]interface{}, response interface{}) error {
	return client.request("POST", path, nil, bodyMap, response)
}

// Put is a method for put request
func (client *Client) Put(path string, bodyMap map[string]interface{}, response interface{}) error {
	return client.request("PUT", path, nil, bodyMap, response)
}
