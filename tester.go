// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
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

	app        Application
	rollup     *RollupMock
	env        *env
	inputIndex int
}

// NewTester creates a Tester for the given application
func NewTester(app Application) *Tester {
	rollup := &RollupMock{}
	return &Tester{
		MsgSender:  common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa"),
		app:        app,
		rollup:     rollup,
		env:        newEnv(NewAddressBook(), rollup),
		inputIndex: 0,
	}
}

// Advance sends an advance input to the application.
// It returns the metadata sent to the app and the outputs received from the app.
func (t *Tester) Advance(payload []byte) TestAdvanceResult {
	t.rollup.reset()
	metadata := Metadata{
		InputIndex:     t.inputIndex,
		MsgSender:      t.MsgSender,
		BlockNumber:    int64(t.inputIndex),
		BlockTimestamp: time.Now().Unix(),
	}
	input := advanceInput{
		Metadata: metadata,
		Payload:  payload,
	}
	err := t.env.handleAdvance(t.app, &input)
	t.inputIndex++
	return TestAdvanceResult{
		RollupMock: *t.rollup,
		Metadata:   metadata,
		Err:        err,
	}
}

// Inspect sends an inspect input to the application.
// It returns the outputs received from the app.
func (t *Tester) Inspect(payload []byte) TestInspectResult {
	t.rollup.reset()
	err := t.env.handleInspect(t.app, payload)
	return TestInspectResult{
		Reports: t.rollup.Reports,
		Err:     err,
	}
}
