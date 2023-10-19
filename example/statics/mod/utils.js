layui.define(['laytpl'], function(exports){
    var laytpl = layui.laytpl;
    var $ = layui.$;
    laytpl.config({open: '{%',close: '%}'});

    var obj = {
        formatTime: function(t) {
            var d = new Date(parseInt(t));
            return d.getFullYear() + '-' + (d.getMonth() + 1) + '-' + d.getDate() + ' ' + d.getHours() + ':' + d.getMinutes() + ':' + d.getSeconds();
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
        }
    };
    
    
    exports('utils', obj);
});

