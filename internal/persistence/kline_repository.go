package persistence

import (
	"github.com/duolacloud/crud-core/repositories"
	models "github.com/yzimhao/trading_engine/v2/internal/models/kline"
)

type KlineRepository interface {
	repositories.CrudRepository[models.Kline, models.CreateKline, models.UpdateKline]
}
