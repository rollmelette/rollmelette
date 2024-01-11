// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ERC20Deposit ////////////////////////////////////////////////////////////////////////////////////

// ERC20Deposit represents an deposit that arrived to the ERC20 wallet.
type ERC20Deposit struct {
	// Token is the address of the ERC20 token.
	Token common.Address

	// Sender is the account that sent the deposit.
	Sender common.Address

	// Amount is the amount of tokens sent.
	Amount *big.Int
}

func (d *ERC20Deposit) String() string {
	value := d.Amount.String()
	return fmt.Sprintf("%v deposited %v of %v token", d.Sender, value, d.Token)
}

// erc20Wallet /////////////////////////////////////////////////////////////////////////////////////

// erc20Wallet is a wallet that manages ERC20 tokens.
type erc20Wallet struct {
	balance map[common.Address]map[common.Address]big.Int
}

func newErc20Wallet() *erc20Wallet {
	return &erc20Wallet{
		balance: make(map[common.Address]map[common.Address]big.Int),
	}
}

func (w *erc20Wallet) tokens() []common.Address {
	var tokens []common.Address
	for t := range w.balance {
		tokens = append(tokens, t)
	}
	sortAddresses(tokens)
	return tokens
}

func (w *erc20Wallet) addresses(token common.Address) []common.Address {
	var addresses []common.Address
	for a := range w.balance[token] {
		addresses = append(addresses, a)
	}
	sortAddresses(addresses)
	return addresses
}

func (w *erc20Wallet) setBalance(token common.Address, address common.Address, value *big.Int) {
	if value.Sign() == 0 {
		if w.balance[token] != nil {
			delete(w.balance[token], address)
			if len(w.balance[token]) == 0 {
				delete(w.balance, token)
			}
		}
	} else {
		if w.balance[token] == nil {
			w.balance[token] = make(map[common.Address]big.Int)
		}
		w.balance[token][address] = *value
	}
}

func (w *erc20Wallet) balanceOf(token common.Address, address common.Address) *big.Int {
	balance := w.balance[token][address]
	return &balance
}

func (w *erc20Wallet) transfer(
	token common.Address,
	src common.Address,
	dst common.Address,
	value *big.Int,
) error {
	if src == dst {
		return fmt.Errorf("can't transfer to self")
	}
	newSrcBalance := new(big.Int).Sub(w.balanceOf(token, src), value)
	if newSrcBalance.Sign() < 0 {
		return fmt.Errorf("insuficient funds")
	}
	newDstBalance := new(big.Int).Add(w.balanceOf(token, dst), value)
	if newDstBalance.Cmp(MaxUint256) > 0 {
		return fmt.Errorf("balance overflow")
	}

	// commit
	w.setBalance(token, src, newSrcBalance)
	w.setBalance(token, dst, newDstBalance)
	return nil
}

func (w *erc20Wallet) withdraw(
	token common.Address,
	address common.Address,
	value *big.Int,
) ([]byte, error) {
	newBalance := new(big.Int).Sub(w.balanceOf(token, address), value)
	if newBalance.Sign() < 0 {
		return nil, fmt.Errorf("insuficient funds")
	}
	w.setBalance(token, address, newBalance)
	return encodeERC20Withdraw(address, value), nil
}

func (w *erc20Wallet) deposit(payload []byte) (Deposit, []byte, error) {
	if len(payload) < 1+20+20+32 {
		return nil, nil, fmt.Errorf("invalid erc20 deposit size; got %v", len(payload))
	}

	// This field will be removed in rollups contracts v2.0
	if payload[0] == 0 {
		return nil, nil, fmt.Errorf("received failed erc20 transfer")
	}
	payload = payload[1:]

	token := common.BytesToAddress(payload[:20])
	payload = payload[20:]

	sender := common.BytesToAddress(payload[:20])
	payload = payload[20:]

	amount := new(big.Int).SetBytes(payload[:32])
	payload = payload[32:]

	newBalance := new(big.Int).Add(w.balanceOf(token, sender), amount)
	if newBalance.Cmp(MaxUint256) > 0 {
		// This should not be possible in real world, but we handle it anyway.
		slog.Warn("overflow erc20 balance", "account", sender)
		newBalance = MaxUint256
	}
	w.setBalance(token, sender, newBalance)

	deposit := &ERC20Deposit{token, sender, amount}
	return deposit, payload, nil
}

// auxiliary functions /////////////////////////////////////////////////////////////////////////////

// encodeERC20Withdraw encodes the voucher to withdraw the asset from the portal.
func encodeERC20Withdraw(token common.Address, value *big.Int) []byte {
	abiJson := `[{
		"type": "function",
		"name": "transfer",
		"inputs": [
			{"type": "address"},
			{"type": "uint256"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		log.Panicf("failed to decode ABI: %v", err)
	}
	voucher, err := abiInterface.Pack("transfer", token, value)
	if err != nil {
		log.Panicf("failed to pack: %v", err)
	}
	return voucher
}
