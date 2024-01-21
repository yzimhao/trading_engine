package app

import (
	"fmt"
	"io"
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gomodule/redigo/redis"
	"github.com/gookit/goutil/fsutil"
	"github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yzimhao/xormlog"

	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Version   = "v0.0.0"
	Goversion = ""
	Commit    = ""
	Build     = "0000-00-00"

	Logger *logrus.Logger

	//
	dbPrefix  = ""
	redisPool *redis.Pool
	database  *xorm.Engine
	runDaemon bool
)

func ShowVersion() {
	fmt.Println("version:", Version)
	fmt.Println("go:", Goversion)
	fmt.Println("commit:", Commit)
	fmt.Println("build:", Build)
}

func ConfigInit(config_file string, conf any) {
	if config_file != "" {
		viper.SetConfigFile(config_file)
		err := viper.ReadInConfig()
		if err != nil {
			logrus.Fatal(err)
		}

		if conf != nil {
			if err := viper.Unmarshal(&conf); err != nil {
				logrus.Fatalf("Error unmarshaling config: %s %s\n", config_file, err)
			}

		}
	}

}

func TimeZoneInit(timezone string) {
	time.LoadLocation(timezone)
}

func LogsInit(logname string, log_path string, log_level string, show bool) {
	Logger = logrus.New()
	level, _ := logrus.ParseLevel(log_level)
	Logger.SetLevel(level)

	Logger.SetFormatter(&formatter.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	})

	output := []io.Writer{}
	if show {
		output = append(output, os.Stdout)
	}
	if log_path != "" {
		err := fsutil.Mkdir(log_path, 0755)
		if err != nil {
			Logger.Fatal(err)
		}

		file := fmt.Sprintf("%s/%s_%s.log", log_path, logname, time.Now().Format("20060102150405"))
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

func DatabaseInit(driver, dsn string, show_sql bool, prefix string) (err error) {
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	if database == nil {
		fmt.Println("dsn:", dsn)
		conn, err := xorm.NewEngine(driver, dsn)
		if err != nil {
			return err
		}

		if prefix != "" {
			tbMapper := names.NewPrefixMapper(names.SnakeMapper{}, prefix)
			conn.SetTableMapper(tbMapper)
			dbPrefix = prefix
		}

		logctx := xormlog.NewLogCtx(Logger)
		conn.SetLogger(logctx)

		conn.DatabaseTZ = time.Local
		conn.TZLocation = time.Local

		if err := conn.Ping(); err != nil {
			return err
		}

		if show_sql {
			conn.ShowSQL(true)
			conn.SetLogLevel(log.LOG_INFO)
		} else {
			conn.SetLogLevel(log.LOG_ERR)
		}
		database = conn
	}
	return nil
}

func Database() *xorm.Engine {
	return database
}

func TablePrefix() string {
	return dbPrefix
}

func RedisPool() *redis.Pool {
	return redisPool
}

func SetDaemon(v bool) {
	runDaemon = v
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
