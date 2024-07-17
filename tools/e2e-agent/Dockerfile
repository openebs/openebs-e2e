From ubuntu:20.04
ARG GO_VERSION
ARG E2EA_VERSION
# Set environment variable to prevent interactive prompts
ENV DEBIAN_FRONTEND=noninteractive
# Set the timezone environment variable to UTC
ENV TZ=UTC
COPY  ./ /
# Parted is used to create and manipulate disk partition.
# Udev provides a dynamic way of setting up device.
# It ensures that devices are configured as soon as they are plugged in and discovered.
# It propagates information about a processed device.
RUN apt-get update; apt-get install net-tools iptables wget parted udev nvme-cli build-essential gettext gettext-base \
        libinih-dev uuid-dev liburcu-dev  libblkid-dev btrfs-progs lvm2 -y;
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz; \
        tar -C /usr/local/ -xzf go${GO_VERSION}.linux-amd64.tar.gz; \
        rm -rf go${GO_VERSION}.linux-amd64.tar.gz; \
        export PATH=$PATH:/usr/local/go/bin; go build -ldflags="-X 'main.Version="${E2EA_VERSION}"'"; rm -rf /usr/local/go; \
        mv e2e-agent /usr/local/bin; wget https://mirrors.edge.kernel.org/pub/linux/utils/fs/xfs/xfsprogs/xfsprogs-6.6.0.tar.xz; \
        tar xvf xfsprogs-6.6.0.tar.xz; cd xfsprogs-6.6.0; ./configure; make; make install
CMD ["e2e-agent"]
