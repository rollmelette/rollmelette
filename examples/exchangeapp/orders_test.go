// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package exchangeapp

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestOrderBookSuite(t *testing.T) {
	suite.Run(t, new(OrderBookSuite))
}

type OrderBookSuite struct {
	suite.Suite
	book *OrderBook
}

func (s *OrderBookSuite) SetupTest() {
	s.book = NewOrderBook()
}

func (s *OrderBookSuite) TestItAddsBuyOrders() {
	reports, err := s.book.AddOrder(&Order{
		OrderData: OrderData{
			Kind:     OrderBuy,
			Price:    big.NewInt(100),
			Quantity: big.NewInt(10),
		},
		ID: 1,
	})
	s.Nil(err)
	s.checkReports(reports,
		`{"kind":"OrderAdded","data":{"kind":"Buy","price":100,"quantity":10,"id":1}}`,
	)

	reports, err = s.book.AddOrder(&Order{
		OrderData: OrderData{
			Kind:     OrderBuy,
			Price:    big.NewInt(200),
			Quantity: big.NewInt(20),
		},
		ID: 2,
	})
	s.Nil(err)
	s.checkReports(reports,
		`{"kind":"OrderAdded","data":{"kind":"Buy","price":200,"quantity":20,"id":2}}`,
	)

	// TODO check internal structures
}

func (s *OrderBookSuite) TestItFailsToAddRepeatedOrder() {
	order := &Order{
		OrderData: OrderData{
			Kind:     OrderBuy,
			Price:    big.NewInt(100),
			Quantity: big.NewInt(10),
		},
		ID: 1,
	}

	reports, err := s.book.AddOrder(order)
	s.Nil(err)
	s.NotNil(reports)

	_, err = s.book.AddOrder(order)
	s.ErrorContains(err, "order already in book")
}

func (s *OrderBookSuite) TestSomething() {
	order := &Order{
		OrderData: OrderData{
			Kind:     OrderBuy,
			Price:    big.NewInt(100),
			Quantity: big.NewInt(10),
		},
		ID: 1,
	}
	reports, err := s.book.AddOrder(order)
	s.Nil(err)
	s.NotNil(reports)

	order = &Order{
		OrderData: OrderData{
			Kind:     OrderBuy,
			Price:    big.NewInt(200),
			Quantity: big.NewInt(20),
		},
		ID: 2,
	}
	reports, err = s.book.AddOrder(order)
	s.Nil(err)
	s.NotNil(reports)

	order = &Order{
		OrderData: OrderData{
			Kind:     OrderSell,
			Price:    big.NewInt(150),
			Quantity: big.NewInt(40),
		},
		ID: 3,
	}
	reports, err = s.book.AddOrder(order)
	s.Nil(err)
	s.NotNil(reports)

	reportsPretty, err := json.MarshalIndent(reports, "", "  ")
	s.Require().Nil(err)
	s.T().Log(string(reportsPretty))
}

func (s *OrderBookSuite) checkReports(reports []Report, expectedReports ...string) {
	expected := "[" + strings.Join(expectedReports, ",") + "]"
	reportsJson, err := json.Marshal(reports)
	s.Require().Nil(err)
	s.Equal(expected, string(reportsJson))
}
