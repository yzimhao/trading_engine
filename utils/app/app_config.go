package app

import "github.com/spf13/viper"

func Cstring(key string) string {
	return viper.GetString(key)
}

func Cbool(key string) bool {
	return viper.GetBool(key)
}

func Cint(key string) int {
	return viper.GetInt(key)
}

func CstringSlice(key string) []string {
	return viper.GetStringSlice(key)
}
