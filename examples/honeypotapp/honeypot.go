// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package honeypotapp

import (
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/rollmelette"
)

// HoneypotApplication is a honeypot[1] application that stores Ether.
// Only the application owner can withdraw from the honeypot.
// The application emits a report with the current balance for every input.
// [1] https://en.wikipedia.org/wiki/Honeypot_(computing)
type HoneypotApplication struct {
	Owner common.Address
}

func (a *HoneypotApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	var err error
	if deposit != nil {
		err = a.deposit(env, deposit)
	} else {
		err = a.withdraw(env, metadata)
	}
	a.Inspect(env, nil)
	return err
}

func (a *HoneypotApplication) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	balance := env.EtherBalanceOf(a.Owner)
	env.Report(balance.FillBytes(make([]byte, 32)))
	return nil
}

func (a *HoneypotApplication) deposit(env rollmelette.Env, deposit rollmelette.Deposit) error {
	switch deposit := deposit.(type) {
	case *rollmelette.EtherDeposit:
		if deposit.Sender != a.Owner {
			env.EtherTransfer(deposit.Sender, a.Owner, deposit.Value)
		}
		return nil
	default:
		return fmt.Errorf("unsupported deposit: %T", deposit)
	}
}

func (a *HoneypotApplication) withdraw(env rollmelette.Env, metadata rollmelette.Metadata) error {
	if metadata.MsgSender != a.Owner {
		return fmt.Errorf("input not from owner")
	}
	balance := env.EtherBalanceOf(a.Owner)
	if balance.Sign() == 0 {
		return fmt.Errorf("nothing to withdraw")
	}
	_, err := env.EtherWithdraw(a.Owner, balance)
	if err != nil {
		return err
	}
	slog.Info("withdrawn", "value", balance)
	return nil
}
