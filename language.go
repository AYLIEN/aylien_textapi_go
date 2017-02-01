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
	"errors"
	"net/url"
)

// LanguageParams is the set of parameters that defines a document whose language needs to be calculated.
type LanguageParams struct {
	// Either URL or Text is required.
	URL  string
	Text string
}

// LanguageResponse is the JSON description of language response.
type LanguageResponse struct {
	Text       string  `json:"text"`
	Language   string  `json:"lang"`
	Confidence float32 `json:"confidence"`
}

// Language calculates the language in which the document defined by the given params information is written in.
func (c *Client) Language(params *LanguageParams) (*LanguageResponse, error) {
	body := &url.Values{}

	if len(params.Text) > 0 {
		body.Add("text", params.Text)
	} else if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must either provide url or text")
	}

	language := &LanguageResponse{}
	err := c.call("/language", body, language)
	if err != nil {
		return nil, err
	}

	return language, err
}
