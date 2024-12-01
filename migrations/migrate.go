package migrations

import (
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"

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
	tables := []any{
		&entities.Asset{},
		&entities.AssetFreeze{},
		&entities.AssetLog{},
		&entities.UnfinishedOrder{},
		&entities.TradeVariety{},
		&entities.Variety{},
	}

	for _, table := range tables {
		indexes, _ := db.Migrator().GetIndexes(table)
		for _, index := range indexes {
			db.Migrator().DropIndex(table, index.Name())
		}
		db.Migrator().DropTable(table)
	}

	return nil
}
