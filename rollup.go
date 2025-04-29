// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Metadata of the rollup advance input.
type Metadata struct {
	// ChainId is the chain id of the base layer.
	ChainId int64

	// AppContract is the address of the application contract.
	AppContract common.Address

	// InputIndex is the advance input index.
	InputIndex int

	// Sender is the account or contract that added the input to the input box.
	MsgSender common.Address

	// BlockNumber is the number of the block when the input was added to the L1 chain.
	BlockNumber int64

	// BlockNumber is the timestamp of the block when the input was added to the L1 chain.
	BlockTimestamp int64

	// PrevRandao is the previous randao value of the block when the input was added to the L1 chain.
	PrevRandao string
}

// finishStatus is the status when finishing a rollup input.
type finishStatus string

const (
	finishStatusAccept finishStatus = "accept"
	finishStatusReject finishStatus = "reject"
)

// advanceInput represents an advance input from finish.
type advanceInput struct {
	Metadata Metadata
	Payload  []byte
}

// inspectInput represent an inspect input from finish.
type inspectInput struct {
	Payload []byte
}

// rollupEnv is the interface of the Rollup API used by the env struct.
type rollupEnv interface {

	// sendVoucher sends a voucher to the Rollup API and returns its index.
	sendVoucher(ctx context.Context, destination common.Address, value *big.Int, payload []byte) (int, error)

	// sendDelegateCallVoucher sends a delegate call voucher to the Rollup API and returns its index.
	sendDelegateCallVoucher(ctx context.Context, destination common.Address, payload []byte) (int, error)

	// sendNotice sends a notice to the Rollup API and returns its index.
	sendNotice(ctx context.Context, payload []byte) (int, error)

	// sendReport sends a report to the Rollup API.
	sendReport(ctx context.Context, payload []byte) error
}

// rollupRun is the interface of the Rollup API used by the run function.
// type rollupRun interface {
//
// 	// finishAndGetNext sends a finish request to the Rollup API.
// 	// If there is no error, it returns an advanceInput or an inspectInput.
// 	finishAndGetNext(status finishStatus) (any, error)
// }
