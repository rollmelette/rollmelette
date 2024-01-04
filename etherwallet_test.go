// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

func TestEtherWalletSuite(t *testing.T) {
	suite.Run(t, new(EtherWalletSuite))
}

type EtherWalletSuite struct {
	suite.Suite
	wallet *etherWallet
	src    common.Address
	dst    common.Address
}

func (s *EtherWalletSuite) SetupTest() {
	s.wallet = newEtherWallet()
	s.src = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	s.dst = common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
}

func (s *EtherWalletSuite) TestDepositString() {
	value := big.NewInt(123000000000000000)
	deposit := &EtherDeposit{s.src, value}
	expected := "0xFafafAfafAFaFAFaFafafafAfaFaFAfAfAfAFaFA deposited 0.123000000000000000 Ether"
	s.Equal(expected, deposit.String())
}

func (s *EtherWalletSuite) TestAddresses() {
	addresses := s.wallet.addresses()
	s.Empty(addresses)

	s.wallet.setBalance(s.src, big.NewInt(1))
	addresses = s.wallet.addresses()
	s.Equal(addresses, []common.Address{s.src})

	s.wallet.setBalance(s.dst, big.NewInt(1))
	addresses = s.wallet.addresses()
	s.Equal(addresses, []common.Address{s.src, s.dst})

	s.wallet.setBalance(s.src, big.NewInt(0))
	s.wallet.setBalance(s.dst, big.NewInt(0))
	addresses = s.wallet.addresses()
	s.Empty(addresses)
}

func (s *EtherWalletSuite) TestBalanceOf() {
	balance := s.wallet.balanceOf(s.src)
	s.Equal(big.NewInt(0), balance)

	s.wallet.setBalance(s.src, big.NewInt(50))
	balance = s.wallet.balanceOf(s.src)
	s.Equal(big.NewInt(50), balance)
}

func (s *EtherWalletSuite) TestValidTransfer() {
	s.wallet.setBalance(s.src, big.NewInt(50))
	s.wallet.setBalance(s.dst, big.NewInt(50))
	err := s.wallet.transfer(s.src, s.dst, big.NewInt(50))
	s.Nil(err)
	srcBalance := s.wallet.balanceOf(s.src)
	s.Equal(big.NewInt(0), srcBalance)
	dstBalance := s.wallet.balanceOf(s.dst)
	s.Equal(big.NewInt(100), dstBalance)
}

func (s *EtherWalletSuite) TestZeroTransfer() {
	err := s.wallet.transfer(s.src, s.dst, big.NewInt(0))
	s.Nil(err)
	srcBalance := s.wallet.balanceOf(s.src)
	s.Equal(big.NewInt(0), srcBalance)
	dstBalance := s.wallet.balanceOf(s.dst)
	s.Equal(big.NewInt(0), dstBalance)
}

func (s *EtherWalletSuite) TestSelfTransfer() {
	s.wallet.setBalance(s.src, big.NewInt(50))
	err := s.wallet.transfer(s.src, s.src, big.NewInt(50))
	s.ErrorContains(err, "can't transfer to self")
}

func (s *EtherWalletSuite) TestInsuficientFundsTransfer() {
	s.wallet.setBalance(s.src, big.NewInt(50))
	err := s.wallet.transfer(s.src, s.dst, big.NewInt(100))
	s.ErrorContains(err, "insuficient funds")
}

func (s *EtherWalletSuite) TestBalanceOverflowTransfer() {
	s.wallet.setBalance(s.src, big.NewInt(50))
	s.wallet.setBalance(s.dst, MaxUint256)
	err := s.wallet.transfer(s.src, s.dst, big.NewInt(50))
	s.ErrorContains(err, "balance overflow")
}

func (s *EtherWalletSuite) TestEncodeWithdraw() {
	voucher := encodeEtherWithdraw(s.src, big.NewInt(100))
	expected := common.Hex2Bytes("522f6815000000000000000000000000fafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064")
	s.Equal(expected, voucher)
}

func (s *EtherWalletSuite) TestInsuficientFundsWithdraw() {
	s.wallet.setBalance(s.src, big.NewInt(50))
	voucher, err := s.wallet.withdraw(s.src, big.NewInt(100))
	s.Nil(voucher)
	s.ErrorContains(err, "insuficient funds")
	balance := s.wallet.balanceOf(s.src)
	s.Equal(big.NewInt(50), balance)
}

func (s *EtherWalletSuite) TestValidWithdraw() {
	s.wallet.setBalance(s.src, big.NewInt(100))
	voucher, err := s.wallet.withdraw(s.src, big.NewInt(100))
	s.Nil(err)
	s.NotNil(voucher)
	balance := s.wallet.balanceOf(s.src)
	s.Equal(big.NewInt(0), balance)
}

func (s *EtherWalletSuite) TestValidDeposit() {
	payload := common.Hex2Bytes("fafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064deadbeef")
	deposit, input, err := s.wallet.deposit(payload)
	s.Nil(err)
	etherDeposit := deposit.(*EtherDeposit)
	s.Equal(s.src, etherDeposit.Sender)
	s.Equal(big.NewInt(100), etherDeposit.Value)
	s.Equal(common.Hex2Bytes("deadbeef"), input)
	balance := s.wallet.balanceOf(s.src)
	s.Equal(big.NewInt(100), balance)
}

func (s *EtherWalletSuite) TestValidDepositWithEmptyInput() {
	payload := common.Hex2Bytes("fafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064")
	deposit, input, err := s.wallet.deposit(payload)
	s.Nil(err)
	s.NotNil(deposit)
	s.Empty(input)
}

func (s *EtherWalletSuite) TestOverflowDeposit() {
	// deposit int max
	payload := common.Hex2Bytes("fafafafafafafafafafafafafafafafafafafafaffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	deposit, input, err := s.wallet.deposit(payload)
	s.Nil(err)
	s.NotNil(deposit)
	s.Empty(input)

	// deposit more ether
	payload = common.Hex2Bytes("fafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000001000")
	deposit, input, err = s.wallet.deposit(payload)
	s.Nil(err)
	s.NotNil(deposit)
	s.Empty(input)

	// check balance
	balance := s.wallet.balanceOf(s.src)
	s.Equal(MaxUint256, balance)
}

func (s *EtherWalletSuite) TestMalformedDeposit() {
	payload := common.Hex2Bytes("fafafa")
	_, _, err := s.wallet.deposit(payload)
	s.ErrorContains(err, "invalid eth deposit size; got 3")
}
