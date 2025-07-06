<template>
  <view class="content">
    <view class="left">
        <view class="kline-chart">
            <text>chart loading...</text>
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
                        <text class="logout">退出</text>
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
                    <view class="line1">
                        <text class="item-title">价格</text>
                        <uni-easyinput type="digit" placeholder="1.00" style="width: 200px;" />
                    </view>
                    <view class="line1">
                        <text class="item-title">数量</text>
                        <uni-easyinput type="digit" placeholder="10" style="width: 200px;"  />
                    </view>
                    <view class="line1">
                        <button type="primary" size="mini">卖出</button>
                        <text>可用BTC: 1000.00</text>
                    </view>
                </view>
                <view class="buy">
                    <view class="line1">
                        <text class="item-title">类型</text>
                        <uni-data-checkbox v-model="range.buyOrderTypeVal" :localdata='range.orderType'></uni-data-checkbox>
                    </view>
                    <view class="line1">
                        <text class="item-title">价格</text>
                        <uni-easyinput type="digit" placeholder="1.00" />
                    </view>
                    <view class="line1">
                        <text class="item-title">数量</text>
                        <uni-easyinput type="digit" placeholder="10" />
                    </view>
                    <view class="line1">
                        <button type="primary" size="mini">买入</button>
                        <text>可用USDT: 98.00</text>
                    </view>
                </view>
            </view>
        </view>
    </view>
    <view class="right">
        <view class="orderbook">
            <text class="orderbook-title">orderbook</text>
            
            <view class="ask">
                <uni-row>
                    <uni-col :span="12">1.03</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">1.04</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">1.01</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">1.02</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">1.03</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">1.04</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">1.01</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">1.02</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">1.03</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">1.04</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
            </view>
        
            <view class="latest-price">
                <text>最新价格: 1.00</text>
                <text style="margin-left: 10px;">24H涨跌幅: 2.0%</text>
            </view>

            <view class="bid">
                <uni-row>
                    <uni-col :span="12">0.99</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">0.99</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">0.99</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">0.99</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">0.99</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">0.99</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">0.99</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">0.99</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">0.99</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">0.99</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                
            </view>
        </view>
        <view class="tradehistory mtop10">
            <view>
                <uni-row>
                    <uni-col :span="12">成交时间</uni-col>
                    <uni-col :span="6">成交价格</uni-col>
                    <uni-col :span="6">成交量</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">2025-07-05 20:11:56.789</uni-col>
                    <uni-col :span="6">10.00</uni-col>
                    <uni-col :span="6">100</uni-col>
                </uni-row>
            </view>
        </view>
    </view>

    
  </view>
</template>

<script>
import {
    request
} from '@/common/request.js'

export default {
  data() {
    return {
        range: {
            orderType: [{"value": "limit","text": "限价"	},{"value": "market","text": "市价"}],
            sellOrderTypeVal: "limit",
            buyOrderTypeVal: "limit",
            assetType:[{"value":"BTC", "text": "BTC"}],
        },
        current: {
            symbol: "",
            baseSymbol: "",
            targetSymbol: ""
        },
        recharge: {
            asset: "",
            volume: ""
        },
        user: {
            name: "",
            isLogin: false,
            token: ""
        },
        data: {
            
        }
    }
  },
  onLoad(options) {
    const user = uni.getStorageSync("user");
    if(user){
        this.user = user;
    }

    this.current.symbol = options.symbol.toUpperCase();
    console.log(this.current);
    if(this.current.symbol){
        this.loadCurrentSymbol();
    }
  },
  methods: {
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
    actionRecharge(){
        const me = this;
        request("/api/example/deposit", {"asset": me.recharge.asset, "volume": me.recharge.volume},  "GET").then(res=>{
            console.log("/api/example/deposit info: ", res);
        }).catch(err=>{
            console.log("/api/example/deposit ", err);
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
        }).catch(err=>{
            console.log("api/v1/product ", err);
        })
    }
  },
}
</script>

<style lang="scss">
    @use '@/style/main.scss';
</style>
