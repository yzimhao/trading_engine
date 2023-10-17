package app

import (
	"fmt"
	"io"
	"os"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gomodule/redigo/redis"
	"github.com/gookit/goutil/fsutil"
	"github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Mode string

const (
	ModeProd  Mode = "prod"
	ModeDev   Mode = "dev"
	ModeDebug Mode = "debug"
	ModeDemo  Mode = "demo"
)

var (
	Version        = "v0.0.0"
	Goversion      = ""
	Commit         = ""
	Build          = ""
	RunMode   Mode = ModeProd

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

func ConfigInit(fp string) {
	viper.SetConfigFile(fp)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	//
	time.LoadLocation(viper.GetString("main.time_zone"))

	RunMode = Mode(viper.GetString("main.mode"))
}

func Cstring(key string) string {
	return viper.GetString(key)
}

func Cbool(key string) bool {
	return viper.GetBool(key)
}

func Cint(key string) int {
	return viper.GetInt(key)
}

func RedisInit(addr, password string, db int) {
	if redisPool == nil {
		redisPool = &redis.Pool{
			MaxIdle:     10, //空闲数
			IdleTimeout: 240 * time.Second,
			MaxActive:   20, //最大数
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

func LogsInit(fn string, is_daemon bool) {
	level, _ := logrus.ParseLevel(viper.GetString("main.log_level"))
	logrus.SetLevel(level)
	// logrus.SetReportCaller(true)

	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	})

	output := []io.Writer{}
	if !is_daemon {
		output = append(output, os.Stdout)
	}
	if viper.GetString("main.log_path") != "" {
		save_path := viper.GetString("main.log_path")
		err := fsutil.Mkdir(save_path, 0755)
		if err != nil {
			logrus.Fatal(err)
		}

		file := fmt.Sprintf("%s/%s_%d.log", save_path, fn, time.Now().Unix())
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			logrus.Fatal(err)
		}
		output = append(output, f)
	}

	mw := io.MultiWriter(output...)
	logrus.SetOutput(mw)
}

func DatabaseInit(driver, dsn string, show_sql bool) {
	if database == nil {
		// dsn := viper.GetString("database.dsn")
		// driver := viper.GetString("database.driver")
		conn, err := xorm.NewEngine(driver, dsn)
		if err != nil {
			logrus.Panic(err)
		}

		// tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "ex_")
		// conn.SetTableMapper(tbMapper)
		// if viper.GetBool("database.show_sql") {
		// 	conn.ShowSQL(true)
		// }
		if show_sql {
			conn.ShowSQL(true)
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
