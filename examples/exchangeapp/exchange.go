// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// ExchangeApp contains an example of a decentralized exchange implemented with an order book.
package exchangeapp

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/rollmelette"
)

// Advance Inputs //////////////////////////////////////////////////////////////////////////////////

// InputKind is an enum that represents the kind of advance input.
type InputKind string

const (
	// InputAddOrder adds an order to the exchange.
	InputAddOrder InputKind = "AddOrder"

	// InputDeleteOrder deletes an order from the exchange.
	InputDeleteOrder InputKind = "DeleteOrder"

	// InputWithdrawEther withdraws Ether from the exchange.
	InputWithdrawEther InputKind = "WithdrawEther"

	// InputWithdrawToken withdraws the ERC20 token from the exchange.
	InputWithdrawToken InputKind = "WithdrawToken"
)

// Reports /////////////////////////////////////////////////////////////////////////////////////////

// ExchangeApplication is a decentralized exchange that allows exchanging a predetermined ERC20
// token for Ether. The exchange uses JSON to encode inputs and outputs.
type ExchangeApplication struct {
	token common.Address
}

func NewExchangeApplication(token common.Address) *ExchangeApplication {
	return &ExchangeApplication{
		token: token,
	}
}

func (a *ExchangeApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	return nil
}

func (a *ExchangeApplication) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	return nil
}
