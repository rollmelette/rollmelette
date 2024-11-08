// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/rollmelette/rollmelette/examples/echoapp"
	"github.com/rollmelette/rollmelette/examples/errorapp"
	"github.com/rollmelette/rollmelette/examples/honeypotapp"
	"github.com/rollmelette/rollmelette/examples/panicapp"
)

func main() {
	examples := map[string]rollmelette.Application{
		"echo":  &echoapp.EchoApplication{},
		"error": &errorapp.ErrorApplication{},
		"honeypot": &honeypotapp.HoneypotApplication{
			Owner: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		},
		"panic": &panicapp.PanicApplication{},
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
	err := rollmelette.Run(context.Background(), opts, app)
	if err != nil {
		slog.Error("application exited with error", "error", err)
		os.Exit(1)
	}
}
