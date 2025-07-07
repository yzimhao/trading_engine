import Highcharts from 'highcharts/highstock'; 


export const KChartManager = {
    chart: null,
    title: "K线",
    period: "d1",
    data: {
        "d1": [
            [Date.UTC(2023, 0, 1), 29.9, 71.5, 30.0, 50.0],
            [Date.UTC(2023, 0, 2), 71.5, 105.0, 60.0, 80.0],
            [Date.UTC(2023, 0, 3), 106.4, 129.2, 90.0, 120.0],
            [Date.UTC(2023, 0, 4), 129.2, 144.0, 100.0, 130.0],
            [Date.UTC(2023, 0, 5), 144.0, 176.0, 130.0, 150.0]
        ]
    },
    init(id) {
        const me = this;
        
        this.chart = Highcharts.stockChart(id, {
            title: {
                text: me.title
            },
            colors: ['#ff3232', '#00aa00'], // 红涨绿跌
            plotOptions: {
                candlestick: {
                    color: '#ff3232',    // 阴线（跌）颜色
                    upColor: '#00aa00',  // 阳线（涨）颜色
                    lineColor: '#333',   // K线边框颜色
                    lineWidth: 1
                }
            },
            series: [{
                type: 'candlestick',
                name: '示例数据',
                data: me.data[me.period],
                tooltip: {
                    valueDecimals: 2
                }
            }]
        });
    }
    
};