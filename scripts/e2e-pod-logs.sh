#!/usr/bin/env bash

# This script extracts logs from a completed pod
# and copies the log file to directory specified

help() {
  cat <<EOF
This script generates logs for a completed pod
and copies the log file to directory specified

Usage: $0 [OPTIONS]

Options:
  --destdir <path>   Location to store log files

If --destdir is not specified the data is dumped to stdout
EOF
}




# args destdir previous namespace podname containername
# destdir == "" -> stdout
# previous == "" -> current logs else -> previous logs
function kubectlEmitLogs {
    destdir=$1
    previous=$2
    ns=$3
    podname=$4
    containername=$5

    if [ -z "$ns" ] || [ -z "$podname" ] || [ -z "$containername" ]  ; then
        echo "ERROR calling kubectlEmitLogs"
        return
    fi

    if [ -z "$previous" ] ; then
        logfile="$destdir/$podname.$containername.log"
        msg="# $podname $containername ----------------------------"
        cmd="kubectl -n $ns logs $podname $containername"
    else
        logfile="$destdir/$podname.$containername.previous.log"
        msg="# $podname $containername previous -------------------"
        cmd="kubectl -n $ns logs -p $podname $containername"
    fi

    if [ -n "$destdir" ]; then
        $cmd >& "$logfile"
    else
        echo "$msg"
        $cmd
    fi
}

# args = namespace destdir podname
# if $destdir != "" then log files are generate in $destdir
#   with the name of the pod and container.
function emitPodContainerLogs {
    ns=$1
    destdir=$2
    podname=$3

    if [ -z "$podname" ] || [ -z "$ns" ]; then
        echo "ERROR calling emitPodContainerLogs"
        return
    fi

    restarts=$(kubectl -n "$ns" get pods "$podname" | grep -v NAME | awk '{print $4}')
    containernames=$(kubectl -n "$ns" get pods "$pod" -o jsonpath="{.spec.containers[*].name}")
    for containername in $containernames
    do
        if [ "$restarts" != "0" ]; then
            kubectlEmitLogs "$destdir" "1" "$ns" "$podname" "$containername"
        fi

        kubectlEmitLogs "$destdir" "" "$ns" "$podname" "$containername"
    done
}

# $1 = namespace
function getPodLogs {
    ns=$1
    dest=$2
    if [ -n "$dest" ];
    then
        mkdir -p "$dest"
    fi

    if [ -z "$ns" ]; then
        echo "ERROR calling getPodLogs"
        return
    fi

    pods=$(kubectl -n "$ns" get pods | grep -v NAME | sed -e 's/ .*//')
    for pod in $pods
    do
        emitPodContainerLogs "$ns" "$2" "$pod"
    done
}


destdir=

# Parse arguments
while [ "$#" -gt 0 ]; do
  case "$1" in
    -d|--destdir)
      shift
      destdir="$1"
      ;;
    *)
      echo "Unknown option: $1"
      help
      exit 1
      ;;
  esac
  shift
done

getPodLogs default "$destdir"
