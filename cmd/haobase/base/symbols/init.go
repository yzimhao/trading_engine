package symbols

import (
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/trading_engine/utils/app"
)

var ()

func Init() {
	init_db()
}

func init_db() {
	db := app.Database()

	err := db.Sync2(
		new(Varieties),
		new(TradingVarieties),
	)

	if err != nil {
		panic(err)
	}
}

func DemoData() {
	symbols := []Varieties{
		Varieties{
			Symbol:        "usd",
			Name:          "美元",
			ShowPrecision: 2,
			MinPrecision:  8,
			Status:        StatusEnabled,
		},
		Varieties{
			Symbol:        "eur",
			Name:          "欧元",
			ShowPrecision: 2,
			MinPrecision:  8,
			Status:        StatusEnabled,
		},
		Varieties{
			Symbol:        "jpy",
			Name:          "日元",
			ShowPrecision: 2,
			MinPrecision:  8,
			Status:        StatusEnabled,
		},
	}

	db := app.Database().NewSession()
	defer db.Close()

	_, err := db.Insert(symbols)
	if err != nil {
		logrus.Error(err)
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
			Status:         StatusEnabled,
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
			Status:         StatusEnabled,
			AllowMinQty:    "0.01",
			AllowMaxQty:    "0",
			AllowMinAmount: "1",
			AllowMaxAmount: "0",
			FeeRate:        "0.001",
		},
	}
	_, err = db.Insert(pairs)
	if err != nil {
		logrus.Error(err)
	}
}
