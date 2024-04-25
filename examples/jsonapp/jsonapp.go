// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// Jsonapp contains an example of a simple game application that uses JSON as inputs and outputs.
package jsonapp

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
)

// GameState is represents the state of the game.
type GameState struct {
	Monsters map[string]Monster `json:"monsters"`
}

// Monster represents a monster in the game.
type Monster struct {
	Name      string `json:"name"`
	HitPoints int    `json:"hitPoints"`
}

// InputKind is an enum that represents the kind of input.
type InputKind string

const (
	// AddMonster adds a monster to the game. It should be called by the GM.
	AddMonster InputKind = "AddMonster"

	// AttackMonster attacks the given monster.
	AttackMonster InputKind = "AttackMonster"
)

// Input has a kind and a payload with specific data.
type Input struct {
	Kind    InputKind       `json:"kind"`
	Payload json.RawMessage `json:"payload"`
}

// AddMonsterPayload is the payload for the AddMonster input.
type AddMonsterPayload = Monster

// AttackMonsterPayload is the payload for the AttackMonster input.
type AttackMonsterPayload struct {
	MonsterName string `json:"monsterName"`
	Damage      int    `json:"damage"`
}

// GameApplication is an application for a simple RPG game.
// The GM can add monsters and the players can attack them.
// Anyone can attack the the monsters.
// For each input, the application emits a report with the current game state.
type GameApplication struct {
	gm    common.Address
	state GameState
}

func NewGameApplication(gm common.Address) *GameApplication {
	return &GameApplication{
		gm: gm,
		state: GameState{
			Monsters: make(map[string]Monster),
		},
	}
}

func (a *GameApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	var input Input
	err := json.Unmarshal(payload, &input)
	if err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}
	switch input.Kind {
	case AddMonster:
		var inputPayload AddMonsterPayload
		err = json.Unmarshal(input.Payload, &inputPayload)
		if err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		err = a.handleAddMonster(metadata, inputPayload)
		if err != nil {
			return err
		}
	case AttackMonster:
		var inputPayload AttackMonsterPayload
		err = json.Unmarshal(input.Payload, &inputPayload)
		if err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		err = a.handleAttackMonster(inputPayload)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid input kind: %v", input.Kind)
	}
	return a.Inspect(env, nil)
}

func (a *GameApplication) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	bytes, err := json.Marshal(a.state)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	env.Report(bytes)
	return nil
}

func (a *GameApplication) handleAddMonster(
	metadata rollmelette.Metadata,
	inputPayload AddMonsterPayload,
) error {
	if metadata.MsgSender != a.gm {
		return fmt.Errorf("only GM can add monsters")
	}
	if inputPayload.HitPoints <= 0 {
		return fmt.Errorf("hit points must be positive")
	}
	_, ok := a.state.Monsters[inputPayload.Name]
	if ok {
		return fmt.Errorf("monster with this name already exists")
	}
	a.state.Monsters[inputPayload.Name] = inputPayload
	return nil
}

func (a *GameApplication) handleAttackMonster(inputPayload AttackMonsterPayload) error {
	if inputPayload.Damage < 0 {
		return fmt.Errorf("negative damage")
	}
	monster, ok := a.state.Monsters[inputPayload.MonsterName]
	if !ok {
		return fmt.Errorf("monster not found")
	}
	monster.HitPoints -= inputPayload.Damage
	if monster.HitPoints <= 0 {
		// killed the monster
		delete(a.state.Monsters, inputPayload.MonsterName)
	} else {
		// update the monster in the map
		a.state.Monsters[inputPayload.MonsterName] = monster
	}
	return nil
}
