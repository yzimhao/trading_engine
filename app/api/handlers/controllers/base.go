package controllers

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/app/common"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type BaseController struct {
	logger      *zap.Logger
	productRepo persistence.ProductRepository
}

type inBaseContext struct {
	fx.In
	Logger      *zap.Logger
	ProductRepo persistence.ProductRepository
}

func NewBaseController(in inBaseContext) *BaseController {
	return &BaseController{
		logger:      in.Logger,
		productRepo: in.ProductRepo,
	}
}

func (ctrl *BaseController) Time(c *gin.Context) {
	common.ResponseOK(c, gin.H{
		"time": time.Now().Nanosecond(),
	})
}

// @Summary exchange info
// @Description get exchange info
// @ID v1.base.exchange_info
// @Tags base
// @Accept json
// @Produce json
// @Query param string true "symbol"
// @Success 200 {string} any
// @Router /api/v1/base/exchange_info [get]
func (ctrl *BaseController) ExchangeInfo(c *gin.Context) {
	symbol := strings.ToUpper(c.Query("symbol"))

	product, err := ctrl.productRepo.Get(symbol)
	if err != nil {
		common.ResponseError(c, err)
		return
	}

	common.ResponseOK(c, product)
}
