layui.define(['form',"baseinfo", 'kline', 'websocket'], function(exports){
    var baseinfo = layui.baseinfo;

    var layer = layui.layer //弹层
        // , table = layui.table //表格
        , form = layui.form
        , $ = layui.$
        , laytpl = layui.laytpl;

        laytpl.config({open: '{%',close: '%}'});

        function rendertradelog(data) {
            var logView = $(".trade-log .log"),
                logTpl = $("#trade-log-tpl").html();

            data['trade_at'] = formatTime(data.trade_at/1e6);
            laytpl(logTpl).render(data, function (html) {
                if ($(".log-item").length > 10) {
                    $(".log-item").last().remove();
                }
                logView.after(html);
            });
        }

        function renderdepth(info) {
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

    var obj = {
        bind: function(){
            form.on('select(order_type)', function (data) {
                if (data.value == "limit") {
                    $(".item-price").show();
                    $(".item-quantity").show();
                    $(".item-amount").hide();
                    $(".item-market-type").hide();
                    $(".qty-tips").hide();
                } else if (data.value == "market") {
                    $(".item-price").hide();
                    $(".item-market-type").show();
                    $(".qty-tips").show();
                }
                form.render('select');
            });
            form.on('radio(market-type)', function (data) {
                if (data.value == "q") {
                    $(".item-quantity").show();
                    $(".item-amount").hide();
                    $(".qty-tips").show();
                } else {
                    $(".item-quantity").hide();
                    $(".qty-tips").hide();
                    $(".item-amount").show();
                }
            });


            $(".opt").on("click", function () {
                var side = $(this).hasClass("sell") ? "sell" : "buy";
                var order_type = $("select[name='order_type']").val();
                var mtype = $("input[name='mtype']:checked").val();

                $.ajax({
                    url: API_HAOBASE_HOST+ "/api/v1/base/order/create",
                    type: "post",
                    dataType: "json",
                    contentType: "application/json",
                    beforeSend: function(r) {
                        r.setRequestHeader("token", Cookies.get("user_id"));
                    },
                    data: function () {
                        var data = {
                            symbol: CURRENT_SYMBOL,
                            side: side,
                            order_type: order_type,
                        };

                        if (order_type == "market") {
                            if (mtype == "market_qty") {
                                data.qty = $("input[name='quantity']").val();
                            } else {
                                data.amount = $("input[name='amount']").val();
                            }
                        } else {
                            data.price = $("input[name='price']").val();
                            data.qty = $("input[name='quantity']").val();
                        }

                        console.log(data);
                        return JSON.stringify(data)
                    }(),
                    success: function (d) {
                        if(d.ok){
                            layer.msg("下单成功");
                        }else{
                            layer.msg(d.reason);
                        }
                    }
                });
            });
        },

        load_depth_data: function(){
            $.get(API_HAOQUOTE_HOST + "/api/v1/quote/depth?symbol="+CURRENT_SYMBOL+"&limit=10", function(d){
                if(d.ok){
                    renderdepth(d.data);
                }
            });
        },
        load_tradelog_data: function(){
            $.get(API_HAOQUOTE_HOST + "/api/v1/quote/trans/record?symbol="+CURRENT_SYMBOL+"&limit=10", function (d) {
                if (d.ok) {
                    var data = d.data.reverse();
                    for(var i=0; i<data.length; i++){
                        rendertradelog(data[i]);
                    }

                }
            });
        },
        load_system_info: function(){
            $.get(API_HAOQUOTE_HOST + "/api/v1/quote/system", function(d){
                $(".version").html(d.version);
                $(".build").html(d.build);
            });
        },
        init: function(){
            console.log(baseinfo);
            this.bind();

        }
    };
    
    obj.init();
    exports('demo', obj);
});

