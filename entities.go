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

// EntitiesParams is the set of parameters that defines a document whose entities needs to be extracted.
type EntitiesParams struct {
	// Either URL or Text is required.
	URL  string
	Text string
}

// A EntitiesResponse is the JSON description of entities extraction response.
type EntitiesResponse struct {
	Text     string              `json:"text"`
	Entities map[string][]string `json:"entities"`
}

// Entities extracts entities mentioned in the document defined by the given params information.
func (c *Client) Entities(params *EntitiesParams) (*EntitiesResponse, error) {
	body := &url.Values{}

	if len(params.Text) > 0 {
		body.Add("text", params.Text)
	} else if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must either provide url or text")
	}

	entities := &EntitiesResponse{}
	err := c.call("/entities", body, entities)
	if err != nil {
		return nil, err
	}

	return entities, err
}
