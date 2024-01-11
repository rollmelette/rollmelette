// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// EtherDeposit ////////////////////////////////////////////////////////////////////////////////////

// EtherDeposit represents an deposit that arrived to the Ether wallet.
type EtherDeposit struct {
	// Sender is the account that sent the deposit.
	Sender common.Address

	// Value is the amount of Wei deposited.
	Value *big.Int
}

func (d *EtherDeposit) String() string {
	value := etherString(d.Value)
	return fmt.Sprintf("%v deposited %v Ether", d.Sender, value)
}

// etherWallet /////////////////////////////////////////////////////////////////////////////////////

// etherWallet is a wallet that manages Ether deposits.
type etherWallet struct {
	balance map[common.Address]big.Int
}

func newEtherWallet() *etherWallet {
	return &etherWallet{
		balance: make(map[common.Address]big.Int),
	}
}

func (w *etherWallet) addresses() []common.Address {
	var addresses []common.Address
	for address := range w.balance {
		addresses = append(addresses, address)
	}
	sortAddresses(addresses)
	return addresses
}

func (w *etherWallet) setBalance(address common.Address, value *big.Int) {
	if value.Sign() == 0 {
		delete(w.balance, address)
	} else {
		w.balance[address] = *value
	}
}

func (w *etherWallet) balanceOf(address common.Address) *big.Int {
	balance := w.balance[address]
	return &balance
}

func (w *etherWallet) deposit(payload []byte) (Deposit, []byte, error) {
	if len(payload) < 20+32 {
		return nil, nil, fmt.Errorf("invalid eth deposit size; got %v", len(payload))
	}

	sender := common.BytesToAddress(payload[:20])
	payload = payload[20:]

	value := new(big.Int).SetBytes(payload[:32])
	payload = payload[32:]

	newBalance := new(big.Int).Add(w.balanceOf(sender), value)
	if newBalance.Cmp(MaxUint256) > 0 {
		// This should not be possible in real world, but we handle it anyway.
		slog.Warn("overflow ether balance", "account", sender)
		newBalance = MaxUint256
	}
	w.setBalance(sender, newBalance)

	deposit := &EtherDeposit{sender, value}
	return deposit, payload, nil
}

func (w *etherWallet) transfer(src common.Address, dst common.Address, value *big.Int) error {
	if src == dst {
		return fmt.Errorf("can't transfer to self")
	}

	newSrcBalance := new(big.Int).Sub(w.balanceOf(src), value)
	if newSrcBalance.Sign() < 0 {
		return fmt.Errorf("insuficient funds")
	}

	newDstBalance := new(big.Int).Add(w.balanceOf(dst), value)
	if newDstBalance.Cmp(MaxUint256) > 0 {
		return fmt.Errorf("balance overflow")
	}

	// commit
	w.setBalance(src, newSrcBalance)
	w.setBalance(dst, newDstBalance)
	return nil
}

func (w *etherWallet) withdraw(address common.Address, value *big.Int) ([]byte, error) {
	newBalance := new(big.Int).Sub(w.balanceOf(address), value)
	if newBalance.Sign() < 0 {
		return nil, fmt.Errorf("insuficient funds")
	}
	w.setBalance(address, newBalance)
	return encodeEtherWithdraw(address, value), nil
}

// auxiliary functions /////////////////////////////////////////////////////////////////////////////

// etherString generates a string with the Ether value given the Wei value.
func etherString(wei *big.Int) string {
	weiFloat := new(big.Float).SetInt(wei)
	tenToEighteen := new(big.Float).SetFloat64(1e18)
	etherFloat := new(big.Float).Quo(weiFloat, tenToEighteen)
	return etherFloat.Text('f', 18)
}

// encodeEtherWithdraw encodes the voucher to withdraw the asset from the portal.
func encodeEtherWithdraw(address common.Address, value *big.Int) []byte {
	abiJson := `[{
		"type": "function",
		"name": "withdrawEther",
		"inputs": [
			{"type": "address"},
			{"type": "uint256"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		log.Panicf("failed to decode ABI: %v", err)
	}
	voucher, err := abiInterface.Pack("withdrawEther", address, value)
	if err != nil {
		log.Panicf("failed to pack: %v", err)
	}
	return voucher
}

// sortAddresses sorts a slice of addresses.
func sortAddresses(addresses []common.Address) {
	slices.SortFunc(addresses, func(a common.Address, b common.Address) int {
		for i := 0; i < len(a); i++ {
			if a[i] < b[i] {
				return -1
			} else if a[i] > b[i] {
				return 1
			}
		}
		return 0
	})
}
