// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package honeypotapp

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

var owner = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
var hacker = common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
var appAddress = common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e")

func TestHoneypotSuite(t *testing.T) {
	suite.Run(t, new(HoneypotSuite))
}

type HoneypotSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *HoneypotSuite) SetupTest() {
	app := new(HoneypotApplication)
	app.Owner = owner
	s.tester = rollmelette.NewTester(app)
}

func (s *HoneypotSuite) TestItDepositsEther() {
	// from owner
	result := s.tester.DepositEther(owner, big.NewInt(100), nil)
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	s.checkBalance(100, result.Reports[0].Payload)

	// from third party
	result = s.tester.DepositEther(hacker, big.NewInt(100), nil)
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	s.checkBalance(200, result.Reports[0].Payload)
}

func (s *HoneypotSuite) TestItWithdrawsEther() {
	// deposit
	result := s.tester.DepositEther(owner, big.NewInt(100), nil)
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	s.checkBalance(100, result.Reports[0].Payload)

	// withdraw
	result = s.tester.Advance(owner, nil)
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	s.checkBalance(0, result.Reports[0].Payload)
	s.Len(result.Vouchers, 1)
	s.Equal(appAddress, result.Vouchers[0].Destination)

	// check voucher
	expectedVoucher := make([]byte, 0, 4+32+32)
	expectedVoucher = append(expectedVoucher, 0x52, 0x2f, 0x68, 0x15)
	expectedVoucher = append(expectedVoucher, make([]byte, 12)...) // padding
	expectedVoucher = append(expectedVoucher, owner[:]...)
	expectedVoucher = append(expectedVoucher, big.NewInt(100).FillBytes(make([]byte, 32))...)
	s.Equal(expectedVoucher, result.Vouchers[0].Payload)
}

func (s *HoneypotSuite) TestItFailsToWithdrawWithoutFunds() {
	result := s.tester.Advance(owner, nil)
	s.ErrorContains(result.Err, "nothing to withdraw")
}

func (s *HoneypotSuite) TestItFailsToWithdrawFromHacker() {
	// deposit
	result := s.tester.DepositEther(owner, big.NewInt(100), nil)
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	s.checkBalance(100, result.Reports[0].Payload)

	// withdraw
	result = s.tester.Advance(hacker, nil)
	s.ErrorContains(result.Err, "input not from owner")
	s.Len(result.Reports, 1)
	s.checkBalance(100, result.Reports[0].Payload)
}

func (s *HoneypotSuite) TestItReportsBalanceInInspect() {
	result := s.tester.Inspect(nil)
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	s.checkBalance(0, result.Reports[0].Payload)

	advanceResult := s.tester.DepositEther(owner, big.NewInt(100), nil)
	s.Nil(advanceResult.Err)

	result = s.tester.Inspect(nil)
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	s.checkBalance(100, result.Reports[0].Payload)
}

func (s *HoneypotSuite) checkBalance(expected int64, payload []byte) {
	balance := new(big.Int).SetBytes(payload)
	s.Zerof(balance.Cmp(big.NewInt(expected)), "expected %v; got %v", expected, balance)
}