package symbols

import "github.com/yzimhao/trading_engine/utils/app"

func NewVarieties(symbol string) *Varieties {
	db := app.Database().NewSession()
	defer db.Close()

	var row Varieties
	db.Where("symbol=?", symbol).Get(&row)
	return &row
}

func NewTradingVarieties(symbol string) *TradingVarieties {
	db := app.Database().NewSession()
	defer db.Close()

	var row TradingVarieties

	db.Where("symbol=?", symbol).Get(&row)
	if row.Id > 0 {
		row.Target = *newVarietiesById(row.TargetSymbolId)
		row.Base = *newVarietiesById(row.BaseSymbolId)
	}
	return &row
}

func newVarietiesById(id int) *Varieties {
	db := app.Database().NewSession()
	defer db.Close()

	var row Varieties
	db.Where("id=?", id).Get(&row)
	return &row
}
