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

// ModelInfoResult represents the structure of the response returned by the FetchModelInfo method.
// It includes various details about a model, such as its configuration, parameters, and additional metadata.
type ModelInfoResult struct {
	Modelfile  string                 `json:"modelfile"`  // The content or path of the model file associated with the model.
	Parameters string                 `json:"parameters"` // A string representation of the model's parameters.
	Template   string                 `json:"template"`   // The template used for the model, possibly indicating its architecture or purpose.
	Details    map[string]interface{} `json:"details"`    // A map containing additional detailed information about the model.
	ModelInfo  map[string]interface{} `json:"model_info"` // A map containing metadata or attributes specific to the model.
}

// FetchModelInfo retrieves information about a specific model from the server.
// It sends a POST request to the "/api/show" endpoint with the model name and verbosity flag.
// The server's response is decoded into a ModelInfoResult struct.
//
// Parameters:
//   - ctx: The context.Context for managing request deadlines and cancellations.
//   - model: The name or identifier of the model whose information is being requested.
//   - verbose: A boolean flag indicating whether to request detailed information (true) or basic information (false).
//
// Returns:
//   - A pointer to a ModelInfoResult containing the model's information.
//   - An error if the request fails or the response cannot be processed.
func (c *Client) FetchModelInfo(
	ctx context.Context,
	model string,
	verbose bool,
) (*ModelInfoResult, error) {
	rel := &url.URL{Path: "/api/show"}
	u := c.BaseURL.ResolveReference(rel)

	return c.sendShowRequest(
		ctx,
		"POST",
		u.String(),
		map[string]interface{}{
			"model":   model,
			"verbose": verbose,
		},
	)
}
