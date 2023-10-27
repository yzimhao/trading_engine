package www

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
	"github.com/yzimhao/trading_engine/utils"
)

// 全部交易品类
type tvarieties struct {
	Symbol         string              `json:"symbol"`
	Name           string              `json:"name"`
	TargetSymbolId int                 `json:"target_symbol_id"`
	BaseSymbolId   int                 `json:"base_symbol_id"`
	PricePrecision int                 `json:"price_precision"`
	QtyPrecision   int                 `json:"qty_precision"`
	AllowMinQty    string              `json:"allow_min_qty"`
	AllowMaxQty    string              `json:"allow_max_qty"`
	AllowMinAmount string              `json:"allow_min_amount"`
	AllowMaxAmount string              `json:"allow_max_amount"`
	FeeRate        string              `json:"fee_rate"`
	UpdateTime     utils.Time          `json:"update_time"`
	Target         varieties.Varieties `json:"target"`
	Base           varieties.Varieties `json:"base"`
}

func trading_varieties(ctx *gin.Context) {
	data := make([]tvarieties, 0)

	for _, v := range base.NewTSymbols().All() {
		item := tvarieties{
			Symbol:         v.Symbol,
			Name:           v.Name,
			TargetSymbolId: v.TargetSymbolId,
			BaseSymbolId:   v.BaseSymbolId,
			PricePrecision: v.PricePrecision,
			QtyPrecision:   v.QtyPrecision,
			AllowMinQty:    v.AllowMinQty.String(),
			AllowMaxQty:    v.AllowMaxQty.String(),
			AllowMinAmount: v.AllowMinAmount.String(),
			AllowMaxAmount: v.AllowMaxAmount.String(),
			FeeRate:        v.FeeRate.String(),
			UpdateTime:     v.UpdateTime,
		}
		data = append(data, item)
	}
	utils.ResponseOkJson(ctx, data)
}

func varieties_config(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	v, _ := base.NewTSymbols().Get(symbol)

	item := tvarieties{
		Symbol:         v.Symbol,
		Name:           v.Name,
		TargetSymbolId: v.TargetSymbolId,
		BaseSymbolId:   v.BaseSymbolId,
		PricePrecision: v.PricePrecision,
		QtyPrecision:   v.QtyPrecision,
		AllowMinQty:    v.AllowMinQty.String(),
		AllowMaxQty:    v.AllowMaxQty.String(),
		AllowMinAmount: v.AllowMinAmount.String(),
		AllowMaxAmount: v.AllowMaxAmount.String(),
		FeeRate:        v.FeeRate.String(),
		UpdateTime:     v.UpdateTime,
	}
	utils.ResponseOkJson(ctx, item)
}
