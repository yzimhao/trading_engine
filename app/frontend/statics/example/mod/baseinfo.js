layui.define(["layer"], function(exports){
    var layer = layui.layer
        , $ = layui.$;

    var base = {
        cfg_info: {},
        load_info: function(){
            var me = this;
            $.ajax({
                url: "/api/v1/base/exchange_info",
                type: "get",
                data: {
                    symbol: CURRENT_SYMBOL,
                    t: Date.now()
                },
                dataType: "json",
                contentType: "application/json",
                success: function (d) {
                    console.log("exchange info:", d);
                    if(d.ok){
                        me.cfg_info = d.data;
                    }
                }
            });
        },
        
        init: function(){
            this.load_info();
        }
    };
    base.init();
    
    
    exports('baseinfo', base);
});
