// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"fmt"

	"github.com/gligneul/rollmelette"
)

// AddressApplication is an application that receives the app address in an advance input and
// returns it in an inspect input.
type AddressApplication struct{}

func (a *AddressApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	payload []byte,
) error {
	// The app address is obtained automatically by rollmelette; other inputs are rejected.
	return fmt.Errorf("input not accepted")
}

func (a *AddressApplication) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	address, ok := env.AppAddress()
	if ok {
		env.Report(address[:])
	}
	return nil
}
