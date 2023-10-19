layui.define(["layer"], function(exports){
    var layer = layui.layer
        , $ = layui.$;

    var login = {
        user_id: "",
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
                $(".user").html(this.user_id);
            }
        }
    };
    exports('login', login);
});

