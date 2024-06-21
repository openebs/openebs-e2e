FROM ubuntu:22.04

RUN apt update && \
    apt install -y wget && \
    apt install -y tar && \
    apt install -y default-jre

RUN echo "JAVA_HOME=$(java -XshowSettings:properties -version 2>&1 | sed '/^[[:space:]]*java\.home/!d;s/^[[:space:]]*java\.home[[:space:]]*=[[:space:]]*//')" >> /etc/environment
RUN tail /etc/environment
CMD ["/bin/sh" "-c" "source /etc/environment"]
WORKDIR /app
RUN wget https://github.com/brianfrankcooper/YCSB/releases/download/0.17.0/ycsb-mongodb-binding-0.17.0.tar.gz
RUN tar -xzf ycsb-mongodb-binding-0.17.0.tar.gz && rm ycsb-mongodb-binding-0.17.0.tar.gz && mv ycsb-* ycsb
ENV YCSB_HOME=/app/ycsb

CMD ["/bin/sh"]