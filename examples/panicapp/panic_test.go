// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package panicapp

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/rollmelette"
	"github.com/stretchr/testify/suite"
)

var payload = common.Hex2Bytes("deadbeef")
var msgSender = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")

func TestPanicSuite(t *testing.T) {
	suite.Run(t, new(PanicSuite))
}

type PanicSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *PanicSuite) SetupTest() {
	app := new(PanicApplication)
	s.tester = rollmelette.NewTester(app)
}

func (s *PanicSuite) TestItRejectsAdvance() {
	result := s.tester.Advance(msgSender, payload)
	s.ErrorContains(result.Err, "a panic occurred: input not accepted")
}

func (s *PanicSuite) TestItRejectsInspect() {
	result := s.tester.Inspect(payload)
	s.ErrorContains(result.Err, "input not accepted")
}
