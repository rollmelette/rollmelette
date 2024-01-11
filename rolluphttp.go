// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// rollupHttp implements the Rollup API by calling the Rollup HTTP server.
type rollupHttp struct {
	url string
}

// newRollupHttp create a new rollup HTTP client.
func newRollupHttp(url string) *rollupHttp {
	return &rollupHttp{
		url: url,
	}
}

// rollup interface ////////////////////////////////////////////////////////////////////////////////

func (r *rollupHttp) finishAndGetNext(ctx context.Context, status finishStatus) (any, error) {
	request := struct {
		Status string `json:"status"`
	}{
		Status: string(status),
	}
	resp, err := r.sendPost(ctx, "finish", request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusAccepted {
		// if we get StatusAccepted we should trying again
		return r.finishAndGetNext(ctx, status)
	}
	if err = checkStatusOk(resp); err != nil {
		return nil, err
	}
	var finishResp struct {
		RequestType string          `json:"request_type"`
		Data        json.RawMessage `json:"data"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&finishResp); err != nil {
		return nil, fmt.Errorf("rollup: decode finish response: %w", err)
	}
	switch finishResp.RequestType {
	case "advance_state":
		return parseAdvanceInput(finishResp.Data)
	case "inspect_state":
		return parseInspectInput(finishResp.Data)
	default:
		return nil, fmt.Errorf("rollup: invalid request type: %v", finishResp.RequestType)
	}
}

func (r *rollupHttp) sendVoucher(ctx context.Context, destination common.Address, payload []byte) (int, error) {
	request := struct {
		Destination string `json:"destination"`
		Payload     string `json:"payload"`
	}{
		Destination: hexutil.Encode(destination[:]),
		Payload:     hexutil.Encode(payload),
	}
	resp, err := r.sendPost(ctx, "voucher", request)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if err = checkStatusOk(resp); err != nil {
		return 0, err
	}
	return parseOutputIndex(resp.Body)
}

func (r *rollupHttp) sendNotice(ctx context.Context, payload []byte) (int, error) {
	request := struct {
		Payload string `json:"payload"`
	}{
		Payload: hexutil.Encode(payload),
	}
	resp, err := r.sendPost(ctx, "notice", request)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if err = checkStatusOk(resp); err != nil {
		return 0, err
	}
	return parseOutputIndex(resp.Body)
}

func (r *rollupHttp) sendReport(ctx context.Context, payload []byte) error {
	request := struct {
		Payload string `json:"payload"`
	}{
		Payload: hexutil.Encode(payload),
	}
	resp, err := r.sendPost(ctx, "report", request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err = checkStatusOk(resp); err != nil {
		return err
	}
	return nil
}

// helpers /////////////////////////////////////////////////////////////////////////////////////////

// sendPost sends a POST request and returns the HTTP response.
// The callee should close the response body.
func (r *rollupHttp) sendPost(ctx context.Context, route string, request any) (*http.Response, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("rollup: serialize request: %w", err)
	}
	endpoint := fmt.Sprintf("%v/%v", r.url, route)
	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("rollup: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rollup: do request: %w", err)
	}
	return resp, nil
}

func parseOutputIndex(r io.Reader) (int, error) {
	var outputResp struct {
		Index int `json:"index"`
	}
	if err := json.NewDecoder(r).Decode(&outputResp); err != nil {
		return 0, fmt.Errorf("rollup: decode finish response: %w", err)
	}
	return outputResp.Index, nil
}

func parseAdvanceInput(data json.RawMessage) (any, error) {
	var advanceRequest struct {
		Payload  string `json:"payload"`
		Metadata struct {
			MsgSender   string `json:"msg_sender"`
			EpochIndex  int    `json:"epoch_index"`
			InputIndex  int    `json:"input_index"`
			BlockNumber int64  `json:"block_number"`
			Timestamp   int64  `json:"timestamp"`
		}
	}
	if err := json.Unmarshal(data, &advanceRequest); err != nil {
		return nil, fmt.Errorf("rollup: decode advance input: %w", err)
	}
	payload, err := hexutil.Decode(advanceRequest.Payload)
	if err != nil {
		return nil, fmt.Errorf("rollup: decode advance payload: %w", err)
	}
	sender, err := hexutil.Decode(advanceRequest.Metadata.MsgSender)
	if err != nil {
		return nil, fmt.Errorf("rollup: decode advance metadata sender: %w", err)
	}
	metadata := Metadata{
		InputIndex:     advanceRequest.Metadata.InputIndex,
		MsgSender:      common.Address(sender),
		BlockNumber:    advanceRequest.Metadata.BlockNumber,
		BlockTimestamp: advanceRequest.Metadata.Timestamp,
	}
	input := &advanceInput{
		Metadata: metadata,
		Payload:  payload,
	}
	return input, nil
}

func parseInspectInput(data json.RawMessage) (any, error) {
	var inspectRequest struct {
		Payload string `json:"payload"`
	}
	if err := json.Unmarshal(data, &inspectRequest); err != nil {
		return nil, fmt.Errorf("rollup: decode advance request: %v", err)
	}
	payload, err := hexutil.Decode(inspectRequest.Payload)
	if err != nil {
		return nil, fmt.Errorf("rollup: decode advance payload: %v", err)
	}
	input := &inspectInput{
		Payload: payload,
	}
	return input, nil
}

func checkStatusOk(resp *http.Response) error {
	if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("http: read body: %w", err)
		}
		return fmt.Errorf("http: invalid status %v: %v", resp.StatusCode, string(body))
	}
	return nil
}
