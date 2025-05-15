package provider

import (
	"context"
	"crypto/tls"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewRedis(v *viper.Viper, logger *zap.Logger) *redis.Client {
	v.SetDefault("redis.address", "localhost:6379")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.username", "")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.tls_enabled", false)

	address := v.GetString("redis.address")
	db := v.GetInt("redis.db")
	username := v.GetString("redis.username")
	password := v.GetString("redis.password")
	tlsEnabled := v.GetBool("redis.tls_enabled")

	opts := redis.Options{
		Addr:     address,
		DB:       db,
		Username: username,
		Password: password,
	}
	if tlsEnabled {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := redis.NewClient(&opts)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.Sugar().Fatalf("failed to connect to redis: %v, addr: %s db: %d", err, address, db)
	}
	logger.Info("connected to redis", zap.String("addr", address), zap.Int("db", db))
	return client

}
