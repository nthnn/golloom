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

// PushModelResult represents the structure of the response received
// after initiating a model push request, containing status messages.
type PushModelResult struct {
	StatusMessages []string `json:"status_messages"` // List of status messages returned from the model push operation.
}

// PushModel initiates a request to push a specified model to the server.
// It constructs the appropriate URL, sends the request, and returns
// the status messages received in response.
func (c *Client) PushModel(
	ctx context.Context,
	model string,
) (*PushModelResult, error) {
	rel := &url.URL{Path: "/api/push"}
	u := c.BaseURL.ResolveReference(rel)

	res, err := c.sendStatusStreamRequest(
		ctx,
		"POST",
		u.String(),
		map[string]interface{}{
			"model": model,
		},
	)

	if err != nil {
		return nil, err
	}

	return &PushModelResult{
		StatusMessages: res.StatusMessages,
	}, nil
}
