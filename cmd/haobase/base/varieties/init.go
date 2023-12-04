package varieties

import (
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils/app"
)

func Init() {
	CreateTable()
}

func CreateTable() {
	db := app.Database()
	err := db.Sync2(
		new(Varieties),
		new(TradingVarieties),
	)

	if err != nil {
		app.Logger.Error(err)
	}
}

func DemoData() {
	symbols := []Varieties{
		Varieties{
			Symbol:        "usd",
			Name:          "美元",
			ShowPrecision: 2,
			MinPrecision:  8,
			Status:        types.StatusEnabled,
		},
		Varieties{
			Symbol:        "eur",
			Name:          "欧元",
			ShowPrecision: 2,
			MinPrecision:  8,
			Status:        types.StatusEnabled,
		},
		Varieties{
			Symbol:        "jpy",
			Name:          "日元",
			ShowPrecision: 2,
			MinPrecision:  8,
			Status:        types.StatusEnabled,
		},
	}

	db := app.Database().NewSession()
	defer db.Close()

	if empty, _ := db.IsTableEmpty(new(Varieties)); empty {
		_, err := db.Insert(symbols)
		if err != nil {
			app.Logger.Error(err)
		}
	}

	usd := NewVarieties("usd")
	jpy := NewVarieties("jpy")
	eur := NewVarieties("eur")
	pairs := []TradingVarieties{
		TradingVarieties{
			Symbol:         "usdjpy",
			Name:           "美日",
			TargetSymbolId: usd.Id,
			BaseSymbolId:   jpy.Id,
			PricePrecision: 3,
			QtyPrecision:   2,
			Status:         types.StatusEnabled,
			AllowMinQty:    "0.01",
			AllowMaxQty:    "0",
			AllowMinAmount: "1",
			AllowMaxAmount: "0",
			FeeRate:        "0.005",
		},
		TradingVarieties{
			Symbol:         "eurusd",
			Name:           "欧美",
			TargetSymbolId: eur.Id,
			BaseSymbolId:   usd.Id,
			PricePrecision: 5,
			QtyPrecision:   2,
			Status:         types.StatusEnabled,
			AllowMinQty:    "0.01",
			AllowMaxQty:    "0",
			AllowMinAmount: "1",
			AllowMaxAmount: "0",
			FeeRate:        "0.001",
		},
	}
	if empty, _ := db.IsTableEmpty(new(TradingVarieties)); empty {
		_, err := db.Insert(pairs)
		if err != nil {
			app.Logger.Error(err)
		}
	}
}
