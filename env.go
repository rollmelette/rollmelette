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
	appAddress    common.Address
	appAddressSet bool
}

func newEnv(addressBook AddressBook, rollup rollupEnv) *env {
	return &env{
		AddressBook: addressBook,
		rollup:      rollup,
	}
}

// handlers ////////////////////////////////////////////////////////////////////////////////////////

func (e *env) handleAdvance(app Application, input *advanceInput) error {
	slog.Debug("received advance",
		"payload", hexutil.Encode(input.Payload),
		"inputIndex", input.Metadata.InputIndex,
		"msgSender", input.Metadata.MsgSender,
		"blockNumber", input.Metadata.BlockNumber,
		"blockTimestamp", input.Metadata.BlockTimestamp,
	)
	if input.Metadata.MsgSender == e.DAppAddressRelay {
		return e.handleDAppAddressRelay(input.Payload)
	}
	return app.Advance(e, input.Metadata, input.Payload)
}

func (e *env) handleDAppAddressRelay(payload []byte) error {
	if len(payload) != 20 {
		return fmt.Errorf("invalid input from DAppAddressRelay: %x", payload)
	}
	e.appAddress = (common.Address)(payload)
	e.appAddressSet = true
	slog.Debug("got application address from relay", "address", e.appAddress)
	return nil
}

func (e *env) handleInspect(app Application, payload []byte) error {
	slog.Debug("received inspect", "payload", hexutil.Encode(payload))
	return app.Inspect(e, payload)
}

// EnvInspector interface //////////////////////////////////////////////////////////////////////////

func (e *env) Report(payload []byte) {
	slog.Debug("sending report", "payload", hexutil.Encode(payload))
	err := e.rollup.sendReport(payload)
	if err != nil {
		panic(err)
	}
}

func (e *env) Reportf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	slog.Debug("sending report", "payload", message)
	err := e.rollup.sendReport([]byte(message))
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
