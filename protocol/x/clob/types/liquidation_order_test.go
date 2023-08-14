package types_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/stretchr/testify/require"
)

var testLiquidationOrder = types.NewLiquidationOrder(
	constants.Alice_Num0,
	constants.ClobPair_Btc,
	false,
	100,
	200,
)

func TestLiquidationOrder_GetBaseQuantums(t *testing.T) {
	quantums := testLiquidationOrder.GetBaseQuantums()
	require.Equal(t, satypes.BaseQuantums(100), quantums)
}

func TestLiquidationOrder_GetOrderSubticks(t *testing.T) {
	subticks := testLiquidationOrder.GetOrderSubticks()
	require.Equal(t, types.Subticks(200), subticks)
}

func TestLiquidationOrder_IsBuy(t *testing.T) {
	tests := map[string]struct {
		side     bool
		expected bool
	}{
		"Is buy": {
			side:     true,
			expected: true,
		},
		"Is sell": {
			side:     false,
			expected: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			order := types.NewLiquidationOrder(
				constants.Alice_Num0,
				constants.ClobPair_Btc,
				tc.side,
				100,
				200,
			)

			result := order.IsBuy()
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestLiquidationOrder_GetOrderHash(t *testing.T) {
	tests := map[string]struct {
		order        types.LiquidationOrder
		expectedHash types.OrderHash
	}{
		"Can take SHA256 hash of an empty liquidation order": {
			order:        types.LiquidationOrder{},
			expectedHash: constants.LiquidationOrderHash_Empty,
		},
		"Can take SHA256 hash of a regular liquidation order": {
			order: *types.NewLiquidationOrder(
				constants.Alice_Num0,
				constants.ClobPair_Btc,
				true,
				10,
				10,
			),
			expectedHash: constants.LiquidationOrderHash_Alice_Number0_Perpetual0,
		},
		"Changing the subaccount ID changes the hash of the liquidation order": {
			order: *types.NewLiquidationOrder(
				constants.Alice_Num1,
				constants.ClobPair_Btc,
				true,
				10,
				10,
			),
			expectedHash: constants.LiquidationOrderHash_Alice_Number1_Perpetual0,
		},
		"Changing the perpetual ID changes the hash of the liquidation order": {
			order: *types.NewLiquidationOrder(
				constants.Alice_Num1,
				constants.ClobPair_Eth,
				true,
				10,
				10,
			),
			expectedHash: constants.LiquidationOrderHash_Alice_Number0_Perpetual1,
		},
		"Changing the CLOB pair ID doesn't change the hash of the liquidation order": {
			order: *types.NewLiquidationOrder(
				constants.Alice_Num0,
				constants.ClobPair_Btc2,
				true,
				10,
				10,
			),
			expectedHash: constants.LiquidationOrderHash_Alice_Number0_Perpetual0,
		},
		"Changing the order side doesn't change the hash of the liquidation order": {
			order: *types.NewLiquidationOrder(
				constants.Alice_Num0,
				constants.ClobPair_Btc2,
				false,
				10,
				10,
			),
			expectedHash: constants.LiquidationOrderHash_Alice_Number0_Perpetual0,
		},
		"Changing the quantums doesn't change the hash of the liquidation order": {
			order: *types.NewLiquidationOrder(
				constants.Alice_Num0,
				constants.ClobPair_Btc2,
				false,
				77,
				10,
			),
			expectedHash: constants.LiquidationOrderHash_Alice_Number0_Perpetual0,
		},
		"Changing the subticks doesn't change the hash of the liquidation order": {
			order: *types.NewLiquidationOrder(
				constants.Alice_Num0,
				constants.ClobPair_Btc2,
				false,
				10,
				88,
			),
			expectedHash: constants.LiquidationOrderHash_Alice_Number0_Perpetual0,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expectedHash, tc.order.GetOrderHash())
		})
	}
}

func TestLiquidationOrder_GetSubaccountId(t *testing.T) {
	subaccountId := testLiquidationOrder.GetSubaccountId()
	require.Equal(t, constants.Alice_Num0, subaccountId)
}

func TestLiquidationOrder_GetClobPairId(t *testing.T) {
	require.Equal(t, types.ClobPairId(0), testLiquidationOrder.GetClobPairId())
}

func TestLiquidationOrder_IsLiquidation(t *testing.T) {
	order := types.LiquidationOrder{}

	isLiquidation := order.IsLiquidation()
	require.True(t, isLiquidation)
}

func TestLiquidationOrder_MustGetLiquidatedPerpetualId(t *testing.T) {
	perpetualId := testLiquidationOrder.MustGetLiquidatedPerpetualId()
	expectedPerpetualId := constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId
	require.Equal(t, expectedPerpetualId, perpetualId)
}

func TestLiquidationOrder_MustGetOrderPanics(t *testing.T) {
	order := types.LiquidationOrder{}

	require.PanicsWithValue(
		t,
		"MustGetOrder: No underlying order on a LiquidationOrder type.",
		func() {
			order.MustGetOrder()
		},
	)
}

func TestNewLiquidationOrder_PanicsOnNonPerpetualClob(t *testing.T) {
	require.PanicsWithValue(
		t,
		"NewLiquidationOrder: Attempting to create liquidation order with a non-perpetual CLOB pair",
		func() {
			types.NewLiquidationOrder(
				constants.Alice_Num0,
				constants.ClobPair_Asset,
				false,
				1,
				1,
			)
		},
	)
}

func TestLiquidationOrder_IsReduceOnly(t *testing.T) {
	order := types.LiquidationOrder{}
	require.False(t, order.IsReduceOnly())
}

func TestSortedLiquidationOrder(t *testing.T) {
	tests := map[string]struct {
		orders []types.LiquidationOrder

		expectedOrdering []types.LiquidationOrder
	}{
		"sorts by clob pair id": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob1_Buy1ETH_Price3000,
				constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
			},
			expectedOrdering: []types.LiquidationOrder{
				constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
				constants.LiquidationOrder_Carl_Num0_Clob1_Buy1ETH_Price3000,
			},
		},
		"sorts by order side": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
				constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
			},
			expectedOrdering: []types.LiquidationOrder{
				constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
				constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
			},
		},
		"sorts buys in descending order": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC,
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy3_Price50_BTC,
			},
			expectedOrdering: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy3_Price50_BTC,
				constants.LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC,
			},
		},
		"sorts sells ascending order": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
				constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
			},
			expectedOrdering: []types.LiquidationOrder{
				constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
			},
		},
		"sorts by order size": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num1_Clob0_Sell01BTC_Price50000,
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
			},
			expectedOrdering: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
				constants.LiquidationOrder_Dave_Num1_Clob0_Sell01BTC_Price50000,
			},
		},
		"tie breaker ordering": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob1_Buy1ETH_Price3000,
				constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
				constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy3_Price50_BTC,
				constants.LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC,
				constants.LiquidationOrder_Dave_Num1_Clob0_Sell01BTC_Price50000,
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
			},
			expectedOrdering: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy3_Price50_BTC,
				constants.LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC,
				constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
				constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
				constants.LiquidationOrder_Dave_Num1_Clob0_Sell01BTC_Price50000,
				constants.LiquidationOrder_Carl_Num0_Clob1_Buy1ETH_Price3000,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sort.Sort(types.SortedLiquidationOrders(tc.orders))
			require.Equal(
				t,
				tc.expectedOrdering,
				tc.orders,
			)
		})
	}
}
