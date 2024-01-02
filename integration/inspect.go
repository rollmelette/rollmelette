// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Report from the inspect result.
type InspectReport struct {
	Payload string `json:"payload"`
}

// InspectCompletionStatus is the completion status for the inspect result.
type InspectCompletionStatus string

const (
	Accepted           InspectCompletionStatus = "Accepted"
	CycleLimitExceeded InspectCompletionStatus = "CycleLimitExceeded"
	Exception          InspectCompletionStatus = "Exception"
	MachineHalted      InspectCompletionStatus = "MachineHalted"
	Rejected           InspectCompletionStatus = "Rejected"
	TimeLimitExceeded  InspectCompletionStatus = "TimeLimitExceeded"
)

// InspectResult is the response of the inspect API.
type InspectResult struct {
	ExceptionPayload    string                  `json:"exception_payload"`
	ProcessedInputCount int                     `json:"processed_input_count"`
	Reports             []InspectReport         `json:"reports"`
	Status              InspectCompletionStatus `json:"status"`
}

// Inspect makes an HTTP inspect request to the node.
func Inspect(ctx context.Context, url string, payload []byte) (*InspectResult, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read body: %w", err)
		}
		return nil, fmt.Errorf("invalid status %v: %v", resp.StatusCode, string(body))
	}
	var result InspectResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}
