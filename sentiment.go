/*
Copyright 2015 Aylien, Inc. All Rights Reserved.

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

// SentimentParams is the set of parameters that defines a document whose sentiment needs analysis.
type SentimentParams struct {
	// Either URL or Text is required.
	Text string
	URL  string

	// The analyze mode.
	// Valid option are tweet suitable for short text (default).
	// And document which is more suitable for longer bodies of text.
	Mode string
}

// SentimentResponse is the JSON description of sentiment analysis.
type SentimentResponse struct {
	Text                   string  `json:"text"`
	Polarity               string  `json:"polarity"`
	PolarityConfidence     float32 `json:"polarity_confidence"`
	Subjectivity           string  `json:"subjectivity"`
	SubjectivityConfidence float32 `json:"subjectivity_confidence"`
}

// AspectBasedSentimentParams is the set of parameters that defines a document whose sentiment needs analysis.
type AspectBasedSentimentParams struct {
	// Either URL or Text is required.
	Text string
	URL  string

	// The domain that document belongs to.
	Domain string
}

// Aspect is the JSON description of an aspect.
type Aspect struct {
	Aspect             string  `json:"aspect"`
	AspectConfidence   float32 `json:"aspect_confidence"`
	Polarity           string  `json:"polarity"`
	PolarityConfidence float32 `json:"polarity_confidence"`
}

// AspectBasedSentimentResponse is the JSON description of aspect based sentiment analysis.
type AspectBasedSentimentResponse struct {
	Text      string   `json:"text"`
	Domain    string   `json:"domain"`
	Aspects   []Aspect `json:"aspects"`
	Sentences []struct {
		Text               string   `json:"text"`
		Polarity           string   `json:"polarity"`
		PolarityConfidence float32  `json:"polarity_confidence"`
		Aspects            []Aspect `json:"aspects"`
	} `json:"sentences"`
}

// Sentiment detects the sentiment of the document defined by the given params information.
// It detects the sentiment in terms of polarity (positive, negative or neutral).
// And in terms of subjectivity (subjective or objective).
func (c *Client) Sentiment(params *SentimentParams) (*SentimentResponse, error) {
	body := &url.Values{}

	if len(params.Text) > 0 {
		body.Add("text", params.Text)
	} else if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must either provide url or text")
	}

	if len(params.Mode) > 0 {
		body.Add("mode", params.Mode)
	}

	sentiment := &SentimentResponse{}
	err := c.call("/sentiment", body, sentiment)
	if err != nil {
		return nil, err
	}

	return sentiment, err
}

// AspectBasedSentiment given a review for a product or service, analyzes the sentiment of the review towards each of the aspects of the product or review that are mentioned in it.
func (c *Client) AspectBasedSentiment(params *AspectBasedSentimentParams) (*AspectBasedSentimentResponse, error) {
	body := &url.Values{}

	if len(params.Text) > 0 {
		body.Add("text", params.Text)
	} else if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must either provide url or text")
	}

	if len(params.Domain) == 0 {
		return nil, errors.New("you must specify the domain")
	}

	result := &AspectBasedSentimentResponse{}
	err := c.call("/absa/"+params.Domain, body, result)
	if err != nil {
		return nil, err
	}

	return result, err
}
