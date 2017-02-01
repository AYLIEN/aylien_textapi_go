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
	"encoding/json"
	"errors"
	"net/url"
)

type CombinedParams struct {
	URL       string
	Text      string
	Endpoints []string
}

type endpointResult struct {
	Endpoint string      `json:"endpoint"`
	Result   interface{} `json:"result"`
}

type combinedRawResponse struct {
	Text    string           `json:"text"`
	Results []endpointResult `json:"results"`
}

type CombinedResponse struct {
	Text            string
	Article         ExtractResponse
	Summary         SummarizeResponse
	Concepts        ConceptsResponse
	Entities        EntitiesResponse
	Hashtags        HashtagsResponse
	Language        LanguageResponse
	Sentiment       SentimentResponse
	Classifications ClassifyResponse
}

func (c *CombinedResponse) UnmarshalJSON(data []byte) error {
	combinedRaw := &combinedRawResponse{}
	err := json.Unmarshal(data, combinedRaw)
	if err != nil {
		return err
	}
	c.Text = combinedRaw.Text
	for _, r := range combinedRaw.Results {
		o, err := json.Marshal(r.Result)
		if err != nil {
			return err
		}
		switch r.Endpoint {
		case "extract":
			err = json.Unmarshal(o, &c.Article)
		case "language":
			err = json.Unmarshal(o, &c.Language)
		case "entities":
			err = json.Unmarshal(o, &c.Entities)
		case "concepts":
			err = json.Unmarshal(o, &c.Concepts)
		case "classify":
			err = json.Unmarshal(o, &c.Classifications)
		case "hashtags":
			err = json.Unmarshal(o, &c.Hashtags)
		case "sentiment":
			err = json.Unmarshal(o, &c.Sentiment)
		case "summarize":
			err = json.Unmarshal(o, &c.Summary)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Combined(params *CombinedParams) (*CombinedResponse, error) {
	body := &url.Values{}

	if len(params.Text) > 0 {
		body.Add("text", params.Text)
	} else if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must either provide url or text")
	}

	if len(params.Endpoints) < 2 {
		return nil, errors.New("you must provide at least two endpoints")
	}
	for _, p := range params.Endpoints {
		body.Add("endpoint", p)
	}

	response := &CombinedResponse{}
	err := c.call("/combined", body, response)
	if err != nil {
		return nil, err
	}

	return response, err
}
