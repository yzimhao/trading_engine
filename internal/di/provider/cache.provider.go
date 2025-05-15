package provider

import (
	redis "github.com/duolacloud/crud-cache-redis"
	"github.com/duolacloud/crud-core/cache"
	goredis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewCache(v *viper.Viper, cli *goredis.Client) (cache.Cache, error) {
	return redis.New(
		redis.WithClient(cli),
	)
}
