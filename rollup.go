// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import "github.com/ethereum/go-ethereum/common"

// Metadata of the rollup advance input.
type Metadata struct {

	// InputIndex is the advance input index.
	InputIndex int

	// Sender is the account or contract that added the input to the input box.
	MsgSender common.Address

	// BlockNumber is the number of the block when the input was added to the L1 chain.
	BlockNumber int64

	// BlockNumber is the timestamp of the block when the input was added to the L1 chain.
	BlockTimestamp int64
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
	sendVoucher(destination common.Address, payload []byte) (int, error)

	// sendNotice sends a notice to the Rollup API and returns its index.
	sendNotice(payload []byte) (int, error)

	// sendNotice sends a report to the Rollup API.
	sendReport(payload []byte) error
}

// rollupRun is the interface of the Rollup API used by the run function.
// type rollupRun interface {
//
// 	// finishAndGetNext sends a finish request to the Rollup API.
// 	// If there is no error, it returns an advanceInput or an inspectInput.
// 	finishAndGetNext(status finishStatus) (any, error)
// }
