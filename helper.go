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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// sendRequest constructs and sends an HTTP request with the specified method, URL, and body.
// It encodes the body as JSON, sets appropriate headers, and processes the server's response.
// Parameters:
//   - ctx: A context.Context for managing request deadlines and cancellations.
//   - method: The HTTP method (e.g., "GET", "POST") to use for the request.
//   - urlStr: The target URL as a string.
//   - body: The payload to be sent with the request; it will be JSON-encoded.
//
// Returns:
//   - A pointer to a PromptResult containing the server's response data.
//   - An error if the request fails or the response cannot be processed.
func (c *Client) sendRequest(
	ctx context.Context,
	method, urlStr string,
	body interface{},
) (*PromptResult, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		urlStr,
		buf,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"HTTP request failed with status %d: %s",
			resp.StatusCode,
			errorBody,
		)
	}

	var genResp PromptResult
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return nil, err
	}

	return &genResp, nil
}

// sendChatRequest functions similarly to sendRequest but expects a response of type ModelResponse.
// It constructs and sends an HTTP request with the specified method, URL, and body, then decodes the response.
// Parameters:
//   - ctx: A context.Context for managing request deadlines and cancellations.
//   - method: The HTTP method to use for the request.
//   - urlStr: The target URL as a string.
//   - body: The payload to be sent with the request; it will be JSON-encoded.
//
// Returns:
//   - A pointer to a ModelResponse containing the server's response data.
//   - An error if the request fails or the response cannot be processed.
func (c *Client) sendChatRequest(
	ctx context.Context,
	method, urlStr string,
	body interface{},
) (*ModelResponse, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method, urlStr,
		buf,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var chatResp ModelResponse
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil {
		return nil, err
	}

	return &chatResp, nil
}

// sendShowRequest constructs and sends an HTTP request, expecting a response of type ModelInfoResult.
// It encodes the body as JSON, sets appropriate headers, and processes the server's response.
// Parameters:
//   - ctx: A context.Context for managing request deadlines and cancellations.
//   - method: The HTTP method to use for the request.
//   - urlStr: The target URL as a string.
//   - body: The payload to be sent with the request; it will be JSON-encoded.
//
// Returns:
//   - A pointer to a ModelInfoResult containing the server's response data.
//   - An error if the request fails or the response cannot be processed.
func (c *Client) sendShowRequest(
	ctx context.Context,
	method, urlStr string,
	body interface{},
) (*ModelInfoResult, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var showResp ModelInfoResult
	err = json.NewDecoder(resp.Body).Decode(&showResp)
	if err != nil {
		return nil, err
	}

	return &showResp, nil
}

// sendEmbedRequest constructs and sends an HTTP request, expecting a response of type EmbedResult.
// It encodes the body as JSON, sets appropriate headers, and processes the server's response.
// Parameters:
//   - ctx: A context.Context for managing request deadlines and cancellations.
//   - method: The HTTP method to use for the request.
//   - urlStr: The target URL as a string.
//   - body: The payload to be sent with the request; it will be JSON-encoded.
//
// Returns:
//   - A pointer to an EmbedResult containing the server's response data.
//   - An error if the request fails or the response cannot be processed.
func (c *Client) sendEmbedRequest(
	ctx context.Context,
	method, urlStr string,
	body interface{},
) (*EmbedResult, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var embedResp EmbedResult
	err = json.NewDecoder(resp.Body).Decode(&embedResp)
	if err != nil {
		return nil, err
	}

	return &embedResp, nil
}

// sendStatusStreamRequest constructs and sends an HTTP request to a specified URL with the given method and body.
// It expects a streaming response containing status messages and processes them accordingly.
// Parameters:
//   - ctx: A context.Context for managing request deadlines and cancellations.
//   - method: The HTTP method (e.g., "POST") to use for the request.
//   - urlStr: The target URL as a string.
//   - body: The payload to be sent with the request; it will be JSON-encoded.
//
// Returns:
//   - A pointer to a struct containing a slice of status messages received from the server.
//   - An error if the request fails, the response cannot be processed, or if the number of status messages exceeds the maximum allowed.
func (c *Client) sendStatusStreamRequest(
	ctx context.Context,
	method, urlStr string,
	body interface{},
) (
	*struct {
		StatusMessages []string `json:"status_messages"`
	},
	error,
) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		urlStr,
		buf,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		limitedBody, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf(
			"HTTP request failed with status %d: %s",
			resp.StatusCode,
			string(limitedBody),
		)
	}

	var msgs []string
	maxMessages := 1000

	dec := json.NewDecoder(resp.Body)
	for dec.More() && len(msgs) < maxMessages {
		var s struct {
			Status string `json:"status"`
		}

		if err := dec.Decode(&s); err != nil {
			return nil, fmt.Errorf("error decoding stream: %w", err)
		}

		msgs = append(msgs, s.Status)
	}

	if len(msgs) >= maxMessages {
		return nil, fmt.Errorf("streaming response exceeded maximum allowed messages")
	}

	return &struct {
		StatusMessages []string `json:"status_messages"`
	}{StatusMessages: msgs}, nil
}
