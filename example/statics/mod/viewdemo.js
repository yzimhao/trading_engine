layui.define(['form',"baseinfo", 'utils', 'kchart', 'websocket','login'], function(exports){
    var baseinfo = layui.baseinfo;
    var login = layui.login;

    var layer = layui.layer //弹层
        , form = layui.form
        , utils = layui.utils
        , $ = layui.$;
        

    var obj = {
        bind: function(){
            var me = this;
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
                            
                            me.load_assets();
                            me.load_order_unfinished();
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
                    utils.renderdepth(d.data);
                }
            });
        },
        load_tradelog_data: function(){
            $.get(API_HAOQUOTE_HOST + "/api/v1/quote/trans/record?symbol="+CURRENT_SYMBOL+"&limit=10", function (d) {
                if (d.ok) {
                    var data = d.data.reverse();
                    for(var i=0; i<data.length; i++){
                        utils.rendertradelog(data[i]);
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
        load_assets: function(){
            $.ajax({
                url: API_HAOBASE_HOST+ "/api/v1/base/assets",
                type: "get",
                beforeSend: function(r) {
                    r.setRequestHeader("token", login.user_id);
                },
                data:{
                    symbols: baseinfo.cfg_info.target.symbol+ "," + baseinfo.cfg_info.base.symbol
                },
                success: function (d) {
                    console.log("load_assets: ", d);
                    if(!d.ok) {
                        layer.msg(d.reason);
                        return;
                    }

                    var html = [];
                    for(var i=0; i<d.data.length; i++){
                        html.push(" " + d.data[i].symbol.toUpperCase() + ":" + d.data[i].avail);
                    }
                    $(".assets .list").html(html.join(" "));
                }
            });
        },
        load_order_unfinished: function(){
            $.ajax({
                url: API_HAOBASE_HOST+ "/api/v1/base/order/unfinished",
                type: "get",
                beforeSend: function(r) {
                    r.setRequestHeader("token", login.user_id);
                },
                data:{
                    symbol: baseinfo.cfg_info.symbol,
                    limit: 4
                },
                success: function (d) {
                    console.log("load_order_unfinished: ", d);
                    if(d.ok){
                        $(".myorder-item").remove();
                        d.data = d.data.reverse();
                        for(var i=0; i<d.data.length; i++){
                            utils.rendermyorder(d.data[i]);
                        }
                    }
                }
            });
        },
        init: function(){
            console.log(baseinfo);
            login.init();
            this.bind();
            this.load_system_info();
            this.load_depth_data();
            this.load_tradelog_data();
            this.load_assets();
            this.load_order_unfinished();
        }
    };
    
    obj.init();
    exports('viewdemo', obj);
});

