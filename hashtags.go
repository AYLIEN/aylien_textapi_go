/*
Copyright 2016 Aylien, Inc. All Rights Reserved.

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
	"errors"
	"net/url"
)

// HashtagsParams is the set of parameters that defines a document whose hashtags need to be calculated.
type HashtagsParams struct {
	// Either URL or Text is required.
	URL  string
	Text string

	// Valid languages are en, de, fr, es, it, pt and auto.
	// Default is en.
	Language string
}

// HashtagsResponse is the JSON description of a hashtags in a document.
type HashtagsResponse struct {
	Text     string   `json:"text"`
	Language string   `json:"language"`
	Hashtags []string `json:"hashtags"`
}

// Hashtags calculates best hashtags describing the document defined by the given params information.
func (c *Client) Hashtags(params *HashtagsParams) (*HashtagsResponse, error) {
	body := &url.Values{}

	if len(params.Text) > 0 {
		body.Add("text", params.Text)
	} else if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must either provide url or text")
	}

	if len(params.Language) > 0 {
		body.Add("language", params.Language)
	}

	hashtags := &HashtagsResponse{}
	err := c.call("/hashtags", body, hashtags)
	if err != nil {
		return nil, err
	}

	return hashtags, err
}
