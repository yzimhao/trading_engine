package migrations

import (
	"fmt"
	"strings"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const dialect = "postgres"

var migrations = &migrate.EmbedFileSystemMigrationSource{
	FileSystem: Migrations,
	Root:       dialect,
}

func MigrateUp(db *gorm.DB, cfg *viper.Viper, logger *zap.Logger) error {

	rawDb, err := db.DB()
	if err != nil {
		return err
	}

	applied, err := migrate.Exec(rawDb, dialect, migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("failed to apply migrations. %v", err)
	}

	logger.Info("migrations applied", zap.Int("applied", applied))
	return nil

}

func MigrateDown(db *gorm.DB, cfg *viper.Viper, logger *zap.Logger) error {
	rawDb, err := db.DB()
	if err != nil {
		return err
	}

	applied, err := migrate.Exec(rawDb, dialect, migrations, migrate.Down)
	if err != nil {
		return fmt.Errorf("failed to apply migrations. %v", err)
	}
	logger.Info("migrations applied", zap.Int("applied", applied))
	return nil
}

func MigrateClean(db *gorm.DB, cfg *viper.Viper, logger *zap.Logger) error {
	//TODO  justfor development
	needDroptables := []any{
		&entities.Asset{},
		&entities.Product{},
		&entities.UserAssetFreeze{},
		&entities.UserAssetLog{},
		&entities.UnfinishedOrder{},
		&entities.UserAsset{},
		"order_",
		"trade_record_",
		"kline_",
	}

	allTables, err := db.Migrator().GetTables()
	if err != nil {
		logger.Debug("get tables failed", zap.Error(err))
		return err
	}
	logger.Sugar().Debugf("all tables %v", allTables)

	var dropTables []any
	for _, tt := range allTables {
		for _, needDrop := range needDroptables {
			switch t := needDrop.(type) {
			case string:
				if strings.HasPrefix(tt, t) {
					dropTables = append(dropTables, tt)
				}
			default:
				dropTables = append(dropTables, needDrop)
			}
		}
	}

	logger.Sugar().Debugf("dropTables: %#v", dropTables)
	for _, dropTable := range dropTables {
		indexes, err := db.Migrator().GetIndexes(dropTable)
		if err != nil {
			logger.Debug("get indexes failed", zap.Error(err))
			continue
		}
		for _, index := range indexes {
			db.Migrator().DropIndex(dropTable, index.Name())
		}
		db.Migrator().DropTable(dropTable)
	}

	return nil
}
