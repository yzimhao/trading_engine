layui.define(["jquery"], function (exports) {
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

        format_num: function(d, spec) {
            const a = parseFloat(d);
            if(spec > -1) {
                return a.toFixed(spec);
            }
            return a;
        }
    };

    exports("utils", obj);
})