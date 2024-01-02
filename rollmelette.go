// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// Rollmelette is a high-level framework for Cartesi Rollups in Go.
package rollmelette

import (
	"fmt"
	"log/slog"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

// MaxUint256 is the max value for uint256.
var MaxUint256 *big.Int

func init() {
	one := big.NewInt(1)
	// Left shift by 256 bits and then subtract 1 to get the max value of uint256.
	MaxUint256 = new(big.Int).Sub(new(big.Int).Lsh(one, 256), one)
}

// Deposit represents an asset deposit to a portal.
type Deposit interface {
	fmt.Stringer

	// GetSender returns the account that made the deposit.
	GetSender() common.Address
}

// Application is the interface that should be implemented by the application developer.
// The application has one method for the advance request and another for the inspect request.
type Application interface {

	// Advance is called to advance the application state. It receives env to interact with the
	// rollup environment, the advance input metadata, a interface representing the deposit if
	// the input came from a portal, and the input payload. If this method returns an error,
	// Rollmelette reverts the execution.
	Advance(env Env, metadata Metadata, deposit Deposit, payload []byte) error

	// Inspect is called to inspect the application state. It receives env to read the rollup
	// environment and the input payload. If this method returns an error, Rollmelette reverts
	// the execution.
	Inspect(env EnvInspector, payload []byte) error
}

// EnvIspector is the entrypoint for the inspect functions of the Rollup API.
type EnvInspector interface {

	// Report sends a report.
	Report(payload []byte)

	// AppAddress returns the application address sent by the address relay contract.
	// If the contract didn't send the address yet, the function returns false.
	AppAddress() (common.Address, bool)

	// EtherAddresses returns the list of addresses that have Ether.
	EtherAddresses() []common.Address

	// EtherBalanceOf returns the balance of the given address.
	EtherBalanceOf(address common.Address) *big.Int
}

// Env is the entrypoint for the Rollup API and to Rollmelette's asset management.
type Env interface {
	EnvInspector

	// Voucher sends a voucher and returns its index.
	Voucher(destination common.Address, payload []byte) int

	// Notice sends a notice and returns its index.
	Notice(payload []byte) int

	// EtherTransfer transfer the given amount of funds from source to destination.
	// It returns an error if source doesn't have enough funds.
	EtherTransfer(src common.Address, dst common.Address, value *big.Int) error

	// EtherWithdraw withdraws the asset from the wallet, generates the voucher to withdraw
	// it from the portal, and returns the voucher index.
	// Before withdrawing Ether, the application must receive its contract address from the
	// address relay contract.
	// It returns an error if the address doesn't have enough funds.
	EtherWithdraw(address common.Address, value *big.Int) (int, error)
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
