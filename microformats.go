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

// MicroformatsParams is the set of parameters that defines a document whose microformats needs to be extracted.
type MicroformatsParams struct {
	URL string
}

// An Address is the JSON description of an hCard adr
type Address struct {
	ID            string `json:"id"`
	StreetAddress string `json:"streetAddress"`
	Locality      string `json:"locality"`
	Region        string `json:"region"`
	CountryName   string `json:"countryName"`
	PostalCode    string `json:"postalCode"`
}

// A Name is the JSON description of an hCard n
type Name struct {
	ID              string `json:"id"`
	HonorificPrefix string `json:"honorificPrefix"`
	GivenName       string `json:"givenName"`
	AdditionalName  string `json:"additionalName"`
	FamilyName      string `json:"familyName"`
	HonorificSuffix string `json:"honorificSuffix"`
}

// A Location is the JSON description of http://microformats.org/wiki/geo
type Location struct {
	ID        string `json:"id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// An HCard is the JSON description of http://microformats.org/wiki/hcard
type HCard struct {
	ID              string   `json:"id"`
	FullName        string   `json:"fullName"`
	StructuredName  Name     `json:"structuredName"`
	NickName        string   `json:"nickName"`
	Email           string   `json:"email"`
	Photo           string   `json:"photo"`
	URL             string   `json:"url"`
	TelephoneNumber string   `json:"telephoneNumber"`
	Birthday        string   `json:"birthday"`
	Category        string   `json:"category"`
	Note            string   `json:"note"`
	Logo            string   `json:"logo"`
	Location        Location `json:"location"`
	Address         Address  `json:"address"`
	Organization    string   `json:"organization"`
}

// A MicroformatsResponse is the JSON description of microformats extraction response.
type MicroformatsResponse struct {
	HCards []HCard `json:"hCards"`
}

// Microformats extracts microformats from document defined by the given params information.
func (c *Client) Microformats(params *MicroformatsParams) (*MicroformatsResponse, error) {
	body := &url.Values{}

	if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must provide a url")
	}

	microformats := &MicroformatsResponse{}
	err := c.call("/microformats", body, microformats)
	if err != nil {
		return nil, err
	}

	return microformats, err
}
