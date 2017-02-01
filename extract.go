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
	"strings"
	"time"
)

// DateTime embeds time.Time and acts as an alias for custom unmarshalling.
type DateTime struct {
	time.Time
}

// Custom unmarshaller for time.Time which discards errors and will set the time to zero value
// in case parsing failed. Using time.Time.IsZero() is advised.
func (t *DateTime) UnmarshalJSON(buf []byte) error {
	tt, _ := time.Parse(time.RFC3339, strings.Trim(string(buf), `"`))
	t.Time = tt
	return nil
}

// ExtractParams is the set of parameters that defines a web page whose data needs to be extracted.
type ExtractParams struct {
	// Either URL or HTML is required.
	URL  string
	HTML string // Raw HTML of web page

	// Whether to extract the best image of the article.
	BestImage bool
}

// ExtractResponse is the JSON description of extract response.
type ExtractResponse struct {
	Title       string   `json:"title"`
	Article     string   `json:"article"`
	Image       string   `json:"image"`
	Author      string   `json:"author"`
	PublishDate DateTime `json:"publishDate"`
	Videos      []string `json:"videos"`
	Feeds       []string `json:"feeds"`
	Keywords    []string `json:"keywords"`
}

// Extract extracts information from the web page defined by the given params information.
func (c *Client) Extract(params *ExtractParams) (*ExtractResponse, error) {
	body := &url.Values{}

	if len(params.HTML) > 0 {
		body.Add("html", params.HTML)
	} else if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must either provide url or html")
	}

	if params.BestImage {
		body.Add("best_image", "true")
	} else {
		body.Add("best_image", "false")
	}

	article := &ExtractResponse{}
	err := c.call("/extract", body, article)
	if err != nil {
		return nil, err
	}

	return article, err
}
