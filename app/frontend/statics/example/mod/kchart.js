layui.define(["baseinfo"], function(exports){
    var $ = layui.$;
    var kchart = klinecharts.init('klinechart');
    var baseinfo = layui.baseinfo;

    var kk = {
        init_chart: function(){
            // 初始化图表
            // 创建一个主图技术指标
            kchart.createIndicator('MA', false, { id: 'candle_pane' })
            // 创建一个副图技术指标VOL
            kchart.createIndicator('VOL')
            // 创建一个副图技术指标MACD
            kchart.setPriceVolumePrecision(baseinfo.cfg_info.price_precision, baseinfo.cfg_info.qty_precision);
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
                },
                
            });
        },
        load_kline_data: function(){
            $.get("/api/v1/quote/kline?symbol="+CURRENT_SYMBOL+"&interval=m1&limit=1000", function (d) {
                if (d.ok) {
                    var items = d.data.reverse();
                    var chartDataList = items.map(function (data) {
                        return {
                            timestamp: new Date(data[0]).getTime(),
                            open: +data[1],
                            high: +data[2],
                            low: +data[3],
                            close: +data[4],
                            volume: Math.ceil(+data[5]),
                        }
                    })
                    kchart.applyNewData(chartDataList)
                }
            });
        }
    }
    
    
    kk.init_chart();
    kk.load_kline_data();

    exports('kchart', kchart);
});