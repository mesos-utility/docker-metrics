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
 * Attention: recommand go1.5 and above, use GO15VENDOREXPERIMENT and govendor for lib vcs.

采集的指标
--------------------------
| Counters | Notes|
|-----|------|
|docker.cpu.system | 内核态使用的CPU百分比|
|docker.cpu.usage | cpu使用情况百分比|
|docker.cpu.user | 用户态使用的CPU百分比|
|docker.mem.usage | 内存使用百分比|
|docker.mem.rss | 内存使用原值|
|docker.mem.max_usage | 内存总量|
|docker.disk.io.read_bytes | 磁盘io读字节数|
|docker.disk.io.write_bytes | 磁盘io写字节数|
|docker.ifname.inbits | 网络io流入bits数|
|docker.ifname.inpackets | 网络io流入包数|
|docker.ifname.inerrs | 网络io流入出错数|
|docker.ifname.indrop | 网络io流入丢弃数|
|docker.ifname.outbits | 网络io流出bits数|
|docker.ifname.outpackets | 网络io流出包数|
|docker.ifname.outerrs | 网络io流出出错数|
|docker.ifname.outdrop | 网络io流出丢弃数|
