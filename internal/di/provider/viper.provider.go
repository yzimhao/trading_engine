package provider

import (
	"os"
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
	// for unit test
	if IsDevelopment() {
		_, current, _, _ := runtime.Caller(0)
		root := path.Dir(path.Dir(path.Dir(path.Dir(current))))
		logger.Sugar().Infof("root: %s", root)
		v.AddConfigPath(root)
	}

	logger.Sugar().Infof("config file: %s", v.ConfigFileUsed())
	return v
}

func Root() string {
	_, current, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(path.Dir(path.Dir(current))))
	return root
}

func IsDevelopment() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	if strings.Contains(exe, "__debug") {
		return true
	}
	return strings.Contains(exe, os.TempDir())
}
