// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/rollmelette"
	"github.com/stretchr/testify/suite"
)

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

type ErrorSuite struct {
	suite.Suite
	app    *ErrorApplication
	tester *rollmelette.Tester
}

func (s *ErrorSuite) SetupTest() {
	s.app = new(ErrorApplication)
	s.tester = rollmelette.NewTester(s.app)
}

func (s *ErrorSuite) TestItRejectsAdvance() {
	payload := common.Hex2Bytes("deadbeef")
	result := s.tester.Advance(payload)
	s.ErrorContains(result.Err, "input not accepted")
}

func (s *ErrorSuite) TestItRejectsInspect() {
	payload := common.Hex2Bytes("deadbeef")
	result := s.tester.Inspect(payload)
	s.ErrorContains(result.Err, "input not accepted")
}
