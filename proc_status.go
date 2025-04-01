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

// ModelProcessStatus represents the JSON structure returned by the server,
// containing a list of model names currently being processed.
type ModelProcessStatus struct {
	Models []string `json:"models"` // A slice of strings, each representing a model name.
}

// ProcessStatus retrieves the current processing status of models from the server.
// It constructs the appropriate API endpoint, sends a GET request, and decodes the JSON response into a ModelProcessStatus.
//
// Parameters:
//   - ctx: The context.Context for managing request deadlines and cancellations.
//
// Returns:
//   - A pointer to a ModelProcessStatus containing the list of models being processed.
//   - An error if the request fails or the response cannot be processed.
func (c *Client) ProcessStatus(
	ctx context.Context,
) (*ModelProcessStatus, error) {
	rel := &url.URL{Path: "/api/ps"}
	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		u.String(),
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var psResp ModelProcessStatus
	err = json.NewDecoder(resp.Body).Decode(&psResp)
	if err != nil {
		return nil, err
	}

	return &psResp, nil
}
