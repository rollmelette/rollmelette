// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package panicapp

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/rollmelette"
	"github.com/stretchr/testify/suite"
)

func TestPanicSuite(t *testing.T) {
	suite.Run(t, new(PanicSuite))
}

type PanicSuite struct {
	suite.Suite
	app    *PanicApplication
	tester *rollmelette.Tester
}

func (s *PanicSuite) SetupTest() {
	s.app = new(PanicApplication)
	s.tester = rollmelette.NewTester(s.app)
}

func (s *PanicSuite) TestItRejectsAdvance() {
	payload := common.Hex2Bytes("deadbeef")
	result := s.tester.Advance(payload)
	s.ErrorContains(result.Err, "a panic occurred: input not accepted")
}

func (s *PanicSuite) TestItRejectsInspect() {
	payload := common.Hex2Bytes("deadbeef")
	result := s.tester.Inspect(payload)
	s.ErrorContains(result.Err, "a panic occurred: input not accepted")
}
