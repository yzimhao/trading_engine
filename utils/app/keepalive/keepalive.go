package keepalive

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/gookit/goutil/arrutil"
	"github.com/yzimhao/trading_engine/types/redisdb"
	"github.com/yzimhao/trading_engine/utils"
)

var (
	single *Keepalive
)

type App struct {
	Name      string              `json:"name"`
	Pid       string              `json:"pid"`
	Version   string              `json:"version"`
	Extras    map[string][]string `json:"extras"`
	Runos     string              `json:"runos"`
	Runarch   string              `json:"runarch"`
	Hostname  string              `json:"hostname"`
	StartTime utils.Time          `json:"start_time"`
}

type Keepalive struct {
	id       string
	interval int
	rdc      *redis.Pool
	sync.Mutex

	app App
}

func NewKeepalive(rdc *redis.Pool, name, version string, interval int) *Keepalive {
	if single == nil {
		single = &Keepalive{
			id:       uuid.New().String(),
			rdc:      rdc,
			interval: interval,
			app: App{
				Name:    name,
				Version: version,
				Pid:     fmt.Sprintf("%d", os.Getpid()),
				Runos:   runtime.GOOS,
				Runarch: runtime.GOARCH,
				Hostname: func() string {
					n, _ := os.Hostname()
					return n
				}(),
				Extras:    make(map[string][]string),
				StartTime: utils.Time(time.Now()),
			},
		}
		single.run()
	}
	return single
}

func (k *Keepalive) run() {
	go func() {
		for {
			func() {
				rdc := k.rdc.Get()
				defer rdc.Close()

				k.Lock()
				defer k.Unlock()

				_data, _ := json.Marshal(k.app)
				topic := redisdb.Keepalive.Format(redisdb.Replace{"uuid": k.id})
				rdc.Do("set", topic, _data)
				rdc.Do("expire", topic, k.interval+3)
			}()
			time.Sleep(time.Second * time.Duration(k.interval))
		}
	}()
}

func SetExtras(key string, pp ...string) {
	single.Lock()
	defer single.Unlock()

	single.app.Extras[key] = append(single.app.Extras[key], pp...)
}

func HasExtrasKeyValue(key string, value string) bool {
	single.Lock()
	defer single.Unlock()

	if _, ok := single.app.Extras[key]; !ok {
		return false
	}

	if !arrutil.InStrings(value, single.app.Extras[key]) {
		return false
	}
	return true
}

func AppInfoTopic() []string {
	rdc := single.rdc.Get()
	defer rdc.Close()

	//todo 下面的函数redis key多了就遍历不出来了需要的key了
	keys, _ := utils.ScanRedisKeys(rdc, 0, "*keepalive.*")

	return keys
}
