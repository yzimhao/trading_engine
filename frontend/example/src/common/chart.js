import { request } from '@/common/request.js';
import "@/common/klinecharts.min.js";


const upColor = '#00da3c';
const downColor = '#ec0000';

export const KChartManager = {
    period: "1m",  // 使用标准周期格式
    kchart: null,
    dataMap: {},
    
    init(id, price_precision, qty_precision) {
        const kchart = klinecharts.init(document.getElementById(id));
        // 初始化图表
        // 创建一个主图技术指标
        kchart.createIndicator('MA', false, { id: 'candle_pane' })
        // 创建一个副图技术指标VOL
        kchart.createIndicator('VOL')
        // 创建一个副图技术指标MACD
        kchart.setPriceVolumePrecision(price_precision, qty_precision);
        kchart.setBarSpace(10);

        kchart.setStyles({
            grid: {
                show: true,
                horizontal: {
                show: true,
                size: 1,
                color: '#EDEDED',
                style: 'dashed',
                dashedValue: [2, 2]
                },
                vertical: {
                show: true,
                size: 1,
                color: '#EDEDED',
                style: 'dashed',
                dashedValue: [2, 2]
                }
            }
        });
        this.kchart = kchart;
    },

    addData(point) {
        this.kchart.updateData(point);
    },

    loadPeriodData(symbol, period) {
        const me = this;
        this.period = period;
        request("/api/v1/klines", {
            "symbol": symbol,
            "period": period
        }, "GET").then(res => {
            const items = res.data.reverse();
            const chartDataList = items.map(function (data) {
                return {
                    timestamp: new Date(data[0]).getTime(),
                    open: +data[1],
                    high: +data[2],
                    low: +data[3],
                    close: +data[4],
                    volume: Math.ceil(+data[5]),
                }
            })
            me.kchart.applyNewData(chartDataList)
        }).catch(err => {
            console.error("加载K线数据失败:", err);
        });
    }
};