// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"fmt"
	"log/slog"
	"runtime"
)

// RunOpts allows the application developer to pass some parameters to the run function.
type RunOpts struct {
	AddressBook

	// RollupURL is the URL of the Rollup API.
	RollupURL string
}

// NewRunOpts creates a RunOpts struct with sensible default values.
func NewRunOpts() *RunOpts {
	var opts RunOpts
	if runtime.GOARCH == "riscv64" {
		opts.RollupURL = "http://127.0.0.1:5004"
	} else {
		opts.RollupURL = "http://127.0.0.1:8080/rollup"
	}
	return &opts
}

// Run connects to the Rollup API and calls the application.
// If opt is nil, this function creates it with the NewRunOpts function.
func Run(opts *RunOpts, app Application) (err error) {
	defer func() {
		panicObj := recover()
		if panicObj != nil {
			panicErr, ok := panicObj.(error)
			if ok {
				err = panicErr
			} else {
				err = fmt.Errorf("a panic occured: %v", panicObj)
			}
		}
	}()
	if opts == nil {
		opts = NewRunOpts()
	}
	rollup := newRollupHttp(opts.RollupURL)
	env := newEnv(opts.AddressBook, rollup)
	status := finishStatusAccept
	for {
		input, err := rollup.finishAndGetNext(status)
		if err != nil {
			return err
		}
		switch input := input.(type) {
		case *advanceInput:
			err = env.handleAdvance(app, input)
		case *inspectInput:
			err = env.handleInspect(app, input.Payload)
		default:
			// impossible
			panic("invalid input type")
		}
		if err != nil {
			slog.Error("rejecting input", "error", err)
			status = finishStatusReject
		} else {
			status = finishStatusAccept
		}
	}
}
