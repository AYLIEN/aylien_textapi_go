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

/*
ELSA stands for Entity Level Sentiment Analysis,
This endpoint analyzes sentiment towards entities found in text.
Extracts mentions of named entities (Person, Organization, Location), 
associates a type and links them to DBPedia (where possible), and evaluates sentiment towards each of the entities.
This endpoint includes some of the functionality of the entities and concepts endpoints.
*/

// ElsaParams is the set of parameters that defines a document whose entities needs to be extracted.
type ElsaParams struct {
	// Either URL or Text is required.
	URL  string
	Text string
}

// An ElsaResponse is the JSON description of entities extraction response.
type ElsaResponse struct {
	Text     string                 `json:"text"`
	Entities []map[string]interface{} `json:"entities"`
}

// Elsa extracts entities mentioned in the document defined by the given params information.
func (c *Client) Elsa(params *ElsaParams) (*ElsaResponse, error) {
	body := &url.Values{}

	if len(params.Text) > 0 {
		body.Add("text", params.Text)
	} else if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must either provide url or text")
	}

	elsa := &ElsaResponse{}
	err := c.call("/elsa", body, elsa)
	if err != nil {
		return nil, err
	}

	return elsa, err
}
