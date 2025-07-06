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
                    <view class="line1">
                        <text class="item-title">数量/金额</text>
                        <uni-data-checkbox v-model="range.sellQtyOrAmountVal" :localdata='range.qtyOrAmount'></uni-data-checkbox>
                    </view>
                    <view class="line1" v-if="range.sellOrderTypeVal == 'limit'">
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
                    <view class="line1">
                        <text class="item-title">数量/金额</text>
                        <uni-data-checkbox v-model="range.buyQtyOrAmountVal" :localdata='range.qtyOrAmount'></uni-data-checkbox>
                    </view>
                    <view class="line1" v-if="range.buyOrderTypeVal == 'limit'">
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
                <text>最新价格: 1.00</text>
                <text style="margin-left: 10px;">24H涨跌幅: 2.0%</text>
            </view>

            <view class="bid">
                <uni-row v-for="(item, i) in depth.bids">
                    <uni-col :span="12">{{ item[0] }}</uni-col>
                    <uni-col :span="12">{{ item[1] }}</uni-col>
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
import { socketInit } from '@/common/websocket.js'; 

export default {
  data() {
    return {
        range: {
            orderType: [{"value": "limit","text": "限价"	},{"value": "market","text": "市价"}],
            qtyOrAmount: [{"value": "qty", "text": "数量"}, {"value": "amount", "text": "金额"}],
            sellQtyOrAmountVal: "qty",
            buyQtyOrAmountVal: "qty",
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
        data: {
            
        }
    }
  },
  onLoad(options) {
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
                    "market.28h."+me.current.symbol,
                ]
            };
            console.log(JSON.stringify(msg));
            socket.send(JSON.stringify(msg));
        };

        socket.onmessage = (e) => {
            const msg = e.data.split('\n');
            for (var i = 0; i < msg.length; i++) {
                var data = JSON.parse(msg[i]);
                console.log("websocket message: " ,data);

                if (data.type == "depth."+ me.current.symbol) {
                    me.depth.asks = data.body.asks;
                    me.depth.bids = data.body.bids;
                } else if (data.type == "trade." + me.current.symbol) {
                    // utils.rendertradelog(data.body);
                } else if (data.type == "new_order."+ me.current.symbol) {
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
                } else if (data.type == "price."+me.current.symbol) {
                    $(".latest-price").html(msg.body.latest_price);
                } else if (data.type =="kline.m1."+me.current.symbol) {
                    // var data = msg.body;
                    // kchart.updateData({
                    //     timestamp: new Date(data[0]).getTime(),
                    //     open: +data[1],
                    //     high: +data[2],
                    //     low: +data[3],
                    //     close: +data[4],
                    //     volume: Math.ceil(+data[5]),
                    // });
                }else if(data.type=="market.24h."+me.current.symbol) {
                    // $(".price_p").html(msg.body.price_change_percent);
                }else if(data.type =="order.cancel." +me.current.symbol) {
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
            "side": "ask",
            "symbol": this.current.symbol
        };
        if(this.range.sellOrderTypeVal == "limit") {
            data['order_type'] = "limit";
            data['price'] = this.order.sellPrice;
        }else{
            data['order_type'] = "market";
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
            "side": "bid",
            "symbol": this.current.symbol
        };
        if(this.range.buyOrderTypeVal == "limit") {
            data['order_type'] = "limit";
            data['price'] = this.order.buyPrice;
        }else{
            data['order_type'] = "market";
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
