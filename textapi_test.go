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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	auth          Auth
	client        *Client
	testServer    *httptest.Server
	testServerURL string
)

func init() {
	auth = Auth{"test", "test"}
	client, _ = NewClient(auth, false)
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Aylien-Textapi-Application-Key") != "test" ||
			r.Header.Get("X-Aylien-Textapi-Application-Id") != "test" {
			w.WriteHeader(403)
			fmt.Fprintln(w, "Authentication failed")
		} else {
			uri := r.RequestURI
			r.ParseForm()
			var bytes []byte
			switch uri {
			case "/extract":
				bytes, _ = json.Marshal(ExtractResponse{})
			case "/concepts":
				bytes, _ = json.Marshal(ConceptsResponse{})
			case "/classify":
				bytes, _ = json.Marshal(ClassifyResponse{})
			case "/classify/iab-qag":
				bytes, _ = json.Marshal(ClassifyByTaxonomyResponse{})
			case "/entities":
				url := r.FormValue("url")
				if url == "invalid" {
					w.WriteHeader(400)
					bytes, _ = json.Marshal(Error{Message: "requirement failed: provided url is not valid."})
				} else {
					bytes, _ = json.Marshal(EntitiesResponse{})
				}
			case "/hashtags":
				bytes, _ = json.Marshal(HashtagsResponse{})
			case "/sentiment":
				bytes, _ = json.Marshal(SentimentResponse{})
			case "/language":
				w.Header().Add("X-RateLimit-Limit", "1000")
				w.Header().Add("X-RateLimit-Reset", "1420479141")
				w.Header().Add("X-RateLimit-Remaining", "999")
				bytes, _ = json.Marshal(LanguageResponse{})
			case "/summarize":
				url := r.FormValue("url")
				if url == "invalid" {
					w.WriteHeader(400)
					bytes, _ = json.Marshal(Error{Message: "requirement failed: provided url is not valid."})
				} else {
					bytes, _ = json.Marshal(SummarizeResponse{})
				}
			case "/image-tags":
				bytes, _ = json.Marshal(ImageTagsResponse{})
			case "/absa/cars":
				bytes, _ = json.Marshal(AspectBasedSentimentResponse{})
			}
			fmt.Fprintln(w, string(bytes))
		}
	}))
	testServerURL = strings.TrimLeft(testServer.URL, "http://")
	client.apiHostAndPath = testServerURL
}

func TestClientCreation(t *testing.T) {
	auth := Auth{"random", "random"}
	client, _ := NewClient(auth, true)
	if !client.useHTTPS {
		t.Error("must be using HTTPS")
	}
}

func TestInvalidKeys(t *testing.T) {
	badAuth := Auth{"wrongtest", "wrongtest"}
	badClient, _ := NewClient(badAuth, false)
	badClient.apiHostAndPath = testServerURL
	params := &LanguageParams{Text: "Hello"}
	if _, err := badClient.Language(params); err == nil {
		t.Error("did not return error")
	}
}

func TestSentiment(t *testing.T) {
	params := &SentimentParams{}
	_, err := client.Sentiment(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.Text = "John is a very good football player!"
	_, err = client.Sentiment(params)
	if err != nil {
		t.Error(err)
	}
}

func TestLanguage(t *testing.T) {
	params := &LanguageParams{}
	_, err := client.Language(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.Text = "John is a very good football player!"
	_, err = client.Language(params)
	if err != nil {
		t.Error(err)
	}
	if client.RateLimits.Limit != 1000 || client.RateLimits.Remaining != 999 || client.RateLimits.Reset != 1420479141 {
		t.Error("invalid ratelimits")
	}
}

func TestExtract(t *testing.T) {
	params := &ExtractParams{}
	_, err := client.Extract(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.URL = "http://example.com/"
	_, err = client.Extract(params)
	if err != nil {
		t.Error(err)
	}
}

func TestClassify(t *testing.T) {
	params := &ClassifyParams{}
	_, err := client.Classify(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.Text = "Just a random piece of text"
	_, err = client.Classify(params)
	if err != nil {
		t.Error(err)
	}
}

func TestConcepts(t *testing.T) {
	params := &ConceptsParams{}
	_, err := client.Concepts(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.Text = "Another piece of random text"
	_, err = client.Concepts(params)
	if err != nil {
		t.Error(err)
	}
}

func TestEntities(t *testing.T) {
	params := &EntitiesParams{URL: "invalid"}
	_, err := client.Entities(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.Text = "Valid text"
	_, err = client.Entities(params)
	if err != nil {
		t.Error(err)
	}
}

func TestHashtags(t *testing.T) {
	params := &HashtagsParams{}
	_, err := client.Hashtags(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.Text = "Google SDK"
	_, err = client.Hashtags(params)
	if err != nil {
		t.Error(err)
	}
}

func TestSummarize(t *testing.T) {
	params := &SummarizeParams{}
	_, err := client.Summarize(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.URL = "invalid"
	_, err = client.Summarize(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.URL = ""
	params.Title = "title"
	_, err = client.Summarize(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.Text = "text"
	_, err = client.Summarize(params)
	if err != nil {
		t.Error(err)
	}
}

func TestImageTags(t *testing.T) {
	params := &ImageTagsParams{}
	_, err := client.ImageTags(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.URL = "https://developer.aylien.com/images/logo-small.png"
	_, err = client.ImageTags(params)
	if err != nil {
		t.Error(err)
	}
}

func TestClassifyByTaxonomy(t *testing.T) {
	params := &ClassifyByTaxonomyParams{}
	_, err := client.ClassifyByTaxonomy(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.URL = "http://techcrunch.com/2015/07/16/microsoft-will-never-give-up-on-mobile"
	_, err = client.ClassifyByTaxonomy(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.Taxonomy = "iab-qag"
	_, err = client.ClassifyByTaxonomy(params)
	if err != nil {
		t.Error(err)
	}
}

func TestAspectBasedSentiment(t *testing.T) {
	params := &AspectBasedSentimentParams{}
	_, err := client.AspectBasedSentiment(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.Text = "Delicious food. Disappointing service."
	_, err = client.AspectBasedSentiment(params)
	if err == nil {
		t.Error("did not return error")
	}
	params.Domain = "cars"
	_, err = client.AspectBasedSentiment(params)
	if err != nil {
		t.Error(err)
	}
}
