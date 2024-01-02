// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// VoucherMock represents a voucher received by the mock.
type VoucherMock struct {
	Destination common.Address
	Payload     []byte
}

// NoticeMock represents a notice received by the mock.
type NoticeMock struct {
	Payload []byte
}

// ReportMock represents a report received by the mock.
type ReportMock struct {
	Payload []byte
}

// rollupHttp implements the Rollup API by calling the Rollup HTTP server.
type RollupMock struct {
	Vouchers []VoucherMock
	Notices  []NoticeMock
	Reports  []ReportMock
}

// rollup interface ////////////////////////////////////////////////////////////////////////////////

func (m *RollupMock) sendVoucher(
	ctx context.Context,
	destination common.Address,
	payload []byte,
) (int, error) {
	m.Vouchers = append(m.Vouchers, VoucherMock{
		Destination: destination,
		Payload:     payload,
	})
	return len(m.Vouchers), nil
}

func (m *RollupMock) sendNotice(ctx context.Context, payload []byte) (int, error) {
	m.Notices = append(m.Notices, NoticeMock{
		Payload: payload,
	})
	return len(m.Notices), nil
}

func (m *RollupMock) sendReport(ctx context.Context, payload []byte) error {
	m.Reports = append(m.Reports, ReportMock{
		Payload: payload,
	})
	return nil
}

// mock functions /////////////////////////////////////////////////////////////////////////////////

func (m *RollupMock) reset() {
	m.Vouchers = nil
	m.Notices = nil
	m.Reports = nil
}
