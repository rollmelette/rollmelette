// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// TestAdvanceResult is the result of the test advance function.
type TestAdvanceResult struct {
	RollupMock
	Metadata
	Err error
}

// TestInspectResult
type TestInspectResult struct {
	Reports []ReportMock
	Err     error
}

// Tester is an unit tester for the Application.
type Tester struct {

	// MsgSender is forwarded to the metadata when sending something.
	MsgSender common.Address

	// AppAddress is the address of the application.
	AppAddress common.Address

	rollup     *RollupMock
	env        *env
	inputIndex int
}

// NewTester creates a Tester for the given application
func NewTester(app Application) *Tester {
	rollup := &RollupMock{}
	return &Tester{
		MsgSender:  common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"),
		AppAddress: common.HexToAddress("0x70ac08179605AF2D9e75782b8DEcDD3c22aA4D0C"),
		rollup:     rollup,
		env:        newEnv(context.Background(), NewAddressBook(), rollup, app),
		inputIndex: 0,
	}
}

// Advance sends an advance input to the application.
// It returns the metadata sent to the app and the outputs received from the app.
func (t *Tester) Advance(payload []byte) TestAdvanceResult {
	return t.sendAdvance(t.MsgSender, payload)
}

// RelayAppAddress simulates an advance input from the app address relay.
func (t *Tester) RelayAppAddress() TestAdvanceResult {
	return t.sendAdvance(t.env.AppAddressRelay, t.AppAddress[:])
}

// Inspect sends an inspect input to the application.
// It returns the outputs received from the app.
func (t *Tester) Inspect(payload []byte) TestInspectResult {
	t.rollup.reset()
	input := inspectInput{
		Payload: payload,
	}
	err := t.env.handle(&input)
	return TestInspectResult{
		Reports: t.rollup.Reports,
		Err:     err,
	}
}

func (t *Tester) sendAdvance(msgSender common.Address, payload []byte) TestAdvanceResult {
	t.rollup.reset()
	metadata := Metadata{
		InputIndex:     t.inputIndex,
		MsgSender:      msgSender,
		BlockNumber:    int64(t.inputIndex),
		BlockTimestamp: time.Now().Unix(),
	}
	input := advanceInput{
		Metadata: metadata,
		Payload:  payload,
	}
	err := t.env.handle(&input)
	t.inputIndex++
	return TestAdvanceResult{
		RollupMock: *t.rollup,
		Metadata:   metadata,
		Err:        err,
	}
}
