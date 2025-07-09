<template>
    <div style="position: fixed; top: 0; right: 0; border: 0; z-index:9999;">
        <a target="_blank" href="https://github.com/yzimhao/trading_engine" class="github-corner" aria-label="View source on GitHub">
            <svg width="80" height="80" viewBox="0 0 250 250"
                style="fill:#64CEAA; color:#fff; position: absolute; top: 0; border: 0; right: 0;" aria-hidden="true">
                <path d="M0,0 L115,115 L130,115 L142,142 L250,250 L250,0 Z"></path>
                <path
                    d="M128.3,109.0 C113.8,99.7 119.0,89.6 119.0,89.6 C122.0,82.7 120.5,78.6 120.5,78.6 C119.2,72.0 123.4,76.3 123.4,76.3 C127.3,80.9 125.5,87.3 125.5,87.3 C122.9,97.6 130.6,101.9 134.4,103.2"
                    fill="currentColor" style="transform-origin: 130px 106px;" class="octo-arm"></path>
                <path
                    d="M115.0,115.0 C114.9,115.1 118.7,116.5 119.8,115.4 L133.7,101.6 C136.9,99.2 139.9,98.4 142.2,98.6 C133.8,88.0 127.5,74.4 143.8,58.0 C148.5,53.4 154.0,51.2 159.7,51.0 C160.3,49.4 163.2,43.6 171.4,40.1 C171.4,40.1 176.1,42.5 178.8,56.2 C183.1,58.6 187.2,61.8 190.9,65.4 C194.5,69.0 197.7,73.2 200.1,77.6 C213.8,80.2 216.3,84.9 216.3,84.9 C212.7,93.1 206.9,96.0 205.4,96.6 C205.1,102.4 203.0,107.8 198.3,112.5 C181.9,128.9 168.3,122.5 157.7,114.1 C157.9,116.9 156.7,120.9 152.7,124.9 L141.0,136.5 C139.8,137.7 141.6,141.9 141.8,141.8 Z"
                    fill="currentColor" class="octo-body"></path>
            </svg>
        </a>
    </div>


  <view class="content">
    <view class="left">
        <view class="kline-chart" >
            <div id="kline" style="height: 500px;"></div>
        </view>
        <view class="user-form mtop10">
            <view class="user-assets">
                <view v-if="!user.isLogin" class="notlogin">
                    <view>
                        <text>游客，请登陆</text>
                    </view>
                    <view class="line1">
                        <text>用户名</text>
                        <uni-easyinput type="text" v-model="user.name" placeholder="随便输入一串字符串" />
                    </view>
                    <view class="line1">
                        <button type="primary" size="mini" @click="actionLogin">登陆</button>
                    </view>
                </view>
                <view v-if="user.isLogin">
                    <view class="line1">
                        <text>Hi, {{user.name}}</text>
                        <text class="logout" @click="actionLogout">退出</text>
                    </view>
                    <text>自助充值</text>
                    <view class="line1">
                        <text>Asset:</text> 
                        <uni-data-checkbox v-model="recharge.asset" :localdata='range.assetType'></uni-data-checkbox>
                    </view>
                    <view class="line1">
                        <text>Volume:</text> 
                        <uni-easyinput type="digit" v-model="recharge.volume" placeholder="1000" />
                    </view>
                    <view class="line1">
                        <button type="primary" size="mini" @click="actionRecharge">充值</button>
                    </view>
                </view>
            </view>
            <view class="user-area">
                <view class="sell">
                    <view class="line1">
                        <text class="item-title">类型</text>
                        <uni-data-checkbox v-model="range.sellOrderTypeVal" :localdata='range.orderType'></uni-data-checkbox>
                    </view>
                    <view class="line1" v-if="range.sellOrderTypeVal == 'MARKET'">
                        <text class="item-title">数量/金额</text>
                        <uni-data-checkbox v-model="range.sellQtyOrAmountVal" :localdata='range.qtyOrAmount'></uni-data-checkbox>
                    </view>
                    <view class="line1" v-if="range.sellOrderTypeVal == 'LIMIT'">
                        <text class="item-title">价格</text>
                        <uni-easyinput type="digit" v-model="order.sellPrice" placeholder="1.00" style="width: 200px;" />
                    </view>
                    <view class="line1" v-if="range.sellQtyOrAmountVal == 'qty'">
                        <text class="item-title">数量</text>
                        <uni-easyinput type="digit" v-model="order.sellQty" placeholder="10" style="width: 200px;"  />
                    </view>
                    <view class="line1" v-if="range.sellQtyOrAmountVal == 'amount'">
                        <text class="item-title">金额</text>
                        <uni-easyinput type="digit" v-model="order.sellAmount" placeholder="10" style="width: 200px;"  />
                    </view>
                    <view class="line1">
                        <button type="primary" size="mini" @click="actionSellOrder">卖出</button>
                        <button type="primary" size="mini" @click="actionSellOrderRand" v-if="range.sellOrderTypeVal=='LIMIT'">随机挂5单</button>
                    </view>
                    <view class="line1 asset-info">
                        <text>{{ current.targetSymbol.toUpperCase() }} 可用: {{ user.targetAsset.avail }}</text>
                        <text style="margin-left: 10px;">冻结: {{ user.targetAsset.freeze }}</text>
                        <text style="margin-left: 10px;">总数: {{ user.targetAsset.total }}</text>
                    </view>
                </view>
                <view class="buy">
                    <view class="line1">
                        <text class="item-title">类型</text>
                        <uni-data-checkbox v-model="range.buyOrderTypeVal" :localdata='range.orderType'></uni-data-checkbox>
                    </view>
                    <view class="line1" v-if="range.buyOrderTypeVal == 'MARKET'">
                        <text class="item-title">数量/金额</text>
                        <uni-data-checkbox v-model="range.buyQtyOrAmountVal" :localdata='range.qtyOrAmount'></uni-data-checkbox>
                    </view>
                    <view class="line1" v-if="range.buyOrderTypeVal == 'LIMIT'">
                        <text class="item-title">价格</text>
                        <uni-easyinput type="digit" placeholder="1.00" v-model="order.buyPrice" />
                    </view>
                    <view class="line1" v-if="range.buyQtyOrAmountVal == 'qty'">
                        <text class="item-title">数量</text>
                        <uni-easyinput type="digit" placeholder="10" v-model="order.buyQty" />
                    </view>
                    <view class="line1" v-if="range.buyQtyOrAmountVal == 'amount'">
                        <text class="item-title">金额</text>
                        <uni-easyinput type="digit" placeholder="10" v-model="order.buyAmount" />
                    </view>
                    <view class="line1">
                        <button type="primary" size="mini" @click="actionBuyOrder">买入</button>
                        <button type="primary" size="mini" @click="actionBuyOrderRand" v-if="range.buyOrderTypeVal=='LIMIT'">随机挂5单</button>
                    </view>
                    <view class="line1 asset-info">
                        <text>{{ current.baseSymbol.toUpperCase() }} 可用: {{ user.baseAsset.avail }}</text>
                        <text style="margin-left: 10px;">冻结: {{ user.baseAsset.freeze }}</text>
                        <text style="margin-left: 10px;">总数: {{ user.baseAsset.total }}</text>
                    </view>
                </view>
            </view>
        </view>
    </view>
    <view class="right">
        <view class="orderbook">
            <text class="orderbook-title">orderbook</text>
            
            <view class="ask">
                <uni-row v-for="(item,i) in depth.asks">
                    <uni-col :span="12">{{ item[0] }}</uni-col>
                    <uni-col :span="12">{{ item[1] }}</uni-col>
                </uni-row>
            </view>
        
            <view class="latest-price">
                <text>最新价格: {{ current.latestPrice }}</text>
                <text style="margin-left: 10px;">24H涨跌幅: {{ current.upRate24h }}%</text>
            </view>

            <view class="bid">
                <uni-row v-for="(item, i) in depth.bids">
                    <uni-col :span="12">{{ item[0] }}</uni-col>
                    <uni-col :span="12">{{ item[1] }}</uni-col>
                </uni-row>
            </view>
        </view>
        <view class="tradehistory mtop10">
            <uni-row>
                <uni-col :span="12">成交时间</uni-col>
                <uni-col :span="4">价格</uni-col>
                <uni-col :span="4">数量</uni-col>
                <uni-col :span="4">金额</uni-col>
            </uni-row>
            <view class="tradeRecord">
                <uni-row v-for="(item, i) in tradeRecords">
                    <uni-col :span="12">{{ formatTimestamp(item.trade_at) }}</uni-col>
                    <uni-col :span="4">{{ item.price }}</uni-col>
                    <uni-col :span="4">{{ item.qty }}</uni-col>
                    <uni-col :span="4">{{ item.amount }}</uni-col>
                </uni-row>
                
            </view>
        </view>
    </view>

    
  </view>

  <view class="footer">
    <view class="version">
        <uni-row>
            <text>version: {{ version.version }} build: {{ version.build }}</text>
        </uni-row>
        <uni-row>
            <text v-if="version.go.length > 0">go: {{ version.go }} commit: {{ version.commit }}</text>
        </uni-row>
    </view>
  </view>
</template>

<script setup>
const formatTimestamp = (value) => {
  // 1. 强制转换为字符串并处理空值
  const strValue = String(value ?? '0').padEnd(19, '0'); // 确保至少19位
  
  // 2. 安全提取毫秒部分（兼容不足19位的情况）
  const msPart = strValue.length >= 6 
    ? strValue.slice(-6, -3) 
    : '000';
  const milliseconds = msPart.padStart(3, '0');

  // 3. 转换时间戳
  const timestamp = parseInt(strValue) / 1e6;
  const date = new Date(timestamp);
  
  // 4. 格式化日期
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${year}/${month}/${day} ${hours}:${minutes}:${seconds}.${milliseconds}`;
}
</script>

<script>
import { request } from '@/common/request.js';
import { socketInit } from '@/common/websocket.js'; 
import { KChartManager } from '@/common/chart.js'; 

export default {
  data() {
    return {
        range: {
            orderType: [{"value": "LIMIT","text": "限价"	},{"value": "MARKET","text": "市价"}],
            qtyOrAmount: [{"value": "qty", "text": "数量"}, {"value": "amount", "text": "金额"}],
            sellQtyOrAmountVal: "qty",
            buyQtyOrAmountVal: "qty",
            sellOrderTypeVal: "LIMIT",
            buyOrderTypeVal: "LIMIT",
            assetType:[{"value":"BTC", "text": "BTC"}],
        },
        current: {
            symbol: "",
            baseSymbol: "",
            targetSymbol: "",
            latestPrice: "-",
            upRate24h:   "0"
        },
        recharge: {
            asset: "",
            volume: ""
        },
        order:{
            sellPrice: "",
            sellQty: "",
            sellAmount: "",
            buyPrice: "",
            buyQty:"",
            buyAmount:""
        },
        depth: {
            asks: [],
            bids:[]
        },
        tradeRecords: [],
        user: {
            name: "",
            isLogin: false,
            token: "",
            targetAsset: {
                avail: 0,
                freeze: 0,
                total: 0
            },
            baseAsset: {
                avail: 0,
                freeze: 0,
                total: 0
            }
        },
        version: {
            build: "",
            go:"",
            version:"",
            commit:""
        }
    }
  },
  
  onLoad(options) {
    const me = this;
    const user = uni.getStorageSync("user");
    if(user){
        this.user = user;
    }

    this.current.symbol = options.symbol;
    console.log(this.current);
    if(this.current.symbol){
        this.loadCurrentSymbol();
        this.loadDepth();
    }
    
    this.iniWebsocket();
    this.loadTradesRecord();
    this.loadAppVersion();
   
  },
  mounted(){
    KChartManager.init("kline", 2, 4);
    KChartManager.loadPeriodData(this.current.symbol, "m1");
  },
  methods: {
    iniWebsocket(){
        const me = this;
        const socket = socketInit();

        socket.onclose = (evt) => {
            console.log("websocket close.");
            setTimeout(() => {
                me.iniWebsocket();
            }, 3e3);
        };
        socket.onopen = () => {
            var msg = {
                "subscribe": [
                    "depth."+ me.current.symbol,
                    "trade." + me.current.symbol,
                    "price."+me.current.symbol,
                    "kline.m1."+me.current.symbol,
                    "market.24h."+me.current.symbol,
                    "market.28h."+me.current.symbol,
                    // "token."+ Cookies.get("jwt"),
                ],
                "unsubscribe":[
                    "MARKET.28h."+me.current.symbol,
                ]
            };
            console.log(JSON.stringify(msg));
            socket.send(JSON.stringify(msg));
        };

        socket.onmessage = (e) => {
            const msgs = e.data.split('\n');
            for (var i = 0; i < msgs.length; i++) {
                var msg = JSON.parse(msgs[i]);
                console.log("websocket message: " , msg);

                if (msg.type == "depth."+ me.current.symbol) {
                    me.depth.asks = msg.body.asks;
                    me.depth.bids = msg.body.bids;
                } else if (msg.type == "trade." + me.current.symbol) {
                    me.tradeRecords.push({
                        amount: msg.body.amount,
                        price: msg.body.price,
                        qty: msg.body.qty,
                        trade_at: msg.body.trade_at
                    })
                } else if (msg.type == "new_order."+ me.current.symbol) {
                    // var myorderView = $(".myorder"),
                    //     myorderTpl = $("#myorder-tpl").html();
                    
                    // var data = msg.body;
                    // data['create_time'] = utils.formatTime(data.create_time);
                    // laytpl(myorderTpl).render(data, function (html) {
                    //     if ($(".order-item").length > 30) {
                    //         $(".order-item").last().remove();
                    //     }
                    //     myorderView.after(html);
                    // });
                } else if (msg.type == "price."+me.current.symbol) {
                    me.current.latestPrice = msg.body.latest_price;
                    me.current.upRate24h = "-";
                } else if (msg.type =="kline.m1."+me.current.symbol) {
                    var data = msg.body;
                    KChartManager.addData({
                        timestamp: new Date(data[0]).getTime(),
                        open: +data[1],
                        high: +data[2],
                        low: +data[3],
                        close: +data[4],
                        volume: Math.ceil(+data[5]),
                    });
                }else if(msg.type=="market.24h."+me.current.symbol) {
                    // $(".price_p").html(msg.body.price_change_percent);
                }else if(msg.type =="order.cancel." +me.current.symbol) {
                    // var order_id = msg.body.order_id;
                    // layer.msg("订单 "+ order_id +" 取消成功");
                    // $(".myorder-item").each(function(){
                    //     if ($(this).attr("order-id")== order_id){
                    //         $(this).remove();
                    //     }
                    // })
                }
            }
        };
    },
    
    actionLogin () {
        const me = this;
        console.log(me.user);
        request("/api/v1/login", {
            "username": me.user.name,
            "password": "123456"
        }, "POST").then(res=>{
            console.log("token: ", res.data.token);
            console.log("expire: ", res.data.expire);
            if(res.data.token){
                me.user.token = res.data.token;
                me.user.isLogin  = true;
                uni.setStorageSync("user", me.user);
            }
            
        }).catch(err=>{
            console.log("/api/v1/login ", err);
        })
    },
    actionLogout(){
        uni.removeStorageSync('user');
        window.location.reload();
    },
    actionRecharge(){
        const me = this;
        request("/api/example/deposit", {"asset": me.recharge.asset, "volume": me.recharge.volume},  "GET").then(res=>{
            uni.showToast({
                title: "充值成功",
                icon: "none"
            });
            me.loadUserAssets();
        }).catch(err=>{
            console.log("/api/example/deposit ", err);
        })
    },
    actionSellOrder(){
        const me = this;
        let data = {
            "side": "SELL",
            "symbol": this.current.symbol
        };
        if(this.range.sellOrderTypeVal == "LIMIT") {
            data['order_type'] = "LIMIT";
            data['price'] = this.order.sellPrice;
        }else{
            data['order_type'] = "MARKET";
        }
        if(this.range.sellQtyOrAmountVal == "qty"){
            data['qty'] = this.order.sellQty;
        }else if(this.range.sellQtyOrAmountVal == "amount") {
            data['amount'] = this.order.sellAmount;
        }

        console.log("actionSellOrder: ", data);
        request("/api/v1/order", data, "POST").then(res=>{
            console.log("/api/v1/order ", data, res);
            uni.showToast({
                title: "挂单成功",
                icon: "none"
            });
            me.loadUserAssets();
        }).catch(err=>{
            console.log("/api/v1/order ", err);
            uni.showToast({title:err.data.msg, icon: "none"});
        })
    },
    actionBuyOrder(){
        const me = this;
        let data = {
            "side": "BUY",
            "symbol": this.current.symbol
        };
        if(this.range.buyOrderTypeVal == "LIMIT") {
            data['order_type'] = "LIMIT";
            data['price'] = this.order.buyPrice;
        }else{
            data['order_type'] = "MARKET";
        }
        if(this.range.buyQtyOrAmountVal == "qty"){
            data['qty'] = this.order.buyQty;
        }else if(this.range.buyQtyOrAmountVal == "amount") {
            data['amount'] = this.order.buyAmount;
        }

        console.log("actionBuyOrder: ", data);
        request("/api/v1/order", data, "POST").then(res=>{
            console.log("/api/v1/order ", data, res);
            uni.showToast({
                title: "挂单成功",
                icon: "none"
            });
            me.loadUserAssets();
        }).catch(err=>{
            console.log("/api/v1/order ", err);
            uni.showToast({title:err.data.msg, icon: "none"});
        })
    },


    getRandomInt(min, max) {
        min = Math.ceil(min);
        max = Math.floor(max);
        return Math.floor(Math.random() * (max - min + 1)) + min;
    },
    actionSellOrderRand(){
        const me = this;
        for(var i=0; i<5; i++) {
            let data = {
                "side": "SELL",
                "symbol": this.current.symbol,
                "order_type": "LIMIT",
                "price": 3 + Math.random(),
                "qty": me.getRandomInt(1, 5)
            };
            
            request("/api/v1/order", data, "POST").then(res=>{
                console.log("/api/v1/order ", data, res);
                me.loadUserAssets();
            }).catch(err=>{})
        }
    },
    actionBuyOrderRand(){
        const me = this;
        for(var i=0; i<5; i++) {
            let data = {
                "side": "BUY",
                "symbol": this.current.symbol,
                "order_type": "LIMIT",
                "price": 1 + Math.random(),
                "qty": me.getRandomInt(1, 5)
            };
            
            request("/api/v1/order", data, "POST").then(res=>{
                console.log("/api/v1/order ", data, res);
                me.loadUserAssets();
            }).catch(err=>{})
        }
    },


    loadAppVersion() {
        const me = this;
        request("/api/v1/version", {}, "GET", false)
        .then(res => {
            me.version = res.data;
            console.log(me.version);
        })
    },
    loadCurrentSymbol() {
        const me = this;
        request("/api/v1/product/"+this.current.symbol, {},  "GET").then(res=>{
            console.log("product info: ", res);
            me.current.baseSymbol = res.data.base.symbol;
            me.current.targetSymbol = res.data.target.symbol;
            me.range.assetType = [];
            
            me.recharge.asset = me.current.targetSymbol;
            me.recharge.volume = 1000;
            me.range.assetType.push({"value":me.current.targetSymbol, "text": me.current.targetSymbol.toUpperCase()});
            me.range.assetType.push({"value":me.current.baseSymbol, "text": me.current.baseSymbol.toUpperCase()});

            if(me.user.isLogin) {
                me.loadUserAssets();
            }
        }).catch(err=>{
            console.log("api/v1/product ", err);
        })
    },
    loadDepth(){
        const me = this;
        request("/api/v1/depth", {
            "symbol": me.current.symbol,
            "limit": 10
        },  "GET").then(res=>{
            console.log("/api/v1/depth ", res);
            const depth = res.data;
            me.depth.asks = depth.asks;
            me.depth.bids = depth.bids;
        }).catch(err=>{
            console.log("/api/v1/depth ", err);
        })
    },
    loadTradesRecord(){
        const me = this;
        request("/api/v1/trades", {
            "symbol": me.current.symbol,
            "limit": 10
        },  "GET").then(res=>{
            console.log("/api/v1/trades ", res);
            me.tradeRecords = res.data;
        }).catch(err=>{
            console.log("/api/v1/trades ", err);
        })
    },
    loadUserAssets() {
        const me = this;
        request("/api/v1/user/asset/query", {"symbols": me.current.baseSymbol + "," +me.current.targetSymbol},  "GET").then(res=>{
            console.log("/api/v1/user/asset/query ", res);
            const assets = res.data;
            for(var i=0; i<assets.length; i++) {
                if(me.current.baseSymbol == assets[i].symbol){
                    me.user.baseAsset = {
                        avail: assets[i].avail_balance,
                        freeze: assets[i].freeze_balance,
                        total: assets[i].total_balance
                    };
                }
                if(me.current.targetSymbol == assets[i].symbol){
                    me.user.targetAsset = {
                        avail: assets[i].avail_balance,
                        freeze: assets[i].freeze_balance,
                        total: assets[i].total_balance
                    };
                }
            }
        }).catch(err=>{
            console.log("/api/v1/user/asset/query ", err);
        })
    },
  },
}
</script>

<style lang="scss">
    @use '@/style/main.scss';
</style>
