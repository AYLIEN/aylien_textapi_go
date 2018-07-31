/*
Copyright 2017 Aylien, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package textapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// version is SDK's version.
const version = "0.6.0"

// An Auth is an authentication token that will be used to
// authenticate client.
type Auth struct {
	// ApplicationID and ApplicationKey identify the client using Text API.
	// These fields are always required
	ApplicationID  string
	ApplicationKey string
}

// A RateLimits is the HTTP X-RateLimit-* headers of last response.
type RateLimits struct {
	Limit     int
	Remaining int
	Reset     int
}

// A Client can make calls to the Text API.
type Client struct {
	auth           Auth
	useHTTPS       bool
	apiHostAndPath string

	RateLimits *RateLimits

	Client *http.Client
}

// An Error is the JSON response whenever an error occurs.
type Error struct {
	Message string `json:"error"`
}

// NewClient returns a new client using the given auth information.
// To use HTTPS, pas useHttps = true.
func NewClient(auth Auth, useHTTPS bool) (*Client, error) {
	if len(auth.ApplicationID) == 0 || len(auth.ApplicationKey) == 0 {
		return nil, errors.New("invalid application ID or application key")
	}
	client := &Client{
		auth:           auth,
		useHTTPS:       useHTTPS,
		apiHostAndPath: "api.aylien.com/api/v1",
		RateLimits:     &RateLimits{},
	}

	return client, nil
}

func (c *Client) call(path string, form *url.Values, v interface{}) error {
	req, err := c.newRequest(path, form)
	if err != nil {
		return err
	}

	if err := c.do(req, v); err != nil {
		return err
	}

	return nil
}

func (c *Client) newRequest(path string, form *url.Values) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	protocol := "http"
	if c.useHTTPS {
		protocol = "https"
	}

	url := protocol + "://" + c.apiHostAndPath + path

	var body io.Reader
	if form != nil && len(*form) > 0 {
		data := form.Encode()
		body = bytes.NewBufferString(data)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Add("User-Agent", "Aylien Text API Go "+version)
	req.Header.Add("X-AYLIEN-TextAPI-Application-ID", c.auth.ApplicationID)
	req.Header.Add("X-AYLIEN-TextAPI-Application-Key", c.auth.ApplicationKey)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	client := &http.Client{}
	if c.Client != nil {
		client = c.Client
	}

	res, err := client.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		var e Error
		if err = json.Unmarshal(resBody, &e); err != nil {
			return errors.New(string(resBody))
		}

		return errors.New(e.Message)
	}

	c.RateLimits.Limit, _ = strconv.Atoi(res.Header.Get("X-RateLimit-Limit"))
	c.RateLimits.Reset, _ = strconv.Atoi(res.Header.Get("X-RateLimit-Reset"))
	c.RateLimits.Remaining, _ = strconv.Atoi(res.Header.Get("X-RateLimit-Remaining"))

	if v != nil {
		err = json.Unmarshal(resBody, v)
		if err != nil {
			return errors.New("invalid response")
		}

		return nil
	}

	return nil
}
