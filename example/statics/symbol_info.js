(function(){
    layui.use(['layer'], function () {
        var layer = layui.layer
        , $ = layui.$;

        var symbol = {
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
                            window["SymbolInfo"] = d.data;
                        }
                    }
                });
            },
            
            init: function(){
                this.load_info();
            }
        };
        symbol.init();
        
    })

})()