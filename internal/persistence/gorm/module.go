package gorm

import (
	"github.com/duolacloud/crud-core/datasource"
	"go.uber.org/fx"
	_gorm "gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(
		datasource.NewDataSource[_gorm.DB],
		NewAssetsRepo,
	),
)
