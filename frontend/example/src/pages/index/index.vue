<template>
  <view class="content">
    <view class="left">
        <view class="kline-chart">
            <text>chart loading...</text>
        </view>
        <view class="user-form mtop10">
            <view class="user-assets">
                <view  class="notlogin">
                    <view>
                        <text>游客，请登陆</text>
                    </view>
                    <view class="line1">
                        <text>用户名</text>
                        <uni-easyinput type="text" placeholder="随便输入一串字符串" style="width: 200px;" />
                    </view>
                    <view class="line1">
                        <button type="primary" size="mini">登陆</button>
                    </view>
                </view>
                <view v-if="false" class="line1">
                    <text>Hi, 132</text>
                    <text> [退出] </text>
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
export default {
  data() {
    return {
        range: {
            orderType: [{"value": 0,"text": "限价"	},{"value": 1,"text": "市价"}],
            sellOrderTypeVal: 0,
            buyOrderTypeVal: 0
        },
        data: {
            
        }
    }
  },
  onLoad() {},
  methods: {},
}
</script>

<style lang="scss">
    @use '@/style/main.scss';
</style>
