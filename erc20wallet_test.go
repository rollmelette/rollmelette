// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

func TestERC20WalletSuite(t *testing.T) {
	suite.Run(t, new(ERC20WalletSuite))
}

type ERC20WalletSuite struct {
	suite.Suite
	wallet *erc20Wallet
	tokens []common.Address
	src    common.Address
	dst    common.Address
}

func (s *ERC20WalletSuite) SetupTest() {
	s.wallet = newErc20Wallet()
	s.tokens = []common.Address{
		common.HexToAddress("0xbabababababababababababababababababababa"),
		common.HexToAddress("0xbebebebebebebebebebebebebebebebebebebebe"),
	}
	s.src = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	s.dst = common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
}

func (s *ERC20WalletSuite) TestDepositString() {
	value := big.NewInt(123)
	deposit := &ERC20Deposit{s.tokens[0], s.src, value}
	expectedString := "0xFafafAfafAFaFAFaFafafafAfaFaFAfAfAfAFaFA deposited 123 of " +
		"0xBAbAbabAbabaBABaBAbABabaBAbAbaBaBAbABaBa token"
	s.Equal(expectedString, deposit.String())
}

func (s *ERC20WalletSuite) TestTokens() {
	// test zero tokens
	tokens := s.wallet.tokens()
	s.Empty(tokens)

	// test single token
	s.wallet.setBalance(s.tokens[0], s.src, big.NewInt(1))
	tokens = s.wallet.tokens()
	expectedTokens := []common.Address{s.tokens[0]}
	s.Equal(expectedTokens, tokens)

	// test two tokens
	s.wallet.setBalance(s.tokens[1], s.dst, big.NewInt(1))
	tokens = s.wallet.tokens()
	s.Equal(s.tokens, tokens)
}

func (s *ERC20WalletSuite) TestAddresses() {
	// test zero addresses
	addresses := s.wallet.addresses(s.tokens[0])
	s.Empty(addresses)

	// test single address
	s.wallet.setBalance(s.tokens[0], s.src, big.NewInt(1))
	addresses = s.wallet.addresses(s.tokens[0])
	expected := []common.Address{s.src}
	s.Equal(expected, addresses)

	// test two addresses
	s.wallet.setBalance(s.tokens[0], s.dst, big.NewInt(1))
	addresses = s.wallet.addresses(s.tokens[0])
	expected = []common.Address{s.src, s.dst}
	s.Equal(expected, addresses)

	// test adding balance to another token
	s.wallet.setBalance(s.tokens[1], s.src, big.NewInt(1))
	addresses = s.wallet.addresses(s.tokens[1])
	expected = []common.Address{s.src}
	s.Equal(expected, addresses)
}

func (s *ERC20WalletSuite) TestBalanceOf() {
	// test zero balance
	balance := s.wallet.balanceOf(s.tokens[0], s.src)
	s.Equal(big.NewInt(0), balance)

	// test non-zero balance
	s.wallet.setBalance(s.tokens[0], s.src, big.NewInt(50))
	balance = s.wallet.balanceOf(s.tokens[0], s.src)
	s.Equal(big.NewInt(50), balance)

	// test adding balance to another token
	s.wallet.setBalance(s.tokens[1], s.src, big.NewInt(100))
	balance = s.wallet.balanceOf(s.tokens[0], s.src)
	s.Equal(big.NewInt(50), balance)

	// test setting balance to zero
	s.wallet.setBalance(s.tokens[0], s.src, big.NewInt(0))
	balance = s.wallet.balanceOf(s.tokens[0], s.src)
	s.Equal(big.NewInt(0), balance)
}

func (s *ERC20WalletSuite) TestValidTransfer() {
	s.wallet.setBalance(s.tokens[0], s.src, big.NewInt(50))
	s.wallet.setBalance(s.tokens[0], s.dst, big.NewInt(50))
	err := s.wallet.transfer(s.tokens[0], s.src, s.dst, big.NewInt(50))
	s.Nil(err)
	srcBalance := s.wallet.balanceOf(s.tokens[0], s.src)
	s.Equal(big.NewInt(0), srcBalance)
	dstBalance := s.wallet.balanceOf(s.tokens[0], s.dst)
	s.Equal(big.NewInt(100), dstBalance)
}

func (s *ERC20WalletSuite) TestZeroTransfer() {
	err := s.wallet.transfer(s.tokens[0], s.src, s.dst, big.NewInt(0))
	s.Nil(err)
	srcBalance := s.wallet.balanceOf(s.tokens[0], s.src)
	s.Equal(big.NewInt(0), srcBalance)
	dstBalance := s.wallet.balanceOf(s.tokens[0], s.dst)
	s.Equal(big.NewInt(0), dstBalance)
}

func (s *ERC20WalletSuite) TestSelfTransfer() {
	s.wallet.setBalance(s.tokens[0], s.src, big.NewInt(50))
	err := s.wallet.transfer(s.tokens[0], s.src, s.src, big.NewInt(50))
	s.ErrorContains(err, "can't transfer to self")
}

func (s *ERC20WalletSuite) TestInsuficientFundsTransfer() {
	s.wallet.setBalance(s.tokens[0], s.src, big.NewInt(50))
	err := s.wallet.transfer(s.tokens[0], s.src, s.dst, big.NewInt(100))
	s.ErrorContains(err, "insuficient funds")
}

func (s *ERC20WalletSuite) TestBalanceOverflowTransfer() {
	s.wallet.setBalance(s.tokens[0], s.src, big.NewInt(50))
	s.wallet.setBalance(s.tokens[0], s.dst, MaxUint256)
	err := s.wallet.transfer(s.tokens[0], s.src, s.dst, big.NewInt(50))
	s.ErrorContains(err, "balance overflow")
}

func (s *ERC20WalletSuite) TestValidWithdraw() {
	s.wallet.setBalance(s.tokens[0], s.src, big.NewInt(100))
	voucher, err := s.wallet.withdraw(s.tokens[0], s.src, big.NewInt(100))
	s.Nil(err)
	expected := common.Hex2Bytes("a9059cbb000000000000000000000000fafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064")
	s.Equal(expected, voucher)
	balance := s.wallet.balanceOf(s.tokens[0], s.src)
	s.Equal(big.NewInt(0), balance)
}

func (s *ERC20WalletSuite) TestInsuficientFundsWithdraw() {
	s.wallet.setBalance(s.tokens[0], s.src, big.NewInt(50))
	_, err := s.wallet.withdraw(s.tokens[0], s.src, big.NewInt(100))
	s.ErrorContains(err, "insuficient funds")
	balance := s.wallet.balanceOf(s.tokens[0], s.src)
	s.Equal(big.NewInt(50), balance)
}

func (s *ERC20WalletSuite) TestValidDeposit() {
	payload := common.Hex2Bytes("01babababababababababababababababababababafafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064deadbeef")
	deposit, input, err := s.wallet.deposit(payload)
	s.Nil(err)

	// check deposit
	erc20Deposit, ok := deposit.(*ERC20Deposit)
	s.Require().True(ok)
	s.Equal(s.tokens[0], erc20Deposit.Token)
	s.Equal(s.src, erc20Deposit.Sender)
	s.Equal(big.NewInt(100), erc20Deposit.Amount)

	// check input data
	s.Equal(common.Hex2Bytes("deadbeef"), input)

	// check balance
	balance := s.wallet.balanceOf(s.tokens[0], s.src)
	s.Equal(big.NewInt(100), balance)
}

func (s *ERC20WalletSuite) TestValidDepositWithEmptyInput() {
	payload := common.Hex2Bytes("01babababababababababababababababababababafafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064")
	deposit, input, err := s.wallet.deposit(payload)
	s.Nil(err)
	s.NotNil(deposit)
	s.Empty(input)
}

func (s *ERC20WalletSuite) TestOverflowDeposit() {
	// deposit int max
	payload := common.Hex2Bytes("01babababababababababababababababababababafafafafafafafafafafafafafafafafafafafafaffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_, _, err := s.wallet.deposit(payload)
	s.Nil(err)

	// deposit more ether
	payload = common.Hex2Bytes("01babababababababababababababababababababafafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000001000")
	_, _, err = s.wallet.deposit(payload)
	s.Nil(err)

	// check balance
	balance := s.wallet.balanceOf(s.tokens[0], s.src)
	s.Equal(MaxUint256, balance)
}

func (s *ERC20WalletSuite) TestMalformedDeposit() {
	payload := common.Hex2Bytes("fafafa")
	_, _, err := s.wallet.deposit(payload)
	s.ErrorContains(err, "invalid erc20 deposit size; got 3")
}

func (s *ERC20WalletSuite) TestFailedDeposit() {
	payload := common.Hex2Bytes("00babababababababababababababababababababafafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064")
	_, _, err := s.wallet.deposit(payload)
	s.ErrorContains(err, "received failed erc20 transfer")
}
