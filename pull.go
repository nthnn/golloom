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

// PullModelResult represents the structure of the response received
// after initiating a model pull request, containing status messages.
type PullModelResult struct {
	StatusMessages []string `json:"status_messages"` // List of status messages returned from the model pull operation.
}

// PullModel initiates a request to pull a specified model from the server.
// It constructs the appropriate URL, sends the request, and returns
// the status messages received in response.
func (c *Client) PullModel(
	ctx context.Context,
	model string,
) (*PullModelResult, error) {
	rel := &url.URL{Path: "/api/pull"}
	u := c.BaseURL.ResolveReference(rel)

	res, err := c.sendStatusStreamRequest(
		ctx,
		"POST",
		u.String(),
		map[string]string{
			"model": model,
		},
	)

	if err != nil {
		return nil, err
	}

	return &PullModelResult{
		StatusMessages: res.StatusMessages,
	}, nil
}
