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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// CheckBlobExists checks if a blob with the given digest exists on the server.
// It sends a HEAD request to the /api/blobs/{digest} endpoint.
// Parameters:
//   - ctx: A context to control request lifetime (e.g., cancellation).
//   - digest: A string representing the blob's digest or identifier.
//
// Returns:
//   - A boolean value: true if the blob exists, false otherwise.
//   - An error if the request fails or if the digest is invalid.
func (c *Client) CheckBlobExists(
	ctx context.Context,
	digest string,
) (bool, error) {
	if strings.Contains(digest, "/") || strings.Contains(digest, "..") {
		return false, fmt.Errorf("invalid digest: %s", digest)
	}

	safeDigest := url.PathEscape(digest)
	rel := &url.URL{Path: "/api/blobs/" + safeDigest}
	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(
		ctx,
		"HEAD",
		u.String(),
		nil,
	)

	if err != nil {
		return false, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil

	case http.StatusNotFound:
		return false, nil

	default:
		return false, fmt.Errorf(
			"unexpected status code: %d",
			resp.StatusCode,
		)
	}
}

// PushBlob uploads a blob to the server.
// It sends a POST request with the blob data to the /api/blobs/{digest} endpoint.
// Parameters:
//   - ctx: A context to control request lifetime (e.g., cancellation).
//   - digest: A string representing the blob's digest or identifier.
//   - file: An io.Reader that provides the blob's data.
//
// Returns:
//   - An error if the upload fails.
func (c *Client) PushBlob(
	ctx context.Context,
	digest string,
	file io.Reader,
) error {
	rel := &url.URL{Path: path.Join("/api/blobs", digest)}
	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		u.String(),
		file,
	)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"failed to push blob: %s",
			string(body),
		)
	}

	return nil
}
