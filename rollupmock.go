// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// TestVoucher represents a voucher received by the mock.
type TestVoucher struct {
	Destination common.Address
	Value       *big.Int
	Payload     []byte
}

// TestDelegateCallVoucher represents a delegate call voucher received by the mock.
type TestDelegateCallVoucher struct {
	Destination common.Address
	Payload     []byte
}

// TestNotice represents a notice received by the mock.
type TestNotice struct {
	Payload []byte
}

// TestReport represents a report received by the mock.
type TestReport struct {
	Payload []byte
}

// rollupHttp implements the Rollup API by calling the Rollup HTTP server.
type rollupMock struct {
	Vouchers             []TestVoucher
	DelegateCallVouchers []TestDelegateCallVoucher
	Notices              []TestNotice
	Reports              []TestReport
}

// rollup interface ////////////////////////////////////////////////////////////////////////////////

func (m *rollupMock) sendVoucher(
	ctx context.Context,
	destination common.Address,
	value *big.Int,
	payload []byte,
) (int, error) {
	m.Vouchers = append(m.Vouchers, TestVoucher{
		Destination: destination,
		Value:       value,
		Payload:     payload,
	})
	return len(m.Vouchers), nil
}

func (m *rollupMock) sendDelegateCallVoucher(ctx context.Context, destination common.Address, payload []byte) (int, error) {
	m.DelegateCallVouchers = append(m.DelegateCallVouchers, TestDelegateCallVoucher{
		Destination: destination,
		Payload:     payload,
	})
	return len(m.DelegateCallVouchers), nil
}

func (m *rollupMock) sendNotice(ctx context.Context, payload []byte) (int, error) {
	m.Notices = append(m.Notices, TestNotice{
		Payload: payload,
	})
	return len(m.Notices), nil
}

func (m *rollupMock) sendReport(ctx context.Context, payload []byte) error {
	m.Reports = append(m.Reports, TestReport{
		Payload: payload,
	})
	return nil
}

// mock functions /////////////////////////////////////////////////////////////////////////////////

func (m *rollupMock) reset() {
	m.Vouchers = nil
	m.Notices = nil
	m.Reports = nil
}
