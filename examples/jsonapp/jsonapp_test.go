// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package jsonapp

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/rollmelette"
	"github.com/stretchr/testify/suite"
)

var gm = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
var player = common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")

func TestJsonSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}

type GameSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *GameSuite) SetupTest() {
	app := NewGameApplication(gm)
	s.tester = rollmelette.NewTester(app)
}

func (s *GameSuite) TestItAddsMonster() {
	input := `{"kind":"AddMonster","payload":{"name":"dragon","hitPoints":10}}`
	result := s.tester.Advance(gm, []byte(input))
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	expectedState := `{"monsters":{"dragon":{"name":"dragon","hitPoints":10}}}`
	s.Equal(expectedState, string(result.Reports[0].Payload))
}

func (s *GameSuite) TestItAttacksMonster() {
	// add monster
	input := `{"kind":"AddMonster","payload":{"name":"goblin","hitPoints":3}}`
	result := s.tester.Advance(gm, []byte(input))
	s.Nil(result.Err)

	// attack monster
	input = `{"kind":"AttackMonster","payload":{"monsterName":"goblin","damage":2}}`
	result = s.tester.Advance(player, []byte(input))
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	expectedState := `{"monsters":{"goblin":{"name":"goblin","hitPoints":1}}}`
	s.Equal(expectedState, string(result.Reports[0].Payload))

	// kill monster
	result = s.tester.Advance(player, []byte(input))
	s.Nil(result.Err)
	s.Len(result.Reports, 1)
	expectedState = `{"monsters":{}}`
	s.Equal(expectedState, string(result.Reports[0].Payload))
}

func (s *GameSuite) TestItFailsToParseJson() {
	input := `not a json`
	result := s.tester.Advance(player, []byte(input))
	s.ErrorContains(result.Err, "failed to unmarshal input")
}

func (s *GameSuite) TestItRejectsInvalidKind() {
	input := `{"kind":"invalid"}`
	result := s.tester.Advance(player, []byte(input))
	s.ErrorContains(result.Err, "invalid input kind: invalid")
}

func (s *GameSuite) TestItRejectsAddMonsterIfNotGM() {
	input := `{"kind":"AddMonster","payload":{"name":"goblin","hitPoints":2}}`
	result := s.tester.Advance(player, []byte(input))
	s.ErrorContains(result.Err, "only GM can add monsters")
}

func (s *GameSuite) TestItRejectsAddMonsterIfHitPointsIsNegative() {
	input := `{"kind":"AddMonster","payload":{"name":"goblin","hitPoints":-2}}`
	result := s.tester.Advance(gm, []byte(input))
	s.ErrorContains(result.Err, "hit points must be positive")
}

func (s *GameSuite) TestItRejectsAddMonsterIfMonsterAlreadyExists() {
	input := `{"kind":"AddMonster","payload":{"name":"goblin","hitPoints":3}}`
	result := s.tester.Advance(gm, []byte(input))
	s.Nil(result.Err)
	result = s.tester.Advance(gm, []byte(input))
	s.ErrorContains(result.Err, "monster with this name already exists")
}

func (s *GameSuite) TestItRejectsAttackIfDamageIsNegative() {
	input := `{"kind":"AttackMonster","payload":{"monsterName":"goblin","damage":-2}}`
	result := s.tester.Advance(gm, []byte(input))
	s.ErrorContains(result.Err, "negative damage")
}

func (s *GameSuite) TestItRejectsAttackIfMonsterDoesntExist() {
	input := `{"kind":"AttackMonster","payload":{"monsterName":"goblin","damage":2}}`
	result := s.tester.Advance(gm, []byte(input))
	s.ErrorContains(result.Err, "monster not found")
}
