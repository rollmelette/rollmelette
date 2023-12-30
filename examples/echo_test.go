// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/rollmelette"
	"github.com/stretchr/testify/suite"
)

func TestEchoSuite(t *testing.T) {
	suite.Run(t, new(EchoSuite))
}

type EchoSuite struct {
	suite.Suite
	app    *EchoApplication
	tester *rollmelette.Tester
}

func (s *EchoSuite) SetupTest() {
	s.app = new(EchoApplication)
	s.tester = rollmelette.NewTester(s.app)
}

func (s *EchoSuite) TestAdvance() {
	payload := common.Hex2Bytes("deadbeef")
	result := s.tester.Advance(payload)
	s.Nil(result.Err)
	s.Len(result.Vouchers, 1)
	s.Len(result.Notices, 1)
	s.Len(result.Reports, 1)
	s.Equal(payload, result.Vouchers[0].Payload)
	s.Equal(result.MsgSender, result.Vouchers[0].Destination)
	s.Equal(payload, result.Notices[0].Payload)
	s.Equal(payload, result.Reports[0].Payload)
}

func (s *EchoSuite) TestInspect() {
	payload := common.Hex2Bytes("deadbeef")
	result := s.tester.Inspect(payload)
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	s.Equal(payload, result.Reports[0].Payload)
}
