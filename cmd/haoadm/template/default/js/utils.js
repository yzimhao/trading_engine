(function(){
    var $ = layui.$;

    var obj = {
        get: function(url, name){
            const params = {};
            const queryString = url.split('?')[1];
            if (queryString) {
                const keyValuePairs = queryString.split('&');
                keyValuePairs.forEach(pair => {
                    const [key, value] = pair.split('=');
                    params[key] = decodeURIComponent(value);
                });
            }
            return params[name];
        },

        format_timestamp: function(d){
            const date = new Date(d);
            // 获取各个时间部分
            const year = date.getFullYear();
            const month = String(date.getMonth() + 1).padStart(2, "0");
            const day = String(date.getDate()).padStart(2, "0");
            const hours = String(date.getHours()).padStart(2, "0");
            const minutes = String(date.getMinutes()).padStart(2, "0");
            const seconds = String(date.getSeconds()).padStart(2, "0");

            // 格式化为标准格式
            const formattedDate = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
            return formattedDate;
        },

        format_status: function(d) {
            if (d==0) {
                return '<i class="fa fa-check"></i>';
            }
            return '<i class="fa fa-close"></i>';
        },

        format_order_type: function(d){
            if(d=="limit") {
                return "限价";
            }else if(d=="market") {
                return "市价";
            }
        },

        format_order_side: function(d){
            if(d=="buy"){
                return "买入";
            }else{
                return "卖出";
            }
        },

        format_order_status: function(d){
            if(d==0) {
                return "未成交";
            }else if(d == 1){
                return "已成交";
            }else if(d==2) {
                return "已取消";
            }
        },

        format_num: function(d, spec) {
            const a = parseFloat(d);
            if(spec > -1) {
                return a.toFixed(spec);
            }
            return a;
        },

        open_url: function(title, url) {
            var index = layer.open({
                title: title,
                type: 2,
                shade: 0.2,
                maxmin:true,
                shadeClose: true,
                area: ['100%', '100%'],
                content: url,
            });
            $(window).on("resize", function () {
                layer.full(index);
            });
        }
    };



    window["utils"] = obj;
})()