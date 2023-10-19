layui.define(["layer"], function(exports){
    var layer = layui.layer
        , $ = layui.$;

    var base = {
        BaseInfo: {},
        load_info: function(){
            var me = this;
            $.ajax({
                url: API_HAOBASE_HOST+ "/api/v1/base/varieties/config",
                type: "get",
                data: {
                    symbol: CURRENT_SYMBOL,
                },
                dataType: "json",
                contentType: "application/json",
                success: function (d) {
                    console.log(d);
                    if(d.ok){
                        me.BaseInfo = d.data;
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
