docker监控脚本
================================

系统需求
--------------------------------
操作系统：Linux

主要逻辑
--------------------------------
获取docker daemon接口数据，解析返回结果，将key组装成json后push到falcon-agent
接口解释请参照:
 * http://docs.docker.com/v1.9/
 * https://github.com/fsouza/go-dockerclient
 * https://github.com/projecteru/eru-metric

使用方法
--------------------------------
1. 根据实际部署情况，配置docker daemon接口;
 * daemon: "addr": "tcp://127.0.0.1:2375" 

2. 测试： ./control build && ./control start
 * $GOPATH/bin/govendor init && $GOPATH/bin/govendor add +external && GO15VENDOREXPERIMENT=1 go build
