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
	"net/url"
	"time"
)

// PromptInfo represents the structure of a prompt request,
// including various parameters to customize the prompt generation.
type PromptInfo struct {
	Model    string `json:"model"`              // The model to be used for prompt generation.
	Prompt   string `json:"prompt,omitempty"`   // The prompt text; optional field.
	Suffix   string `json:"suffix,omitempty"`   // Suffix to append to the prompt; optional field.
	System   string `json:"system,omitempty"`   // System context or instructions; optional field.
	Template string `json:"template,omitempty"` // Template to structure the prompt; optional field.

	Images  []string               `json:"images,omitempty"`  // List of images associated with the prompt; optional field.
	Format  interface{}            `json:"format,omitempty"`  // Format of the prompt; can be string or map; optional field.
	Options map[string]interface{} `json:"options,omitempty"` // Additional options for prompt customization; optional field.

	Stream *bool `json:"stream,omitempty"` // Flag to indicate streaming response; optional field.
	Raw    *bool `json:"raw,omitempty"`    // Flag to indicate raw response; optional field.

	KeepAlive string      `json:"keep_alive,omitempty"` // Connection keep-alive duration; optional field.
	Context   interface{} `json:"context,omitempty"`    // Contextual information for the prompt; optional field.
}

// PromptResult represents the structure of a prompt response,
// containing details about the generated prompt and its processing metrics.
type PromptResult struct {
	Model      string      `json:"model"`                 // The model used for prompt generation.
	Response   string      `json:"response"`              // The generated response.
	CreatedAt  time.Time   `json:"created_at"`            // Timestamp of when the prompt was created.
	Context    interface{} `json:"context,omitempty"`     // Contextual information associated with the prompt; optional field.
	Done       bool        `json:"done"`                  // Flag indicating if the prompt processing is complete.
	DoneReason string      `json:"done_reason,omitempty"` // Reason for completion; optional field.

	TotalDuration      int64 `json:"total_duration,omitempty"`       // Total time taken for prompt processing; optional field.
	LoadDuration       int64 `json:"load_duration,omitempty"`        // Time taken to load resources; optional field.
	PromptEvalCount    int   `json:"prompt_eval_count,omitempty"`    // Number of prompt evaluations; optional field.
	PromptEvalDuration int64 `json:"prompt_eval_duration,omitempty"` // Duration of prompt evaluations; optional field.
	EvalCount          int   `json:"eval_count,omitempty"`           // Number of evaluations performed; optional field.
	EvalDuration       int64 `json:"eval_duration,omitempty"`        // Duration of evaluations; optional field.
}

// ValidatePromptInfo validates the fields of the PromptInfo struct,
// ensuring that fields like Format and Context have appropriate types.
func (req *PromptInfo) ValidatePromptInfo() error {
	if req.Format != nil {
		switch req.Format.(type) {
		case string, map[string]interface{}:
		default:
			return fmt.Errorf(
				"invalid type for Format field; must be string or map[string]interface{}",
			)
		}
	}

	if req.Context != nil {
		switch ctx := req.Context.(type) {
		case []int:
		case []interface{}:
			for idx, v := range ctx {
				if _, ok := v.(int); !ok {
					return fmt.Errorf(
						"invalid type for Context at index %d; expected int",
						idx,
					)
				}
			}

		default:
			return fmt.Errorf(
				"invalid type for Context field; expected []int or []interface{} of ints",
			)
		}
	}

	return nil
}

// Generate sends a prompt generation request to the server,
// validates the PromptInfo, and returns the generated PromptResult.
func (c *Client) Generate(
	ctx context.Context,
	req *PromptInfo,
) (*PromptResult, error) {
	if err := req.ValidatePromptInfo(); err != nil {
		return nil, err
	}

	rel := &url.URL{Path: "/api/generate"}
	u := c.BaseURL.ResolveReference(rel)

	return c.sendRequest(ctx, "POST", u.String(), req)
}
