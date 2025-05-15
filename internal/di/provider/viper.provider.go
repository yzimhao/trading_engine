package provider

import (
	"strings"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	var v = viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
	)

	v.AutomaticEnv()

	return v
}
