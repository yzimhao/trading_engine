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