#!/usr/bin/env sh

#dest="/mnt/host/tmp"
dest="/tmp"
if [ -d ${dest} ]; then
    echo "archiving ./workspace to ${dest}/e2e_fio-$1.tgz"
    tar cvzf ${dest}/e2e_fio-$1.tgz ./workspace
    echo "{ base64 e2e_fio-$1.tgz"
    echo "------------------"
    base64 ${dest}/e2e_fio-$1.tgz
    echo "------------------"
    echo "} base64 e2e_fio-$1.tgz"
else
    echo "path ${dest} does not exist"
fi
