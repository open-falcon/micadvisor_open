FROM progrium/busybox


MAINTAINER Meng Zhuo <mengzhuo@xiaomi.com>

ENV Interval 60

ADD cadvisor /home/work/uploadCadviosrData/cadvisor

ADD uploadCadvisorData /home/work/uploadCadviosrData/uploadCadvisorData

ADD run /home/work/uploadCadviosrData/run

RUN mkdir -p /home/work/uploadCadviosrData/log

EXPOSE 8080


ENTRYPOINT ["/home/work/uploadCadviosrData/run"]

