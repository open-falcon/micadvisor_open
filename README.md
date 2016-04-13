Docker container monitor plugin for Open-Falcon  Micadvisor-Open
------------------------------------
描述
------------------
micadvisor-open是基于open-falcon的docker容器资源监控插件，监控容器的cpu、内存、diskio以及网络io等，数据采集后上报到open-falcon

biuld方法：
-----------------
./build

启动方法：
-----------------
```bash
    docker run \
    --volume=/:/rootfs:ro \
    --volume=/var/run:/var/run:rw \
    --volume=/sys:/sys:ro \
    --volume=/home/work/log/cadvisor/:/home/work/uploadCadviosrData/log \
    --volume=/var/lib/docker/:/var/lib/docker:ro \
    --volume=/home/docker/containers:/home/docker/containers:ro \
    --publish=18080:18080 \
    --env Interval=60 \
    --detach=true \
    --name=micadvisor \
    --net=host \
    --restart=always \
    micadvisor:latest
```

注：
```
    容器的endpoint获取有两种方式:
    1,对被监控的容器传入env，EndPoint=myendpoint 方式；
    2,通过获取被监控的容器的hosts文件的第一个ip地址
    --volume=/sys:/sys:ro 此volume中包含docker容器监控所需要的重要内容，如/sys/fs/cgroup下的相关内容
    --volume=/home/work/log/cadvisor/:/home/work/uploadCadviosrData/log \ 为日志内容路径
    --env Interval=60 表示提取数据的间隔时间
```

采集的指标
--------------------------
| Counters | Notes|
|-----|------|
|cpu.busy|cpu使用情况百分比|
|cpu.user|用户态使用的CPU百分比|
|cpu.system|内核态使用的CPU百分比|
|cpu.core.busy|每个cpu的使用情况|
|mem.memused.percent|内存使用百分比|
|mem.memused|内存使用原值|
|mem.memtotal|内存总量|
|mem.memused.hot|内存热使用情况|
|disk.io.read_bytes|磁盘io读字节数|
|disk.io.write_bytes|磁盘io写字节数|
|net.if.in.bytes|网络io流入字节数|
|net.if.in.packets|网络io流入包数|
|net.if.in.errors|网络io流入出错数|
|net.if.in.dropped|网络io流入丢弃数|
|net.if.out.bytes|网络io流出字节数|
|net.if.out.packets|网络io流出包数|
|net.if.out.errors|网络io流出出错数|
|net.if.out.dropped|网络io流出丢弃数|

Contributors
------------------------------------------
- mengzhuo: QQ:296142139; MAIL:mengzhuo@xiaomi.com 
