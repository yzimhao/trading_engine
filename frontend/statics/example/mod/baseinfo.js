layui.define(["layer"], function(exports){
    var layer = layui.layer
        , $ = layui.$;

    var base = {
        cfg_info: {},
        load_info: function(){
            var me = this;
            $.ajax({
                url: "/api/v1/product/" + CURRENT_SYMBOL,
                type: "get",
                data: {
                    t: Date.now()
                },
                dataType: "json",
                contentType: "application/json",
                success: function (d) {
                    console.log("product:", d);
                    if(d.code==0){
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
