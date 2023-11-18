package config

var (
	App *Configuration
)

type ModeType string

const (
	ModeProd  ModeType = "prod"
	ModeDev   ModeType = "dev"
	ModeDebug ModeType = "debug"
	ModeDemo  ModeType = "demo"
)

type Configuration struct {
	Main     MainConfiguration     `mapstructure:"main"`
	Redis    RedisConfiguration    `mapstructure:"redis"`
	Database DatabaseConfiguration `mapstructure:"database"`
	Local    LocalConfiguration    `mapstructure:"local"`
	Haobase  HaobaseConfiguration  `mapstructure:"haobase"`
	Haomatch HaomatchConfiguration `mapstructure:"haomatch"`
	Haoquote HaoquoteConfiguration `mapstructure:"haoquote"`
	Haoadm   HaoadmConfiguration   `mapstructure:"haoadm"`
}

type MainConfiguration struct {
	Mode             ModeType `mapstructure:"mode"`
	LogLevel         string   `mapstructure:"log_level"`
	LogPath          string   `mapstructure:"log_path"`
	TimeZone         string   `mapstructure:"time_zone"`
	SecretKey        string   `mapstructure:"secret_key"`
	StaticServerName string   `mapstructure:"static_server_name"`
}

type RedisConfiguration struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Prefix   string `mapstructure:"prefix"`
}

type DatabaseConfiguration struct {
	Driver  string `mapstructure:"driver"`
	DSN     string `mapstructure:"dsn"`
	ShowSQL bool   `mapstructure:"show_sql"`
	Prefix  string `mapstructure:"prefix"`
}

type LocalConfiguration struct {
	Symbols []string `mapstructure:"symbols"`
}

type HaomatchConfiguration struct {
	Cache   string `mapstructure:"cache"`
	LogFile string `mapstructure:"log_file"`
	Listen  string `mapstructure:"listen"`
}

type HaoquoteConfiguration struct {
	Cache  string   `mapstructure:"cache"`
	Period []string `mapstructure:"period"`
	Listen string   `mapstructure:"listen"`
	Debug  bool     `mapstructure:"debug"`
	// LogFile string   `mapstructure:"log_file"`
}

type HaobaseConfiguration struct {
	Listen string `mapstructure:"listen"`
	Debug  bool   `mapstructure:"debug"`
	// LogFile string `mapstructure:"log_file"`
	InternalApiAllowIp []string `mapstructure:"internal_api_allow_ip"`
}

type HaoadmConfiguration struct {
	Listen   string `mapstructure:"listen"`
	Debug    bool   `mapstructure:"debug"`
	SiteName string `mapstructure:"site_name"`
	Readonly bool   `mapstructure:"readonly"`
}
