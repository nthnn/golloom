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
	"net/http"
	"net/url"
	"time"
)

// Client is a struct that encapsulates the configuration and HTTP client used to interact
// with the Ollama server (or similar API endpoints). It holds the base URL for the API and
// an HTTP client instance with a configurable timeout.
type Client struct {
	// BaseURL is the base URL of the server to which requests will be sent.
	BaseURL *url.URL
	// HTTPClient is the HTTP client used to make requests. Its configuration (e.g., timeout)
	// can be set during client initialization.
	HTTPClient *http.Client
}

// NewClient creates a new instance of Client configured to communicate with the server.
// It takes a base URL as a string and a duration (in minutes) for setting the HTTP client's timeout.
// Parameters:
//   - baseURL: A string representing the server's base URL.
//   - minutes: A time.Duration value (in minutes) that sets the timeout for HTTP requests.
//
// Returns:
//   - A pointer to a Client instance properly configured with the base URL and HTTP client.
//   - An error if the provided baseURL cannot be parsed.
func NewClient(
	baseURL string,
	minutes time.Duration,
) (*Client, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		BaseURL: parsed,
		HTTPClient: &http.Client{
			Timeout: minutes * time.Minute,
		},
	}, nil
}
