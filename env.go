// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// env is the implementation of the Env interface.
// This struct isn't public because we don't want users of Rollmelette to create it.
// Instead, it is create by the running and testing functions.
type env struct {
	AddressBook
	ctx         context.Context
	rollup      rollupEnv
	app         Application
	appAddress  common.Address
	etherWallet *etherWallet
	erc20Wallet *erc20Wallet
}

func newEnv(ctx context.Context, addressBook AddressBook, rollup rollupEnv, app Application) *env {
	return &env{
		ctx:         ctx,
		AddressBook: addressBook,
		rollup:      rollup,
		app:         app,
		etherWallet: newEtherWallet(),
		erc20Wallet: newErc20Wallet(),
	}
}

// handlers ////////////////////////////////////////////////////////////////////////////////////////

func (e *env) handle(input any) (err error) {
	defer func() {
		// Recover from panic so we can safely reject the input and print an error message.
		panicObj := recover()
		if panicObj != nil {
			panicErr, ok := panicObj.(error)
			if ok {
				err = panicErr
			} else {
				err = fmt.Errorf("a panic occurred: %v", panicObj)
			}
		}
		if err != nil {
			slog.Error("input rejected", "error", err)
		}
	}()
	switch input := input.(type) {
	case *advanceInput:
		return e.handleAdvance(input)
	case *inspectInput:
		return e.handleInspect(input.Payload)
	default:
		// impossible
		panic("invalid input type")
	}
}

func (e *env) handleAdvance(input *advanceInput) error {
	slog.Debug("received advance",
		"payload", hexutil.Encode(input.Payload),
		"chainId", input.Metadata.ChainId,
		"appContract", input.Metadata.AppContract,
		"inputIndex", input.Metadata.InputIndex,
		"msgSender", input.Metadata.MsgSender,
		"blockNumber", input.Metadata.BlockNumber,
		"blockTimestamp", input.Metadata.BlockTimestamp,
		"prevRandao", input.Metadata.PrevRandao,
	)
	var (
		err     error
		deposit Deposit
		payload []byte = input.Payload
	)
	e.appAddress = (common.Address)(input.Metadata.AppContract)
	switch input.Metadata.MsgSender {
	case e.EtherPortal:
		deposit, payload, err = e.etherWallet.deposit(payload)
	case e.ERC20Portal:
		deposit, payload, err = e.erc20Wallet.deposit(payload)
	}
	if err != nil {
		return err
	}
	if deposit != nil {
		slog.Debug("received deposit", "deposit", deposit)
	}
	return e.app.Advance(e, input.Metadata, deposit, payload)
}

func (e *env) handleInspect(payload []byte) error {
	slog.Debug("received inspect", "payload", hexutil.Encode(payload))
	return e.app.Inspect(e, payload)
}

// EnvInspector interface //////////////////////////////////////////////////////////////////////////

func (e *env) Report(payload []byte) {
	slog.Debug("sending report", "payload", hexutil.Encode(payload))
	err := e.rollup.sendReport(e.ctx, payload)
	if err != nil {
		panic(err)
	}
}

func (e *env) AppAddress() common.Address {
	return e.appAddress
}

func (e *env) EtherAddresses() []common.Address {
	return e.etherWallet.addresses()
}

func (e *env) EtherBalanceOf(address common.Address) *big.Int {
	return e.etherWallet.balanceOf(address)
}

func (e *env) ERC20Tokens() []common.Address {
	return e.erc20Wallet.tokens()
}

func (e *env) ERC20Addresses(token common.Address) []common.Address {
	return e.erc20Wallet.addresses(token)
}

func (e *env) ERC20BalanceOf(token common.Address, address common.Address) *big.Int {
	return e.erc20Wallet.balanceOf(token, address)
}

// Env interface ///////////////////////////////////////////////////////////////////////////////////

func (e *env) Voucher(destination common.Address, value *big.Int, payload []byte) int {
	slog.Debug("sending voucher", "destination", destination, "value", value, "payload", hexutil.Encode(payload))
	index, err := e.rollup.sendVoucher(e.ctx, destination, value, payload)
	if err != nil {
		panic(err)
	}
	return index
}

func (e *env) Notice(payload []byte) int {
	slog.Debug("sending notice", "payload", hexutil.Encode(payload))
	index, err := e.rollup.sendNotice(e.ctx, payload)
	if err != nil {
		panic(err)
	}
	return index
}

func (e *env) EtherTransfer(src common.Address, dst common.Address, value *big.Int) error {
	return e.etherWallet.transfer(src, dst, value)
}

func (e *env) EtherWithdraw(address common.Address, value *big.Int) (int, error) {
	err := e.etherWallet.withdraw(address, value)
	if err != nil {
		return 0, err
	}
	return e.Voucher(e.appAddress, value, nil), nil
}

func (e *env) ERC20Transfer(
	token common.Address,
	src common.Address,
	dst common.Address,
	value *big.Int,
) error {
	return e.erc20Wallet.transfer(token, src, dst, value)
}

func (e *env) ERC20Withdraw(
	token common.Address,
	address common.Address,
	value *big.Int,
) (int, error) {
	payload, err := e.erc20Wallet.withdraw(token, address, value)
	if err != nil {
		return 0, err
	}
	return e.Voucher(token, value, payload), nil
}
