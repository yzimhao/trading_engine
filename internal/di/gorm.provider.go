package di

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGorm(v *viper.Viper) *gorm.DB {
	v.SetDefault("database_host", "localhost")
	v.SetDefault("database_port", 5432)
	v.SetDefault("database_user", "postgres")
	v.SetDefault("database_password", "postgres")
	v.SetDefault("database_name", "wallet_transfer")
	v.SetDefault("database_timezone", "Asia/Shanghai")
	v.SetDefault("database_debug", false)

	user := v.GetString("database_user")
	password := v.GetString("database_password")
	host := v.GetString("database_host")
	port := v.GetInt("database_port")
	dbname := v.GetString("database_name")
	timezone := v.GetString("database_timezone")
	debug := v.GetBool("database_debug")

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		host, port, user, password, dbname, timezone)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	if debug {
		db = db.Debug()
	}

	return db
}
