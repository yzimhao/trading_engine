package symbols

func NewVarieties(symbol string) *Varieties {
	var row Varieties
	db.Where("symbol=?", symbol).Get(&row)
	return &row
}

func NewTradingVarieties(symbol string) *TradingVarieties {
	var row TradingVarieties

	db.Where("symbol=?", symbol).Get(&row)
	if row.Id > 0 {
		row.Target = *newVarietiesById(row.TargetSymbolId)
		row.Standard = *newVarietiesById(row.StandardSymbolId)
	}
	return &row
}

func newVarietiesById(id int) *Varieties {
	var row Varieties
	db.Where("id=?", id).Get(&row)
	return &row
}
