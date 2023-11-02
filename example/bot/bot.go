package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

const (
	BOTSELL = "demobot1"
	BOTBUY  = "demobot2"
)

type bot struct {
	symbol       string
	now_price    string
	remote_price string
}

func StartBot() {
	auto_deposit("bot-test-001", BOTSELL, "usd", "1000000000000")
	auto_deposit("bot-test-002", BOTBUY, "jpy", "1000000000000")

	b1 := bot{symbol: "usdjpy"}
	go b1.run()
}

func (b *bot) run() {
	for {
		b.get_now_price()
		update := b.get_remote_price()
		app.Logger.Infof("%v", b)
		if update {
			b.auto_buy(BOTBUY, b.remote_price, "30")
			b.auto_sell(BOTSELL, b.remote_price, "30")
			b.auto_depth()
		}
		sec := 5 + rand.Int63n(20)
		app.Logger.Infof("sleep: %d sec", sec)
		time.Sleep(time.Second * time.Duration(sec))
	}
}

func (b *bot) get_remote_price() (isupdate bool) {
	latest := get_remote_price(b.symbol)
	if latest != "0" && b.remote_price != latest {
		b.remote_price = latest
		return true
	}
	return false
}

func (b *bot) get_now_price() {
	b.now_price = get_latest_price(b.symbol)
}

func (b *bot) auto_depth() {
	depth := get_depth(b.symbol)
	if len(depth["asks"]) < 10 || utils.D(depth["asks"][0][0]).Cmp(utils.D(b.remote_price)) > 1 {

		for i := 0; i < len(depth["asks"]); i++ {
			float := rand.Float64()
			price := utils.D(b.remote_price).Add(decimal.NewFromFloat(float))
			b.auto_sell(BOTSELL, price.String(), "0.01")
		}
	}

	if len(depth["bids"]) < 10 || utils.D(depth["bids"][0][0]).Cmp(utils.D(b.remote_price)) > 1 {
		for i := 0; i < len(depth["bids"]); i++ {
			float := rand.Float64()
			price := utils.D(b.remote_price).Sub(decimal.NewFromFloat(float))
			b.auto_buy(BOTBUY, price.String(), "0.01")
		}
	}
}

type order_create_request_args struct {
	Symbol    string                 `json:"symbol" binding:"required"`
	Side      trading_core.OrderSide `json:"side" binding:"required"`
	OrderType trading_core.OrderType `json:"order_type" binding:"required"`
	Price     string                 `json:"price" example:"1.00"`
	Quantity  string                 `json:"qty" example:"12"`
	Amount    string                 `json:"amount" example:"100.00"`
}

func (b *bot) auto_buy(user, price, qty string) {
	var buy order_create_request_args
	buy.Symbol = b.symbol
	buy.Side = trading_core.OrderSideBuy
	buy.OrderType = trading_core.OrderTypeLimit
	buy.Price = price
	buy.Quantity = qty
	data, _ := json.Marshal(buy)
	open_order(user, data)
}
func (b *bot) auto_sell(user, price, qty string) {
	var buy order_create_request_args
	buy.Symbol = b.symbol
	buy.Side = trading_core.OrderSideSell
	buy.OrderType = trading_core.OrderTypeLimit
	buy.Price = price
	buy.Quantity = qty
	data, _ := json.Marshal(buy)
	open_order(user, data)
}

func open_order(user string, data []byte) {
	// 要发送的数据，可以是 JSON 数据、表单数据等
	// data := []byte(`{"name": "John", "age": 30}`)

	// 创建 HTTP 请求头
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("Token", user)

	// 发送 POST 请求
	url := fmt.Sprintf("http:%s/api/v1/base/order/create", viper.GetString("api.haobase_host"))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		app.Logger.Warnf("HTTP request creation failed: %s", err.Error())
		return
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.Logger.Warnf("HTTP POST request failed: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	// 读取响应
	if resp.Status == "200 OK" {
		// 在这里处理响应数据
		body, _ := ioutil.ReadAll(resp.Body)
		app.Logger.Infof("body: %s", body)
	} else {
		app.Logger.Warnf("request failed with status: %s", resp.Status)
	}
}

func get_latest_price(symbol string) string {
	// 创建 HTTP 请求头
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	url := fmt.Sprintf("http:%s/api/v1/quote/price?symbol=%s", viper.GetString("api.haoquote_host"), symbol)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		app.Logger.Warnf("HTTP request creation failed: %s", err.Error())
		return "0"
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.Logger.Warnf("HTTP GET request failed: %s", err.Error())
		return "0"
	}
	defer resp.Body.Close()

	type response_data struct {
		Ok   bool              `json:"ok"`
		Data map[string]string `json:"data"`
	}

	// 读取响应
	if resp.Status == "200 OK" {
		// 在这里处理响应数据
		body, _ := ioutil.ReadAll(resp.Body)
		var data response_data
		json.Unmarshal(body, &data)
		return data.Data[symbol]
	} else {
		app.Logger.Warnf("HTTP GET request failed with status: %s", resp.Status)
	}
	return "0"
}

func get_depth(symbol string) map[string][2][]string {
	// 创建 HTTP 请求头
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	// 发送 POST 请求
	url := fmt.Sprintf("http:%s/api/v1/quote/depth?symbol=%s&limit=10", viper.GetString("api.haoquote_host"), symbol)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		app.Logger.Warnf("HTTP request creation failed: %s", err.Error())
		return map[string][2][]string{}
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.Logger.Warnf("HTTP GET request failed: %s", err.Error())
		return map[string][2][]string{}
	}
	defer resp.Body.Close()

	type response_data struct {
		Ok   bool                   `json:"ok"`
		Data map[string][2][]string `json:"data"`
	}

	// 读取响应
	if resp.Status == "200 OK" {
		// 在这里处理响应数据
		body, _ := ioutil.ReadAll(resp.Body)
		var data response_data
		json.Unmarshal(body, &data)
		return data.Data
	} else {
		app.Logger.Warnf("HTTP GET request failed with status: %s", resp.Status)
	}
	return map[string][2][]string{}
}

func get_remote_price(symbol string) string {
	// 创建 HTTP 请求头
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	// 发送 POST 请求
	url := fmt.Sprintf("https://finance.pae.baidu.com/vapi/v1/getquotation?group=huilv_minute&need_reverse_real=0&code=%s&finClientType=pc", symbol)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		app.Logger.Warnf("HTTP request creation failed: %s", err.Error())
		return "0"
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.Logger.Warnf("HTTP GET request failed: %s", err.Error())
		return "0"
	}
	defer resp.Body.Close()

	type response_data struct {
		Result struct {
			Cur struct {
				Price string `json:"price"`
			} `json:"cur"`
		} `json:"Result"`
	}

	// 读取响应
	if resp.Status == "200 OK" {
		// 在这里处理响应数据
		body, _ := ioutil.ReadAll(resp.Body)
		var data response_data
		json.Unmarshal(body, &data)
		return data.Result.Cur.Price
	} else {
		app.Logger.Warnf("request failed with status: %s", resp.Status)
	}
	return "0"
}

func auto_deposit(order_id, user_id, symbol, amount string) {

	type req_deposit_withdraw_args struct {
		OrderId string `json:"order_id" binding:"required"`
		UserId  string `json:"user_id" binding:"required"`
		Symbol  string `json:"symbol" binding:"required"`
		Amount  string `json:"amount" binding:"required"`
	}

	args := req_deposit_withdraw_args{
		OrderId: order_id,
		UserId:  user_id,
		Symbol:  symbol,
		Amount:  amount,
	}

	data, _ := json.Marshal(args)
	headers := make(http.Header)
	// 发送 POST 请求
	url := fmt.Sprintf("http:%s/api/v1/internal/deposit", viper.GetString("api.haobase_host"))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		app.Logger.Warnf("HTTP request creation failed: %s", err.Error())
		return
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.Logger.Warnf("HTTP POST request failed: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	// 读取响应
	if resp.Status == "200 OK" {
		// 在这里处理响应数据
		body, _ := ioutil.ReadAll(resp.Body)
		app.Logger.Infof("body: %s", body)
	} else {
		app.Logger.Warnf("request failed with status: %s", resp.Status)
	}
}
