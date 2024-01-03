// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package panicapp

import (
	"fmt"

	"github.com/gligneul/rollmelette"
)

// PanicApplication is an application that always panic. Rollmelette captures the panic and reject
// the input.
type PanicApplication struct{}

func (a *PanicApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	panic("input not accepted")
}

func (a *PanicApplication) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	panic(fmt.Errorf("input not accepted"))
}
