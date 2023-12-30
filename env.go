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
	rollup rollup
}

func newEnv(rollup rollup) *env {
	return &env{
		rollup: rollup,
	}
}

// handlers ////////////////////////////////////////////////////////////////////////////////////////

func (e *env) handleAdvance(app Application, input *advanceInput) error {
	return app.Advance(e, input.Metadata, input.Payload)
}

// EnvInspector interface //////////////////////////////////////////////////////////////////////////

func (e *env) Report(payload []byte) {
	slog.Info("sending report", "payload", hexutil.Encode(payload))
	err := e.rollup.sendReport(payload)
	if err != nil {
		panic(err)
	}
}

func (e *env) Reportf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	slog.Info("sending report", "payload", message)
	err := e.rollup.sendReport([]byte(message))
	if err != nil {
		panic(err)
	}
}

// EnvInspector interface //////////////////////////////////////////////////////////////////////////

func (e *env) Voucher(destination common.Address, payload []byte) int {
	slog.Info("sending voucher", "destination", destination, "payload", hexutil.Encode(payload))
	index, err := e.rollup.sendVoucher(destination, payload)
	if err != nil {
		panic(err)
	}
	return index
}

func (e *env) Notice(payload []byte) int {
	slog.Info("sending notice", "payload", hexutil.Encode(payload))
	index, err := e.rollup.sendNotice(payload)
	if err != nil {
		panic(err)
	}
	return index
}
