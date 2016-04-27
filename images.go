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

// ImageTagsParams defines the image whose tags needs to be calculated.
type ImageTagsParams struct {
	URL string
}

// ImageTag is the JSON description of a pair of tag and confidence
type ImageTag struct {
	Tag        string  `json:"tag"`
	Confidence float32 `json:"confidence"`
}

// ImageTagsResponse is the JSON description of tag image response.
type ImageTagsResponse struct {
	Image string     `json:"string"`
	Tags  []ImageTag `json:"tags"`
}

// TagImage tags image defined by the given params information.
func (c *Client) ImageTags(params *ImageTagsParams) (*ImageTagsResponse, error) {
	body := &url.Values{}

	if len(params.URL) > 0 {
		body.Add("url", params.URL)
	} else {
		return nil, errors.New("you must provide a url")
	}

	imageTags := &ImageTagsResponse{}
	err := c.call("image-tags", body, imageTags)
	if err != nil {
		return nil, err
	}

	return imageTags, err
}
