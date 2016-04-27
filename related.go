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

// RelatedParams is the set of parameters that defines a phrase whose related phrases needs to be retrieved.
type RelatedParams struct {
	Count  int
	Phrase string
}

// Related is the JSON description of a related phrase.
type Related struct {
	Phrase   string  `json:"phrase"`
	Distance float32 `json:"distance"`
}

// RelatedResponse is the JSON description of related response.
type RelatedResponse struct {
	Phrase  string    `json:"phrase"`
	Related []Related `json:"related"`
}

// Related returns related phrases to the phrase defined by the given params information.
func (c *Client) Related(params *RelatedParams) (*RelatedResponse, error) {
	body := &url.Values{}

	if len(params.Phrase) > 0 {
		body.Add("phrase", params.Phrase)
	} else {
		return nil, errors.New("you must provide a phrase")
	}

	if params.Count > 0 {
		body.Add("count", strconv.Itoa(params.Count))
	}

	related := &RelatedResponse{}
	err := c.call("/related", body, related)
	if err != nil {
		return nil, err
	}

	return related, err
}
