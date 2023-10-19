layui.define(function(exports){
    

    var obj = {
        formatTime: function(t) {
            var d = new Date(parseInt(t));
            return d.getFullYear() + '-' + (d.getMonth() + 1) + '-' + d.getDate() + ' ' + d.getHours() + ':' + d.getMinutes() + ':' + d.getSeconds();
        }
    };
    
    
    exports('utils', obj);
});

