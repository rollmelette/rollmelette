// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package exchangeapp

import (
	"fmt"
	"math/big"
	"slices"
)

// OrderKind is the kind of the order, which can be buy or sell.
type OrderKind string

const (
	OrderBuy  OrderKind = "Buy"
	OrderSell OrderKind = "Sell"
)

// OrderData is the data of an order.
type OrderData struct {
	Kind     OrderKind `json:"kind"`
	Price    *big.Int  `json:"price"`
	Quantity *big.Int  `json:"quantity"`
}

// Order is the internal representation of an order.
type Order struct {
	OrderData
	ID int `json:"id"` // input index
}

// ReportKind is the kind of report.
type ReportKind string

const (
	ReportOrderAdded   ReportKind = "OrderAdded"
	ReportOrderRemoved ReportKind = "OrderRemoved"
	ReportTrade        ReportKind = "Trade"
)

// Report contains information about the order book transactions.
type Report struct {
	Kind ReportKind `json:"kind"`
	Data any        `json:"data"`
}

type ReportOrderAddedData = Order

func NewReportOrderAdded(order *Order) Report {
	return Report{
		Kind: ReportOrderAdded,
		Data: ReportOrderAddedData(*order),
	}
}

type ReportOrderRemovedData = Order

func NewReportOrderRemoved(order *Order) Report {
	return Report{
		Kind: ReportOrderRemoved,
		Data: ReportOrderRemovedData(*order),
	}
}

type ReportTradeData struct {
	BuyOrder  Order    `json:"buyOrder"`
	SellOrder Order    `json:"sellOrder"`
	Quantity  *big.Int `json:"quantity"`
}

func NewReportTrade(buyOrder *Order, sellOrder *Order, quantity *big.Int) Report {
	return Report{
		Kind: ReportTrade,
		Data: ReportTradeData{
			BuyOrder:  *buyOrder,
			SellOrder: *sellOrder,
			Quantity:  quantity,
		},
	}
}

// OrderBook stores the orders and provides methods to operate on them.
type OrderBook struct {
	// Orders contains all orders in the book.
	Orders map[int]*Order

	// BuyOrders contains the buy orders sorted by highest price.
	BuyOrders []*Order

	// SellOrders contains the sell orders sorted by lowest price.
	SellOrders []*Order
}

// NewOrderBook creates a new order book.
func NewOrderBook() *OrderBook {
	return &OrderBook{
		Orders: make(map[int]*Order),
	}
}

// AddOrder adds the order to the book and tries to fulfill it.
// It returns the list of reports.
func (b *OrderBook) AddOrder(order *Order) ([]Report, error) {
	if _, ok := b.Orders[order.ID]; ok {
		return nil, fmt.Errorf("order already in book")
	}
	switch order.Kind {
	case OrderBuy:
		b.BuyOrders = append(b.BuyOrders, order)
		slices.SortFunc(b.BuyOrders, func(a *Order, b *Order) int {
			return b.Price.Cmp(a.Price)
		})
	case OrderSell:
		b.SellOrders = append(b.SellOrders, order)
		slices.SortFunc(b.SellOrders, func(a *Order, b *Order) int {
			return a.Price.Cmp(b.Price)
		})
	default:
		return nil, fmt.Errorf("invalid order kind")
	}
	b.Orders[order.ID] = order
	reports := []Report{NewReportOrderAdded(order)}
	return b.match(reports), nil
}

// match tries to match orders in the order book.
func (b *OrderBook) match(reports []Report) []Report {
	i := 0
	j := 0

	// Match orders
	for i < len(b.BuyOrders) && j < len(b.SellOrders) {
		buy := b.BuyOrders[i]
		sell := b.SellOrders[j]

		if buy.Price.Cmp(sell.Price) < 0 {
			// Buy price is smaller than the sell price
			break
		}

		quantity := bigMin(buy.Quantity, sell.Quantity)
		reports = append(reports, NewReportTrade(buy, sell, quantity))

		buy.Quantity = new(big.Int).Sub(buy.Quantity, quantity)
		sell.Quantity = new(big.Int).Sub(sell.Quantity, quantity)

		if buy.Quantity.Cmp(big.NewInt(0)) == 0 {
			i++
		}
		if sell.Quantity.Cmp(big.NewInt(0)) == 0 {
			j++
		}
	}

	// Remove fulfilled orders
	if i > 0 {
		for ii := 0; ii < i; ii++ {
			order := b.BuyOrders[ii]
			delete(b.Orders, order.ID)
			reports = append(reports, NewReportOrderRemoved(order))
		}
		b.BuyOrders = slices.Delete(b.BuyOrders, 0, i)
	}
	if j > 0 {
		for jj := 0; jj < j; jj++ {
			order := b.SellOrders[jj]
			delete(b.Orders, order.ID)
			reports = append(reports, NewReportOrderRemoved(order))
		}
		b.SellOrders = slices.Delete(b.SellOrders, 0, j)
	}

	return reports
}

// RemoveOrder removes an order from the book.
func (b *OrderBook) RemoveOrder(id int) ([]Report, error) {
	order, ok := b.Orders[id]
	if !ok {
		return nil, fmt.Errorf("order not found")
	}
	switch order.Kind {
	case OrderBuy:
		b.BuyOrders = slices.DeleteFunc(b.BuyOrders, func(o *Order) bool {
			return o.ID == id
		})
	case OrderSell:
		b.SellOrders = slices.DeleteFunc(b.SellOrders, func(o *Order) bool {
			return o.ID == id
		})
	default:
		panic("impossible")
	}
	delete(b.Orders, id)
	reports := []Report{NewReportOrderRemoved(order)}
	return reports, nil
}

func bigMin(a *big.Int, b *big.Int) *big.Int {
	if a.Cmp(b) < 0 {
		return a
	}
	return b
}
