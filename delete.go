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

// DeleteModelRequest represents the request payload for deleting a model from the server.
// It contains the identifier of the model that should be removed.
type DeleteModelRequest struct {
	// Model specifies the unique identifier or name of the model to be deleted.
	Model string `json:"model"`
}

// DeleteModelResult encapsulates the response returned by the server after attempting to delete a model.
// It includes a slice of status messages that describe the outcome of the operation.
type DeleteModelResult struct {
	// StatusMessages holds any messages returned by the server regarding the deletion operation.
	StatusMessages []string `json:"status_messages"`
}

// DeleteModel sends a request to delete a model from the server.
// It constructs the API endpoint URL for deletion and issues a POST request with the provided DeleteModelRequest.
// Parameters:
//   - ctx: A context.Context for controlling cancellation and timeouts during the HTTP request.
//   - req: A pointer to a DeleteModelRequest containing the model identifier to be deleted.
//
// Returns:
//   - A pointer to a DeleteModelResult containing the server's status messages about the deletion.
//   - An error if the request fails or the server returns an error.
func (c *Client) DeleteModel(
	ctx context.Context,
	req *DeleteModelRequest,
) (*DeleteModelResult, error) {
	rel := &url.URL{Path: "/api/delete"}
	u := c.BaseURL.ResolveReference(rel)

	res, err := c.sendStatusStreamRequest(
		ctx,
		"POST",
		u.String(),
		req,
	)

	if err != nil {
		return nil, err
	}

	return &DeleteModelResult{
		StatusMessages: res.StatusMessages,
	}, nil
}
