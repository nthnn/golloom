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
	"time"
)

// EmbedResult represents the response from an embedding operation.
// It includes details about the model used, the creation timestamp, and the resulting embedding data.
type EmbedResult struct {
	Model     string      `json:"model"`      // The identifier of the model used to generate the embedding.
	CreatedAt time.Time   `json:"created_at"` // The timestamp indicating when the embedding was created.
	Embedding interface{} `json:"embedding"`  // The actual embedding data; its structure depends on the model's output.
}

// Embed sends a request to generate an embedding for the given input using the specified model and options.
// It constructs the appropriate API endpoint and sends a POST request with the necessary parameters.
// Parameters:
//   - ctx: A context.Context object for managing request deadlines and cancellations.
//   - model: The name or identifier of the model to use for generating the embedding.
//   - input: The input data for which the embedding is to be generated.
//   - options: A map of additional options to customize the embedding process.
//
// Returns:
//   - A pointer to an EmbedResult struct containing the embedding and related metadata.
//   - An error if the request fails or the server returns an unexpected response.
func (c *Client) Embed(
	ctx context.Context,
	model, input string,
	options map[string]interface{},
) (*EmbedResult, error) {
	rel := &url.URL{Path: "/api/embed"}
	u := c.BaseURL.ResolveReference(rel)

	return c.sendEmbedRequest(
		ctx,
		"POST",
		u.String(),
		map[string]interface{}{
			"model":   model,
			"input":   input,
			"options": options,
		},
	)
}
