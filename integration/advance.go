// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package integration

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rollmelette/rollmelette"
)

const ApplicationAddress = "0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e"
const SenderAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
const TestMnemonic = "test test test test test test test test test test test junk"

// Advance sends an input using the devnet contract addresses.
func Advance(ctx context.Context, url string, payload []byte) error {
	if len(payload) == 0 {
		return fmt.Errorf("cannot send empty payload")
	}
	book := rollmelette.NewAddressBook()
	input := hexutil.Encode(payload)
	cmd := exec.CommandContext(ctx,
		"cast", "send",
		"--mnemonic", TestMnemonic,
		"--rpc-url", url,
		book.InputBox.String(),             // TO
		"addInput(address,bytes)(bytes32)", // SIG
		ApplicationAddress, input,          // ARGS
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cast: %w: %v", err, string(output))
	}
	return nil
}
