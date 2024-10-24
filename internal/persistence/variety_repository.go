package persistence

import (
	"github.com/duolacloud/crud-core/repositories"
	models_variety "github.com/yzimhao/trading_engine/v2/internal/models/variety"
)

type VarietyRepository interface {
	repositories.CrudRepository[models_variety.Variety, models_variety.CreateVariety, models_variety.UpdateVariety]
}
