<template>
  <view class="content">
    <view class="left">
        <view class="kline-chart"></view>
        <view class="user-form mtop10"></view>
    </view>
    <view class="right">
        <view class="orderbook">
            <text class="orderbook-title">orderbook</text>
            
            <view class="ask">
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
                    <uni-col :span="12">1.00</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                <uni-row>
                    <uni-col :span="12">0.99</uni-col>
                    <uni-col :span="12">1000</uni-col>
                </uni-row>
                
            </view>
        </view>
        <view class="tradehistory mtop10"></view>
    </view>

    
  </view>
</template>

<script>
export default {
  data() {
    return {
      title: 'Home',
    }
  },
  onLoad() {},
  methods: {},
}
</script>

<style lang="scss">
    @use '@/style/main.scss';
</style>
