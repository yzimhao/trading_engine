package base

import (
	"github.com/redis/go-redis/v9"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/symbols"
	"xorm.io/xorm"
)

var (
	db  *xorm.Engine
	rdc *redis.Client
)

func Init(_db *xorm.Engine, _rdc *redis.Client) {
	db = _db
	rdc = _rdc
	symbols.Init(db, rdc)
}

func DB() *xorm.Engine {
	return db
}

func RDC() *redis.Client {
	return rdc
}
