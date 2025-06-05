package provider

import (
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewViper(logger *zap.Logger) *viper.Viper {
	var v = viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
	)

	v.AutomaticEnv()

	logger.Sugar().Infof("config file: %s", v.ConfigFileUsed())
	return v
}

func Root() string {
	_, current, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(path.Dir(path.Dir(current))))
	return root
}
