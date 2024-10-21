package di

import (
	"github.com/duolacloud/broker-core"
	rocketmq "github.com/duolacloud/broker-rocketmq"
	"github.com/spf13/viper"
)

func NewBroker(v *viper.Viper) (broker.Broker, error) {
	v.SetDefault("broker.addrs", "127.0.0.1:9876")
	v.SetDefault("rocketmq.concurrency", 10)
	v.SetDefault("rocketmq.retry", 10)
	v.SetDefault("rocketmq.groupname", "default")

	bopts := make([]broker.Option, 0)

	addrs := v.GetString("broker.addrs")
	concurrency := v.GetInt("rocketmq.concurrency")
	retry := v.GetInt("rocketmq.retry")
	groupName := v.GetString("rocketmq.groupname")

	bopts = append(bopts, broker.Addrs(addrs))
	bopts = append(bopts, rocketmq.WithConsumeGoroutineNums(concurrency))
	bopts = append(bopts, rocketmq.WithRetry(retry))
	bopts = append(bopts, rocketmq.WithGroupName(groupName))

	b := rocketmq.NewBroker(bopts...)

	if err := b.Connect(); err != nil {
		return nil, err
	}

	return b, nil
}
