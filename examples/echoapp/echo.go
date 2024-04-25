// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package echoapp

import (
	"github.com/rollmelette/rollmelette"
)

// EchoApplication is an application that emits a voucher, a notice, and a report for each advance
// input; and a report for each inspect input.
type EchoApplication struct{}

func (a *EchoApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	env.Voucher(metadata.MsgSender, payload)
	env.Notice(payload)
	env.Report(payload)
	return nil
}

func (a *EchoApplication) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	env.Report(payload)
	return nil
}
