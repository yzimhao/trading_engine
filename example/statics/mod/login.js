layui.define(["layer"], function(exports){
    var layer = layui.layer
        , $ = layui.$;

    var login = {
        user_id: "",
        show_setting_user_id: function(){
            var me = this;
            layer.prompt({
                formType: 0,
                placeholder: "10000",
                title: '请设置您的用户名，只允许数字和字母,长度4~10',
                area: ['400px', '50px'] // 自定义文本域宽高
            }, function(value, index, elem){
                var pp = new RegExp(/^[a-z0-9]{4,10}$/);
                if (pp.test(value)) {
                    me.user_id = value;
                    Cookies.set("user_id", value, { expires: 7, path: '' });
                    layer.close(index); // 关闭层
                    window.location.reload();
                }else{
                    layer.msg("用户名不符合规则");
                }
            });
        },
        logout: function(){
            Cookies.remove("user_id");
            window.location.reload();
        },

        
        init: function(){
            var me = this;
            if(!Cookies.get("user_id")) {
                this.show_setting_user_id();
            }else{
                this.user_id = Cookies.get("user_id");
                $(".user").html(this.user_id);
                $(".logout").on("click", function(){
                    me.logout();
                }).show();
            }
        }
    };
    exports('login', login);
});

