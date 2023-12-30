// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"log/slog"
	"os"

	"github.com/gligneul/rollmelette"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

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
