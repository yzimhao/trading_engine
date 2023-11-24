package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gomodule/redigo/redis"
	"github.com/gookit/goutil/fsutil"
	"github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/utils/app/config"

	"xorm.io/xorm"
	"xorm.io/xorm/log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Version   = "v0.0.0"
	Goversion = ""
	Commit    = ""
	Build     = "2023-10-31"

	Logger *logrus.Logger

	//
	redisPool *redis.Pool
	database  *xorm.Engine
)

func ShowVersion() {
	fmt.Println("version:", Version)
	fmt.Println("go:", Goversion)
	fmt.Println("commit:", Commit)
	fmt.Println("build:", Build)
}

func ConfigInit(config_file string, is_daemon bool) {
	if config_file != "" {
		viper.SetConfigFile(config_file)
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(&config.App); err != nil {
			logrus.Panicf("Error unmarshaling config: %s %s\n", config_file, err)
		}
	} else {
		config.App = &config.Configuration{}
	}

	//时区
	time.LoadLocation(config.App.Main.TimeZone)
	exename := filepath.Base(os.Args[0])
	initLogs(exename, is_daemon)

	if config.App.Main.Mode != config.ModeProd {
		Logger.Infof("当前运行在%s模式下", config.App.Main.Mode)
	}

}

func initLogs(logname string, isdaemon bool) {
	Logger = logrus.New()
	level, _ := logrus.ParseLevel(config.App.Main.LogLevel)
	Logger.SetLevel(level)
	// Logger.SetReportCaller(true)

	Logger.SetFormatter(&formatter.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	})

	output := []io.Writer{}
	if !isdaemon {
		output = append(output, os.Stdout)
	}
	if config.App.Main.LogPath != "" {
		err := fsutil.Mkdir(config.App.Main.LogPath, 0755)
		if err != nil {
			Logger.Fatal(err)
		}

		file := fmt.Sprintf("%s/%s_%s.log", config.App.Main.LogPath, logname, time.Now().Format("20060102150405"))
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			Logger.Fatal(err)
		}
		output = append(output, f)
	}

	mw := io.MultiWriter(output...)
	Logger.SetOutput(mw)
}

func RedisInit(addr, password string, db int) {
	if redisPool == nil {
		redisPool = &redis.Pool{
			MaxIdle:     50, //空闲数
			IdleTimeout: 300 * time.Second,
			MaxActive:   0, //最大数
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", addr)
				if err != nil {
					return nil, err
				}
				if password != "" {
					if _, err := c.Do("AUTH", password); err != nil {
						c.Close()
						return nil, err
					}
				}

				if _, err := c.Do("SELECT", db); err != nil {
					return nil, err
				}

				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	}

}

func DatabaseInit(driver, dsn string, show_sql bool, prefix string) {
	if database == nil {
		conn, err := xorm.NewEngine(driver, dsn)
		if err != nil {
			Logger.Panic(err)
		}

		if prefix != "" {
			// tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, prefix)
			// conn.SetTableMapper(tbMapper)
		}
		if show_sql {
			conn.ShowSQL(true)
		} else {
			conn.SetLogLevel(log.LOG_ERR)
		}

		conn.DatabaseTZ = time.Local
		conn.TZLocation = time.Local
		database = conn
	}
}

func Database() *xorm.Engine {
	return database
}

func RedisPool() *redis.Pool {
	return redisPool
}

func Deamon(pid string, logfile string) (*daemon.Context, *os.Process, error) {

	if err := fsutil.MkParentDir(pid); err != nil {
		return nil, nil, err
	}

	cntxt := &daemon.Context{
		PidFileName: pid,
		PidFilePerm: 0644,
		LogFileName: logfile,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		// Args:        ,
	}

	child, err := cntxt.Reborn()
	return cntxt, child, err
}

func Keepalive(app string, sec int) {
	go func() {
		for {
			data := make(map[string]any)
			data["version"] = Version

			_data, _ := json.Marshal(data)

			rdc := redisPool.Get()
			topic := fmt.Sprintf("keepalive.%s", app)
			rdc.Do("set", topic, _data)
			rdc.Do("expire", topic, sec+5)
			time.Sleep(time.Second * time.Duration(sec))
			rdc.Close()
		}
	}()
}

func KeepaliveInfo() []string {
	rdc := redisPool.Get()
	defer rdc.Close()

	keys, _ := scanKeys(rdc, 0, "keepalive.*")
	return keys
}

func scanKeys(conn redis.Conn, cursor int, pattern string) ([]string, error) {
	const count = 100 // 一次扫描的数量

	var keys []string
	for {
		values, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", pattern, "COUNT", count))
		if err != nil {
			return nil, err
		}

		values, err = redis.Scan(values, &cursor, &keys)
		if err != nil {
			return nil, err
		}

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}
