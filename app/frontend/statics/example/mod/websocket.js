layui.define(["layer", "utils", "kchart"],function(exports){
    var layer = layui.layer;
    var utils = layui.utils;
    var kchart = layui.kchart;
    var $ = layui.$;

    var socket = function () {
        if (window["WebSocket"]) {
            var protocol = window.location.protocol == "https:" ? "wss:" : "ws:";
            conn = new WebSocket("/ws");
            conn.onclose = function (evt) {
                layer.msg("<b>WebSocket Connection closed</b>");
                setTimeout(function () {
                    socket();
                }, 5e3);
            };

            conn.onopen = function(e){
                var msg = {
                    "subscribe": [
                        "depth."+CURRENT_SYMBOL,
                        "trade." + CURRENT_SYMBOL,
                        "price."+CURRENT_SYMBOL,
                        "kline.m1."+CURRENT_SYMBOL,
                        "market.24h."+CURRENT_SYMBOL,
                        "market.28h."+CURRENT_SYMBOL,
                        //订阅登陆用户相关消息，会对应多种消息格式
                        "token."+ Cookies.get("jwt"),
                    ],
                    "unsubscribe":[
                        "market.28h."+CURRENT_SYMBOL,
                    ]
                };
                console.log(JSON.stringify(msg));
                conn.send(JSON.stringify(msg));
            }

            
            conn.onmessage = function (evt) {
                var messages = evt.data.split('\n');
                for (var i = 0; i < messages.length; i++) {
                    var msg = JSON.parse(messages[i]);
                    console.log(msg);
                    if (msg.type == "depth."+CURRENT_SYMBOL) {
                        utils.renderdepth(msg.body);
                    } else if (msg.type == "trade." +CURRENT_SYMBOL) {
                        utils.rendertradelog(msg.body);
                    } else if (msg.type == "new_order."+ CURRENT_SYMBOL) {
                        var myorderView = $(".myorder"),
                            myorderTpl = $("#myorder-tpl").html();
                        
                        var data = msg.body;

                        data['create_time'] = utils.formatTime(data.create_time);
                        laytpl(myorderTpl).render(data, function (html) {
                            if ($(".order-item").length > 30) {
                                $(".order-item").last().remove();
                            }
                            myorderView.after(html);
                        });
                    } else if (msg.type == "price."+CURRENT_SYMBOL) {
                        $(".latest-price").html(msg.body.latest_price);
                    } else if (msg.type =="kline.m1."+CURRENT_SYMBOL) {
                        var data = msg.body;
                        kchart.updateData({
                            timestamp: new Date(data[0]).getTime(),
                            open: +data[1],
                            high: +data[2],
                            low: +data[3],
                            close: +data[4],
                            volume: Math.ceil(+data[5]),
                        });
                    }else if(msg.type=="market.24h."+CURRENT_SYMBOL) {
                        $(".price_p").html(msg.body.price_change_percent);
                    }else if(msg.type =="order.cancel." +CURRENT_SYMBOL) {
                        var order_id = msg.body.order_id;
                        layer.msg("订单 "+ order_id +" 取消成功");
                        $(".myorder-item").each(function(){
                            if ($(this).attr("order-id")== order_id){
                                $(this).remove();
                            }
                        })
                    }
                }
            };

            
        } else {
            layer.msg("<b>Your browser does not support WebSockets.</b>");
        }
    };

    socket();
    var obj = {
        init: function(){
            socket();
        }
    };
    exports("websocket", obj);
})