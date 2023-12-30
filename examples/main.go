// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"log/slog"
	"os"

	"github.com/gligneul/rollmelette"
)

func main() {
	examples := map[string]rollmelette.Application{
		"echo": &EchoApplication{},
	}

	if len(os.Args) < 2 {
		slog.Error("missing example name")
		os.Exit(1)
	}

	app, ok := examples[os.Args[1]]
	if !ok {
		slog.Error("example not found")
		os.Exit(1)
	}

	opts := rollmelette.NewRunOpts()
	err := rollmelette.Run(opts, app)
	if err != nil {
		slog.Error("application exited with error", "error", err)
		os.Exit(1)
	}
}
