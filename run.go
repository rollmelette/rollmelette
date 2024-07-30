// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"context"
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
	opts.AddressBook = NewAddressBook()
	opts.RollupURL = "http://127.0.0.1:5004"
	return &opts
}

// Run connects to the Rollup API and calls the application.
// If opt is nil, this function creates it with the NewRunOpts function.
func Run(ctx context.Context, opts *RunOpts, app Application) (err error) {
	if opts == nil {
		opts = NewRunOpts()
	}
	rollup := newRollupHttp(opts.RollupURL)
	env := newEnv(ctx, opts.AddressBook, rollup, app)
	status := finishStatusAccept
	for {
		input, err := rollup.finishAndGetNext(ctx, status)
		if err != nil {
			return err
		}
		err = env.handle(input)
		if err != nil {
			status = finishStatusReject
		} else {
			status = finishStatusAccept
		}
	}
}
