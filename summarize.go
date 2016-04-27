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
	"strconv"
)

// SummarizeParams is the set of parameters that defines a document whose needs to be summarized.
type SummarizeParams struct {
	// Either URL or a pair of Title and Text is required.
	URL   string
	Text  string
	Title string

	// Summarize mode
	// Valid options are default, and short.
	// short mode produces relatively shorter sentences.
	Mode string

	// Quantity of sentences in default mode.
	// Not applicable to short mode.
	NumberOfSentences     int
	PercentageOfSentences int
}

// SummarizeResponse is the JSON description of summarize response.
type SummarizeResponse struct {
	Text      string   `json:"text"`
	Sentences []string `json:"sentences"`
}

// Summarize summarizes the document defined by the given params information.
func (c *Client) Summarize(params *SummarizeParams) (*SummarizeResponse, error) {
	body := &url.Values{}

	if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else if len(params.Title) > 0 && len(params.Text) > 0 {
		body.Add("title", params.Title)
		body.Add("text", params.Text)
	} else {
		return nil, errors.New("you must either provide url or a pair of text and title")
	}

	if len(params.Mode) > 0 {
		body.Add("mode", params.Mode)
	} else {
		body.Add("mode", "default")
	}

	if params.NumberOfSentences > 0 {
		body.Add("sentences_number", strconv.Itoa(params.NumberOfSentences))
	}

	if params.PercentageOfSentences > 0 {
		body.Add("sentences_percentage", strconv.Itoa(params.PercentageOfSentences))
	}

	summary := &SummarizeResponse{}
	err := c.call("/summarize", body, summary)
	if err != nil {
		return nil, err
	}

	return summary, err
}
