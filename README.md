docker监控脚本
================================
[![Build Status](https://travis-ci.org/mesos-utility/docker-metrics.png?branch=dev)](https://travis-ci.org/mesos-utility/docker-metrics)

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
 * 注意点: 推荐使用go1.5版本及以上编译源代码,因为使用了GO15VENDOREXPERIMENT及govendor做lib库的版本管理;

采集的指标
--------------------------
| Counters | Notes|
|-----|------|
|cpu.system.rate|内核态使用的CPU百分比|
|cpu.usage.rate|cpu使用情况百分比|
|cpu.user.rate|用户态使用的CPU百分比|
|mem.usage|内存使用百分比|
|mem.rss|内存使用原值|
|mem.max_usage|内存总量|
|disk.io.read_bytes|磁盘io读字节数|
|disk.io.write_bytes|磁盘io写字节数|
|ifname.inbits.rate|网络io流入bits数|
|ifname.inpackets.rate|网络io流入包数|
|ifname.inerrs.rate|网络io流入出错数|
|ifname.indrop.rate|网络io流入丢弃数|
|ifname.outbits.rate|网络io流出bits数|
|ifname.outpackets.rate|网络io流出包数|
|ifname.outerrs.rate|网络io流出出错数|
|ifname.outdrop.rate|网络io流出丢弃数|
