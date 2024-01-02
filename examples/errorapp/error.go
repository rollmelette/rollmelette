// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package errorapp

import (
	"fmt"

	"github.com/gligneul/rollmelette"
)

// ErrorApplication is an application that always returns an error; rejecting the input.
type ErrorApplication struct{}

func (a *ErrorApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	return fmt.Errorf("input not accepted")
}

func (a *ErrorApplication) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	return fmt.Errorf("input not accepted")
}
