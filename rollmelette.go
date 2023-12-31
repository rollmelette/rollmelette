// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// Rollmelette is a high-level framework for Cartesi Rollups in Go.
package rollmelette

import (
	"log/slog"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

// Application is the interface that should be implemented by the application developer.
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

	// AppAddress returns the application address sent by the address relay contract.
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

// init configures the slog package with the tint handler.
func init() {
	logOpts := new(tint.Options)
	logOpts.Level = slog.LevelDebug
	logOpts.NoColor = !isatty.IsTerminal(os.Stdout.Fd())
	// disable timestamp because it is irrelevant in the cartesi machine
	logOpts.ReplaceAttr = func(groups []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey {
			var zeroattr slog.Attr
			return zeroattr
		}
		return attr
	}
	handler := tint.NewHandler(os.Stdout, logOpts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
