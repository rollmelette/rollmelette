// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package errorapp

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

var payload = common.Hex2Bytes("deadbeef")
var msgSender = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

type ErrorSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *ErrorSuite) SetupTest() {
	app := new(ErrorApplication)
	s.tester = rollmelette.NewTester(app)
}

func (s *ErrorSuite) TestItRejectsAdvance() {
	result := s.tester.Advance(msgSender, payload)
	s.ErrorContains(result.Err, "input not accepted")
}

func (s *ErrorSuite) TestItRejectsInspect() {
	result := s.tester.Inspect(payload)
	s.ErrorContains(result.Err, "input not accepted")
}
