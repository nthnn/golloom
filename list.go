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
	"time"
)

// ModelDetails encapsulates detailed attributes of a machine learning model,
// providing insights into its format, family, parameter size, and quantization level.
type ModelDetails struct {
	Format            string   `json:"format"`             // The storage format of the model (e.g., ONNX, TensorFlow, PyTorch).
	Family            string   `json:"family"`             // The primary category or group the model belongs to (e.g., Vision, NLP).
	Families          []string `json:"families,omitempty"` // Additional categories or groups the model is associated with.
	ParameterSize     string   `json:"parameter_size"`     // The total number of parameters in the model, often indicative of its complexity.
	QuantizationLevel string   `json:"quantization_level"` // The degree of quantization applied to the model, affecting its size and performance.
}

// ModelInfo provides metadata about a specific machine learning model,
// including its name, last modification timestamp, file size, checksum digest, and detailed attributes.
type ModelInfo struct {
	Name       string       `json:"name"`        // The unique identifier or name of the model.
	ModifiedAt time.Time    `json:"modified_at"` // The timestamp indicating the last modification time of the model.
	Size       int64        `json:"size"`        // The size of the model file in bytes.
	Digest     string       `json:"digest"`      // The checksum or hash digest of the model file for integrity verification.
	Details    ModelDetails `json:"details"`     // An embedded struct containing detailed attributes of the model.
}

// ModelList represents a collection of machine learning models,
// serving as a container for multiple ModelInfo entries.
type ModelList struct {
	Models []ModelInfo `json:"models"` // A slice containing information about each available model.
}

// ListModels retrieves a list of available machine learning models from the server.
// It constructs the appropriate API endpoint, sends a GET request, and decodes the JSON response into a ModelList.
//
// Parameters:
//   - ctx: The context.Context for managing request deadlines and cancellations.
//
// Returns:
//   - A pointer to a ModelList containing metadata about available models.
//   - An error if the request fails or the response cannot be processed.
func (c *Client) ListModels(ctx context.Context) (*ModelList, error) {
	rel := &url.URL{Path: "/api/tags"}
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

	var listResp ModelList
	err = json.NewDecoder(resp.Body).Decode(&listResp)
	if err != nil {
		return nil, err
	}

	return &listResp, nil
}
