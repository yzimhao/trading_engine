layui.define(function(exports){
    var kchart = klinecharts.init('klinechart')
    var initKline = function () {
        // 初始化图表
        // 创建一个主图技术指标
        kchart.createIndicator('MA', false, { id: 'candle_pane' })
        // 创建一个副图技术指标VOL
        kchart.createIndicator('VOL')
        // 创建一个副图技术指标MACD
        
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
    };
    initKline()
    exports('kchart', kchart);
});