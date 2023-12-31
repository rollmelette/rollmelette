// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// env is the implementation of the Env interface.
// This struct isn't public because we don't want users of Rollmelette to create it.
// Instead, it is create by the running and testing functions.
type env struct {
	AddressBook
	rollup        rollupEnv
	app           Application
	appAddress    common.Address
	appAddressSet bool
}

func newEnv(addressBook AddressBook, rollup rollupEnv, app Application) *env {
	return &env{
		AddressBook: addressBook,
		rollup:      rollup,
		app:         app,
	}
}

// handlers ////////////////////////////////////////////////////////////////////////////////////////

func (e *env) handle(input any) (err error) {
	defer func() {
		// Recover from panic so we can safely reject the input and print an error message.
		panicObj := recover()
		if panicObj != nil {
			panicErr, ok := panicObj.(error)
			if ok {
				err = panicErr
			} else {
				err = fmt.Errorf("a panic occured: %v", panicObj)
			}
		}
		if err != nil {
			slog.Error("input rejected", "error", err)
		}
	}()
	switch input := input.(type) {
	case *advanceInput:
		return e.handleAdvance(input)
	case *inspectInput:
		return e.handleInspect(input.Payload)
	default:
		// impossible
		panic("invalid input type")
	}
}

func (e *env) handleAdvance(input *advanceInput) error {
	slog.Debug("received advance",
		"payload", hexutil.Encode(input.Payload),
		"inputIndex", input.Metadata.InputIndex,
		"msgSender", input.Metadata.MsgSender,
		"blockNumber", input.Metadata.BlockNumber,
		"blockTimestamp", input.Metadata.BlockTimestamp,
	)
	if input.Metadata.MsgSender == e.AppAddressRelay {
		return e.handleAppAddressRelay(input.Payload)
	}
	return e.app.Advance(e, input.Metadata, input.Payload)
}

func (e *env) handleAppAddressRelay(payload []byte) error {
	if len(payload) != 20 {
		return fmt.Errorf("invalid input from app address relay: %x", payload)
	}
	e.appAddress = (common.Address)(payload)
	e.appAddressSet = true
	slog.Debug("got application address from relay", "address", e.appAddress)
	return nil
}

func (e *env) handleInspect(payload []byte) error {
	slog.Debug("received inspect", "payload", hexutil.Encode(payload))
	return e.app.Inspect(e, payload)
}

// EnvInspector interface //////////////////////////////////////////////////////////////////////////

func (e *env) Report(payload []byte) {
	slog.Debug("sending report", "payload", hexutil.Encode(payload))
	err := e.rollup.sendReport(payload)
	if err != nil {
		panic(err)
	}
}

func (e *env) AppAddress() (common.Address, bool) {
	return e.appAddress, e.appAddressSet
}

// EnvInspector interface //////////////////////////////////////////////////////////////////////////

func (e *env) Voucher(destination common.Address, payload []byte) int {
	slog.Debug("sending voucher", "destination", destination, "payload", hexutil.Encode(payload))
	index, err := e.rollup.sendVoucher(destination, payload)
	if err != nil {
		panic(err)
	}
	return index
}

func (e *env) Notice(payload []byte) int {
	slog.Debug("sending notice", "payload", hexutil.Encode(payload))
	index, err := e.rollup.sendNotice(payload)
	if err != nil {
		panic(err)
	}
	return index
}
