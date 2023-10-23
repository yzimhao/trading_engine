#### 开发

1. 安装依赖
    - mysql
    - redis

2. 配置文件
    修改cmd/config.toml文件中的mysql、redis服务地址和端口

3. 启动各模块
```
    cd cmd/ 
    go run haobase/main.go -c config.toml
    go run haomatch/main.go -c config.toml
    go run haoquote/main.go -c config.toml
```