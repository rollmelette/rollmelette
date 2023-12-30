// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"github.com/gligneul/rollmelette"
)

// EchoApplication is an application that emits a voucher, a notice, and a report for each advance
// input; and a report for each inspect input.
type EchoApplication struct{}

func (e *EchoApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	payload []byte,
) error {
	env.Voucher(metadata.MsgSender, payload)
	env.Notice(payload)
	env.Report(payload)
	return nil
}

func (e *EchoApplication) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	env.Report(payload)
	return nil
}
