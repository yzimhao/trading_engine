layui.define(['laytpl'], function(exports){
    var laytpl = layui.laytpl;
    var $ = layui.$;
    laytpl.config({open: '{%',close: '%}'});

    var obj = {
        formatTime: function(t) {
            var d = new Date(parseInt(t));
            var month = (d.getMonth()+1) <10 ? "0"+d.getMonth()+1 : d.getMonth() +1;
            var date = d.getDate() <10 ? "0"+d.getDate() : d.getDate();
            var hours = d.getHours() <10 ? "0"+d.getHours() : d.getHours();
            var minutes = d.getMinutes() <10 ? "0"+d.getMinutes() : d.getMinutes();
            var seconds = d.getSeconds() <10 ? "0"+d.getSeconds() : d.getSeconds();

            return d.getFullYear() + '-' + month + '-' + date + ' ' + hours + ':' + minutes + ':' + seconds;
        },

        rendertradelog: function(data) {
            var logView = $(".trade-log .log"),
                logTpl = $("#trade-log-tpl").html();

            data['trade_at'] = this.formatTime(data.trade_at/1e6);
            laytpl(logTpl).render(data, function (html) {
                if ($(".log-item").length > 10) {
                    $(".log-item").last().remove();
                }
                logView.after(html);
            });
        },

        renderdepth: function (info) {
            var askTpl = $("#depth-ask-tpl").html()
                , askView = $(".depth-ask")
                , bidTpl = $("#depth-bid-tpl").html()
                , bidView = $(".depth-bid");


            laytpl(askTpl).render(info.asks.reverse(), function (html) {
                askView.html(html);
            });
            laytpl(bidTpl).render(info.bids, function (html) {
                bidView.html(html);
            });
        },

        rendermyorder: function(data) {
            var orderTpl = $("#myorder-tpl").html();

            data['create_time'] = this.formatTime(data.create_time/1e6);
            data['order_side'] = data['order_side'].toUpperCase();
            laytpl(orderTpl).render(data, function(html){
                $(".myorder-table-title").after(html);
            });
        }
    };
    
    
    exports('utils', obj);
});

