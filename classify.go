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
	"strconv"
)

// ClassifyParams is the set of parameters that defines a document whose classification needs to be calculated.
type ClassifyParams struct {
	// Either URL or Text is required.
	URL  string
	Text string

	// Valid languages are en, de, fr, es, it, pt and auto.
	// Default is en.
	Language string
}

// A Category is the JSON description of a classification category.
type Category struct {
	// IPTC SubjectCode
	Code       string  `json:"code"`
	Label      string  `json:"label"`
	Confidence float32 `json:"confidence"`
}

// A ClassifyResponse is the JSON description of classify response.
type ClassifyResponse struct {
	Text       string     `json:"text"`
	Language   string     `json:"language"`
	Categories []Category `json:"categories"`
}

// UnsupervisedClassifyParams is the set of parameters that defines a document whose needs to be classified.
type UnsupervisedClassifyParams struct {
	// Either URL or Text is required.
	URL  string
	Text string

	// List of classes to classify into
	Classes []string

	// Number of concepts used to measure the semantic similarity between two words.
	NumberOfConcepts int
}

// An UnsupervisedClassifyClass is the JSON description of a class.
type UnsupervisedClassifyClass struct {
	Label string  `json:"label"`
	Score float32 `json:"score"`
}

// An UnsupervisedClassifyResponse is the JSON description of unsupervised classification response.
type UnsupervisedClassifyResponse struct {
	Text    string                      `json:"text"`
	Classes []UnsupervisedClassifyClass `json:"classes"`
}

// A ClassifyByTaxonomyParams is the set of parameters that defines a document whose needs to be classified according to a taxonomy.
type ClassifyByTaxonomyParams struct {
	// Either URL or Text is required.
	URL  string
	Text string

	// Valid languages are en, de, fr, es, it, pt and auto.
	// Default is en.
	Language string

	// Valid taxonomies are iab-qag and iptc-subjectcode.
	Taxonomy string
}

// A ClassifyByTaxonomyResponse is the JSON description of classification by taxonomy response.
type ClassifyByTaxonomyResponse struct {
	Text       string `json:"text"`
	Language   string `json:"language"`
	Taxonomy   string `json:"taxonomy"`
	Categories []struct {
		Id        string  `json:"id"`
		Label     string  `json:"label"`
		Score     float32 `json:"score"`
		Confident bool    `json:"confident"`
		Links     []struct {
			Link string `json:"link"`
			Rel  string `json:"rel"`
		} `json:"links"`
	} `json:"categories"`
}

// Classify classifies the document defined by the given params information.
func (c *Client) Classify(params *ClassifyParams) (*ClassifyResponse, error) {
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

	classification := &ClassifyResponse{}
	err := c.call("/classify", body, classification)
	if err != nil {
		return nil, err
	}

	return classification, err
}

// UnsupervisedClassify picks the most semantically relevant class label or tag for the document defined by the given params information.
func (c *Client) UnsupervisedClassify(params *UnsupervisedClassifyParams) (*UnsupervisedClassifyResponse, error) {
	body := &url.Values{}

	if len(params.Text) > 0 {
		body.Add("text", params.Text)
	} else if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must either provide url or text")
	}

	if params.NumberOfConcepts > 0 {
		body.Add("number_of_concepts", strconv.Itoa(params.NumberOfConcepts))
	}

	if len(params.Classes) < 2 {
		return nil, errors.New("you must provide at least two classes")
	}
	for _, c := range params.Classes {
		body.Add("class", c)
	}

	classes := &UnsupervisedClassifyResponse{}
	err := c.call("/classify/unsupervised", body, classes)
	if err != nil {
		return nil, err
	}

	return classes, err
}

// ClassifyByTaxonomy classifies the document defined by the given params information according to the specified taxonomy.
func (c *Client) ClassifyByTaxonomy(params *ClassifyByTaxonomyParams) (*ClassifyByTaxonomyResponse, error) {

	body := &url.Values{}

	if len(params.Text) > 0 {
		body.Add("text", params.Text)
	} else if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must either provide url or text")
	}

	if len(params.Taxonomy) == 0 {
		return nil, errors.New("you must specify the taxonomy")
	}

	if len(params.Language) > 0 {
		body.Add("language", params.Language)
	}

	classifications := &ClassifyByTaxonomyResponse{}
	err := c.call("/classify/"+params.Taxonomy, body, classifications)
	if err != nil {
		return nil, err
	}

	return classifications, err
}
