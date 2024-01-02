// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package addressapp

import (
	"testing"

	"github.com/gligneul/rollmelette"
	"github.com/stretchr/testify/suite"
)

func TestAddressSuite(t *testing.T) {
	suite.Run(t, new(AddressSuite))
}

type AddressSuite struct {
	suite.Suite
	app    *AddressApplication
	tester *rollmelette.Tester
}

func (s *AddressSuite) SetupTest() {
	s.app = new(AddressApplication)
	s.tester = rollmelette.NewTester(s.app)
}

func (s *AddressSuite) TestItAcceptsTheAppAddress() {
	// Get nothing before sending address
	inspectResult := s.tester.Inspect(nil)
	s.Nil(inspectResult.Err)
	s.Empty(inspectResult.Reports)

	// Send address
	advanceResult := s.tester.RelayAppAddress()
	s.Nil(advanceResult.Err)

	// Get address from inspect
	inspectResult = s.tester.Inspect(nil)
	s.Nil(inspectResult.Err)
	s.Len(inspectResult.Reports, 1)
	s.Equal(s.tester.AppAddress[:], inspectResult.Reports[0].Payload)
}
