# VERSION 0.1
FROM alpine:edge

MAINTAINER openebs openebs@openebs.io

RUN apk --no-cache add \
    	fio=3.37-r0 gcc musl-dev

ADD e2e_fio.c /
ADD liveness.c /
ADD liveness.h /
ADD e2e_fio_version.h /
ADD ./archive_wksp.sh /

RUN gcc -Wall -o /e2e_fio e2e_fio.c liveness.c
RUN apk --no-cache del \
        gcc musl-dev

ENTRYPOINT ["/e2e_fio"]
