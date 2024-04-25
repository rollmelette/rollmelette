// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package echoapp

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

var payload = common.Hex2Bytes("deadbeef")
var msgSender = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")

func TestEchoSuite(t *testing.T) {
	suite.Run(t, new(EchoSuite))
}

type EchoSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *EchoSuite) SetupTest() {
	app := new(EchoApplication)
	s.tester = rollmelette.NewTester(app)
}

func (s *EchoSuite) TestAdvance() {
	result := s.tester.Advance(msgSender, payload)
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
	result := s.tester.Inspect(payload)
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	s.Equal(payload, result.Reports[0].Payload)
}
