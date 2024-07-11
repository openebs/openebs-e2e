FROM ubuntu:20.04

# Set environment variable to prevent interactive prompts
ENV DEBIAN_FRONTEND=noninteractive
# Set the timezone environment variable to UTC
ENV TZ=UTC

RUN apt-get update && \
    apt-get install -y wget && \
    apt-get install -y tar && \
    apt-get install -y default-jre

RUN echo "JAVA_HOME=$(java -XshowSettings:properties -version 2>&1 | sed '/^[[:space:]]*java\.home/!d;s/^[[:space:]]*java\.home[[:space:]]*=[[:space:]]*//')" >> /etc/environment
RUN tail /etc/environment
CMD ["/bin/sh" "-c" "source /etc/environment"]
WORKDIR /app
RUN wget https://github.com/brianfrankcooper/YCSB/releases/download/0.17.0/ycsb-mongodb-binding-0.17.0.tar.gz
RUN tar -xzf ycsb-mongodb-binding-0.17.0.tar.gz && rm ycsb-mongodb-binding-0.17.0.tar.gz && mv ycsb-* ycsb
ENV YCSB_HOME=/app/ycsb

CMD ["/bin/sh"]