(function(){
    layui.use(['layer'], function () {
        var layer = layui.layer
        , $ = layui.$;

        var login = {
            user_id: "",
            load_user_assets: function(){
                var me = this;
                $.ajax({
                    url: API_HAOBASE_HOST+ "/api/v1/base/assets",
                    type: "get",
                    data: {
                        symbols: SymbolInfo.target.symbol + "," +SymbolInfo.base.symbol,
                    },
                    dataType: "json",
                    contentType: "application/json",
                    beforeSend: function(r) {
                        r.setRequestHeader("token", me.user_id);
                    },
                    success: function (d) {
                        console.log(d);
                    }
                });
            },
            show_setting_user_id: function(){
                var me = this;
                layer.prompt({
                    formType: 0,
                    value: "1000",
                    title: '第一次，请设置您的用户ID',
                    area: ['400px', '50px'] // 自定义文本域宽高
                }, function(value, index, elem){
                    me.user_id = value
                    Cookies.set("user_id", value, { expires: 7, path: '' });
                    layer.close(index); // 关闭层
                    window.location.reload();
                });
            },
            init: function(){
                if(!Cookies.get("user_id")) {
                    this.show_setting_user_id();
                }else{
                    this.user_id = Cookies.get("user_id");
                    layui.$(".user").html(this.user_id);
                    this.load_user_assets();
                }
            }
        };
        login.init();
    })

})()