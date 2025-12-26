// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// TestAdvanceResult is the result of the test advance function.
type TestAdvanceResult struct {
	Vouchers             []TestVoucher
	DelegateCallVouchers []TestDelegateCallVoucher
	Notices              []TestNotice
	Reports              []TestReport
	Metadata
	Err error
}

// TestInspectResult
type TestInspectResult struct {
	Reports []TestReport
	Err     error
}

// Tester is an unit tester for the Application.
type Tester struct {
	rollup *rollupMock
	book   AddressBook
	env    *env
	index  int
}

// NewTester creates a Tester for the given application
func NewTester(app Application) *Tester {
	rollup := &rollupMock{}
	book := NewAddressBook()
	return &Tester{
		rollup: rollup,
		book:   book,
		env:    newEnv(context.Background(), book, rollup, app),
		index:  0,
	}
}

// Book returns the address book used by the tester.
func (t *Tester) Book() AddressBook {
	return t.book
}

// Advance sends an advance input to the application.
// It returns the metadata sent to the app and the outputs received from the app.
func (t *Tester) Advance(msgSender common.Address, payload []byte) TestAdvanceResult {
	return t.sendAdvance(msgSender, payload)
}

// DepositEther simulates an advance input from the Ether portal.
func (t *Tester) DepositEther(
	msgSender common.Address,
	value *big.Int,
	payload []byte,
) TestAdvanceResult {
	if value.Cmp(MaxUint256) > 0 {
		panic("value too big")
	} else if value.Cmp(big.NewInt(0)) < 0 {
		panic("negative value")
	}
	portalPayload := make([]byte, 0, common.AddressLength+common.HashLength+len(payload))
	portalPayload = append(portalPayload, msgSender[:]...)
	portalPayload = append(portalPayload, value.FillBytes(make([]byte, common.HashLength))...)
	portalPayload = append(portalPayload, payload...)
	return t.sendAdvance(t.env.EtherPortal, portalPayload)
}

// DepositERC20 simulates an advance input from the ERC20 portal.
func (t *Tester) DepositERC20(
	token common.Address,
	msgSender common.Address,
	value *big.Int,
	payload []byte,
) TestAdvanceResult {
	if value.Cmp(MaxUint256) > 0 {
		panic("value too big")
	} else if value.Cmp(big.NewInt(0)) < 0 {
		panic("negative value")
	}
	portalPayload := make([]byte, 0, common.AddressLength+2*common.HashLength+len(payload))
	portalPayload = append(portalPayload, token[:]...)
	portalPayload = append(portalPayload, msgSender[:]...)
	portalPayload = append(portalPayload, value.FillBytes(make([]byte, common.HashLength))...)
	portalPayload = append(portalPayload, payload...)
	return t.sendAdvance(t.env.ERC20Portal, portalPayload)
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
		ChainId:        1,
		AppContract:    common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e"),
		Index:          t.index,
		MsgSender:      msgSender,
		BlockNumber:    int64(t.index),
		BlockTimestamp: time.Now().Unix(),
		PrevRandao:     "0x0000000000000000000000000000000000000000000000000000000000000001",
	}
	input := advanceInput{
		Metadata: metadata,
		Payload:  payload,
	}
	err := t.env.handle(&input)
	t.index++
	return TestAdvanceResult{
		Vouchers:             t.rollup.Vouchers,
		DelegateCallVouchers: t.rollup.DelegateCallVouchers,
		Notices:              t.rollup.Notices,
		Reports:              t.rollup.Reports,
		Metadata:             metadata,
		Err:                  err,
	}
}
