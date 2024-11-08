// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package integration

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/rollmelette/rollmelette/examples/echoapp"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"
)

const TestTimeout = 15 * time.Second

func TestNonodoSuite(t *testing.T) {
	suite.Run(t, new(NonodoSuite))
}

type NonodoSuite struct {
	suite.Suite
	group  *errgroup.Group
	ctx    context.Context
	cancel context.CancelFunc
	nonodo *exec.Cmd
}

// Setup ///////////////////////////////////////////////////////////////////////////////////////////

func (s *NonodoSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), TestTimeout)
	s.group, s.ctx = errgroup.WithContext(s.ctx)

	// start nonodo
	nonodo := exec.CommandContext(s.ctx, "brunodo", "-d")
	out := NewNotifyWriter(os.Stdout, "nonodo: ready")
	nonodo.Stdout = out
	s.group.Go(nonodo.Run)
	s.nonodo = nonodo
	select {
	case <-out.ready:
	case <-s.ctx.Done():
		slog.Debug("Error", "err", s.ctx.Err())
		s.T().Error(s.ctx.Err())
	}

	// start test app
	s.group.Go(func() error {
		opts := rollmelette.NewRunOpts()
		app := new(echoapp.EchoApplication)
		return rollmelette.Run(s.ctx, opts, app)
	})
}

func (s *NonodoSuite) TearDownTest() {
	s.cancel()
	_ = exec.Command("pkill", "brunodo").Run()
	_ = exec.Command("pkill", "nonodo").Run()
}

// Test Cases //////////////////////////////////////////////////////////////////////////////////////

func (s *NonodoSuite) XTestAdvance() {
	payload := common.Hex2Bytes("deadbeef")
	err := Advance(s.ctx, "http://127.0.0.1:8545", payload)
	s.Require().Nil(err)
	client := graphql.NewClient("http://127.0.0.1:8080/graphql", nil)
	err = waitForInput(s.ctx, client, 0)
	s.Require().Nil(err)
	result, err := getNodeState(s.ctx, client)
	s.Require().Nil(err)
	s.Require().Len(result.Inputs.Edges, 1)
	input := result.Inputs.Edges[0].Node

	// nolint
	expectedVoucherPayload := "0x237a816f000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000004deadbeef00000000000000000000000000000000000000000000000000000000"
	expectedNoticePayload := "0xc258d6e500000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004deadbeef00000000000000000000000000000000000000000000000000000000"
	s.Require().Equal(payload, common.Hex2Bytes(input.Payload[2:]))
	s.Require().Equal(expectedVoucherPayload, input.Vouchers.Edges[0].Node.Payload)
	s.Require().Equal(expectedNoticePayload, input.Notices.Edges[0].Node.Payload)
	s.Require().Equal(payload, common.Hex2Bytes(input.Reports.Edges[0].Node.Payload[2:]))
}

func (s *NonodoSuite) TestInspect() {
	payload := common.Hex2Bytes("deadbeef")
	url := fmt.Sprintf("http://127.0.0.1:8080/inspect/%s", common.Address{})
	response, err := Inspect(s.ctx, url, payload)
	s.Require().Nil(err)
	s.Require().Len(response.Reports, 1)
	s.Require().Equal(payload, common.Hex2Bytes(response.Reports[0].Payload[2:]))
}

// Helper functions ////////////////////////////////////////////////////////////////////////////////

func waitForInput(ctx context.Context, client graphql.Client, _index int) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		result, err := getInputStatus(ctx, client, "0")
		if err != nil && !strings.Contains(err.Error(), "input not found") {
			return fmt.Errorf("failed to get input status: %w", err)
		}
		if result.Input.Status == CompletionStatusAccepted {
			return nil
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
