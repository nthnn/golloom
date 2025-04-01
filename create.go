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

// CreateModelRequest represents the payload for creating a new model on the server.
// It includes various configuration options and parameters for the model creation process.
type CreateModelRequest struct {
	Model      string                 `json:"model"`                // The name or identifier for the new model.
	From       string                 `json:"from,omitempty"`       // Optional. The source model to base the new model on.
	Files      map[string]string      `json:"files,omitempty"`      // Optional. A map of file names to their content, used in model creation.
	Adapters   map[string]string      `json:"adapters,omitempty"`   // Optional. A map specifying adapter configurations.
	Template   string                 `json:"template,omitempty"`   // Optional. A template defining the model's structure or behavior.
	License    interface{}            `json:"license,omitempty"`    // Optional. License information for the model.
	System     string                 `json:"system,omitempty"`     // Optional. System-specific parameters or configurations.
	Parameters map[string]interface{} `json:"parameters,omitempty"` // Optional. Additional parameters for model creation.
	Messages   []Message              `json:"messages,omitempty"`   // Optional. A sequence of messages or instructions related to the model.
	Stream     *bool                  `json:"stream,omitempty"`     // Optional. Indicates if the creation process should be streamed.
	Quantize   string                 `json:"quantize,omitempty"`   // Optional. Specifies quantization settings for the model.
}

// CreateModelResult encapsulates the outcome of a model creation operation,
// including any status messages returned during the process.
type CreateModelResult struct {
	// A list of messages detailing the status or progress of the creation operation.
	StatusMessages []string `json:"status_messages"`
}

// CreateModel sends a request to create a new model on the server using the provided configuration.
// It constructs the appropriate API endpoint and sends a POST request with the creation parameters.
// Parameters:
//   - ctx: A context.Context object for managing request deadlines and cancellations.
//   - req: A pointer to a CreateModelRequest struct containing the model creation parameters.
//
// Returns:
//   - A pointer to a CreateModelResult struct containing status messages from the operation.
//   - An error if the request fails or the server returns an unexpected response.
func (c *Client) CreateModel(
	ctx context.Context,
	req *CreateModelRequest,
) (*CreateModelResult, error) {
	rel := &url.URL{Path: "/api/create"}
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

	return &CreateModelResult{
		StatusMessages: res.StatusMessages,
	}, nil
}
