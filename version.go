/*
 * Copyright 2025 Nathanne Isip
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package golloom

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// Version represents the structure of the version information
// returned by the API, including the version string and optional build time.
type Version struct {
	Version   string `json:"version"`              // The version string of the application.
	BuildTime string `json:"build_time,omitempty"` // The build time of the application, optional field.
}

// Version retrieves the current version information of the application
// by sending a GET request to the /api/version endpoint.
// It returns a Version struct containing the version details.
func (c *Client) Version(ctx context.Context) (*Version, error) {
	rel := &url.URL{
		Path: "/api/version",
	}

	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var verResp Version
	if err := json.NewDecoder(resp.Body).Decode(&verResp); err != nil {
		return nil, err
	}

	return &verResp, nil
}
