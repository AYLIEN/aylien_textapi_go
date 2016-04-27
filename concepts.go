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

// ConceptsParams is the set of parameters that defines a document whose concepts needs to be extracted.
type ConceptsParams struct {
	// Either URL or Text is required.
	URL  string
	Text string

	// Valid languages are en, de, fr, es, it, pt and auto.
	// Default is en.
	Language string
}

// A SurfaceForm is the JSON description of a concept's surface form.
type SurfaceForm struct {
	String string  `json:"string"`
	Score  float32 `json:"score"`
	Offset int     `json:"offset"`
}

// A Concept is the JSON description of a concept in document.
type Concept struct {
	SurfaceForms []SurfaceForm `json:"surfaceForms"`
	Types        []string      `json:"types"`
	Support      int           `json:"support"`
}

// A ConceptsResponse is the JSON description of concept extraction response.
type ConceptsResponse struct {
	Text     string             `json:"text"`
	Language string             `json:"language"`
	Concepts map[string]Concept `json:"concepts"`
}

// Concepts extracts concepts mentioned in the document defined by the given params information.
func (c *Client) Concepts(params *ConceptsParams) (*ConceptsResponse, error) {
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

	concepts := &ConceptsResponse{}
	err := c.call("/concepts", body, concepts)
	if err != nil {
		return nil, err
	}

	return concepts, err
}
