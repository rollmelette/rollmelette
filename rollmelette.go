// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// Rollmelette is a high-level framework for Cartesi Rollups in Go.
package rollmelette

import "github.com/ethereum/go-ethereum/common"

// Application is the interface that should be implemented by the DApp developer.
// The application has only two methods: one for the advance request and another for the inspect
// request.
type Application interface {

	// Advance the application state.
	// If this method returns an error, Rollmelette reverts the execution.
	Advance(env Env, metadata Metadata, payload []byte) error

	// Inspect the application state.
	// If this method returns an error, Rollmelette reverts the execution.
	Inspect(env EnvInspector, payload []byte) error
}

// EnvIspector is the entrypoint for the inspect functions of the Rollup API.
type EnvInspector interface {

	// Report sends a report.
	Report(payload []byte)

	// Reportf receives a format string and arguments, and send them as a report.
	Reportf(format string, args ...any)

	// DAppAddress returns the application address sent by the address relay contract.
	// If the contract didn't send the address yet, the function returns false.
	AppAddress() (common.Address, bool)
}

// Env is the entrypoint for the Rollup API and to Rollmelette's asset management.
type Env interface {
	EnvInspector

	// Voucher sends a voucher and returns its index.
	Voucher(destination common.Address, payload []byte) int

	// Notice sends a notice and returns its index.
	Notice(payload []byte) int
}
