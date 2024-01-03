// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package addressapp

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/rollmelette"
	"github.com/stretchr/testify/suite"
)

var msgSender = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
var appAddress = common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")

func TestAddressSuite(t *testing.T) {
	suite.Run(t, new(AddressSuite))
}

type AddressSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *AddressSuite) SetupTest() {
	app := new(AddressApplication)
	s.tester = rollmelette.NewTester(app)
}

func (s *AddressSuite) TestItRejectsAdvance() {
	result := s.tester.Advance(msgSender, nil)
	s.ErrorContains(result.Err, "reject")
}

func (s *AddressSuite) TestItAcceptsTheAppAddress() {
	// Get nothing before sending address
	inspectResult := s.tester.Inspect(nil)
	s.Nil(inspectResult.Err)
	s.Empty(inspectResult.Reports)

	// Send address
	advanceResult := s.tester.RelayAppAddress(appAddress)
	s.Nil(advanceResult.Err)

	// Get address from inspect
	inspectResult = s.tester.Inspect(nil)
	s.Nil(inspectResult.Err)
	s.Len(inspectResult.Reports, 1)
	s.Equal(appAddress[:], inspectResult.Reports[0].Payload)
}
