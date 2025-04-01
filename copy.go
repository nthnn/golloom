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
	"net/url"
)

// CopyModelResult encapsulates the outcome of a model copying operation,
// including any status messages returned during the process.
type CopyModelResult struct {
	// StatusMessages contains a list of messages detailing the status or
	// progress of the copy operation.
	StatusMessages []string `json:"status_messages"`
}

// CopyModel initiates the copying of a model from a source to a destination within the server.
// It constructs the appropriate API endpoint and sends a POST request with the source and destination parameters.
// Parameters:
//   - ctx: A context.Context object for managing request deadlines and cancellations.
//   - source: The name or identifier of the source model to be copied.
//   - destination: The target name or identifier where the model should be copied to.
//
// Returns:
//   - A pointer to a CopyModelResult struct containing status messages from the operation.
//   - An error if the request fails or the server returns an unexpected response.
func (c *Client) CopyModel(
	ctx context.Context,
	source, destination string,
) (*CopyModelResult, error) {
	rel := &url.URL{Path: "/api/copy"}
	u := c.BaseURL.ResolveReference(rel)

	res, err := c.sendStatusStreamRequest(
		ctx,
		"POST",
		u.String(),
		map[string]string{
			"source":      source,
			"destination": destination,
		},
	)

	if err != nil {
		return nil, err
	}

	return &CopyModelResult{
		StatusMessages: res.StatusMessages,
	}, nil
}
