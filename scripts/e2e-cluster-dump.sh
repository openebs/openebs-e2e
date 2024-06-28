#!/usr/bin/env bash

# This script makes the best attempt to dump stuff
# so ignore fails and keep paddling.
# set -e

SCRIPT_DIR=$(dirname "$0")

help() {
  cat <<EOF
This script generates logs for mayastor pods and cluster state.

Usage: $0 [OPTIONS]

Options:
  --destdir <path>   Location to store log files
  --clusteronly   Only generate cluster information

If --destdir is not specified the data is dumped to stdout
EOF
}

pluginPath=""

function cluster-get {
    echo "-- PODS mayastor* --------------------"
    # The CSI tests creates namespaces containing the text mayastor
    mns=$(kubectl get ns | grep -e mayastor -e datacore -e kube-system | sed -e "s/ .*//")
    for ns in $mns
    do
        kubectl -n "$ns" -o wide get pods --sort-by=.metadata.creationTimestamp
    done
    echo "-- PODS ------------------------------"
    kubectl get -o wide pods --sort-by=.metadata.creationTimestamp
    echo "-- PVCS ------------------------------"
    kubectl get pvc --sort-by=.metadata.creationTimestamp
    echo "-- PV --------------------------------"
    kubectl get pv --sort-by=.metadata.creationTimestamp
    echo "-- Volume attachments ----------------"
    kubectl get volumeattachments --sort-by=.metadata.creationTimestamp
    echo "-- Storage Classes -------------------"
    kubectl get sc --sort-by=.metadata.creationTimestamp
    echo "-- K8s Nodes -------------------------"
    kubectl get nodes -o wide --show-labels
    echo "-- Pools -----------------------------"
    kubectl -n datacore get diskpools --sort-by=.metadata.creationTimestamp
    kubectl -n mayastor get diskpools --sort-by=.metadata.creationTimestamp
    echo "-- Snapshot --------------------------"
    kubectl get volumesnapshots
    echo "-- SnapshotContents ------------------"
    kubectl get volumesnapshotcontents
    echo "-- SnapshotClasses -------------------"
    kubectl get volumesnapshotclasses
    echo "-- K8s Deployments -------------------"
    kubectl -n datacore get deployments
    kubectl -n mayastor get deployments
    echo "-- K8s Daemonsets --------------------"
    kubectl -n datacore get daemonsets
    kubectl -n mayastor get daemonsets
    if [ -n "$pluginPath" ]; then
        echo "-- MSV ---------------"
        $pluginPath -n mayastor get volumes
        echo "-- DSP ---------------"
        $pluginPath -n mayastor get pools
        echo "-- MSN ---------------"
        $pluginPath -n mayastor get nodes
        echo "-- Snapshot ----------"
        $pluginPath -n mayastor get volume-snapshots
    fi

}

function cluster-describe {
    echo "-- PODS mayastor* --------------------"
    # The CSI tests creates namespaces containing the text mayastor
    mns=$(kubectl get ns | grep -e mayastor -e datacore -e kube-system | sed -e "s/ .*//")
    for ns in $mns
    do
        kubectl -n "$ns" describe pods
    done
    echo "-- PODS ------------------------------"
    kubectl describe pods
    echo "-- PVCS ------------------------------"
    kubectl describe pvc
    echo "-- PV --------------------------------"
    kubectl describe pv
    echo "-- Volume attachements ---------------"
    kubectl describe volumeattachments
    echo "-- Storage Classes -------------------"
    kubectl describe sc
    echo "-- K8s Nodes -------------------------"
    kubectl describe nodes
    echo "-- Pools --------------------"
    kubectl -n datacore describe diskpools
    kubectl -n mayastor describe diskpools
    echo "-- Snapshot --------------------------"
    kubectl describe volumesnapshots
    echo "-- SnapshotContents ------------------"
    kubectl describe volumesnapshotcontents
    echo "-- SnapshotClasses -------------------"
    kubectl describe volumesnapshotclasses
    echo "-- K8s Deployments -------------------"
    kubectl -n datacore describe deployments
    kubectl -n mayastor describe deployments
    echo "-- K8s Daemonsets --------------------"
    kubectl -n datacore describe daemonsets
    kubectl -n mayastor describe daemonsets
    if [ -n "$pluginPath" ]; then
        echo "-- MSV ---------------"
        $pluginPath -n mayastor get volumes -o yaml
        echo "-- DSP ---------------"
        $pluginPath -n mayastor get pools -o yaml
        echo "-- MSN ---------------"
        $pluginPath -n mayastor get nodes -o yaml
        echo "-- Snapshot ----------"
        $pluginPath -n mayastor get volume-snapshots -o yaml
    fi
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

# $1 = podlogs, 0 => do not generate pod logs
# $2 = [destdir] undefined => dump to stdout,
#                   otherwise generate log files in $destdir
function getLogs {
    podlogs="$1"
    shift
    dest="$1"
    shift

    if [ -n "$dest" ];
    then
        mkdir -p "$dest"
    fi

    if [ "$podlogs" -ne 0 ]; then
        getPodLogs default "$dest"
    fi

    if [ -n "$dest" ];
    then
        cluster-get >& "$dest/cluster.get.txt"
        cluster-describe >& "$dest/cluster.describe.txt"

        echo "logfiles generated in $dest"
        echo ""

    else
        cluster-get
        cluster-describe
    fi
}

podlogs=1
destdir=

# Parse arguments
while [ "$#" -gt 0 ]; do
  case "$1" in
    -d|--destdir)
      shift
      destdir="$1"
      ;;
    -c|--clusteronly)
      podlogs=0
      ;;
    --plugin)
      shift
      pluginPath="$1"
      ;;
    *)
      echo "Unknown option: $1"
      help
      exit 1
      ;;
  esac
  shift
done

getLogs "$podlogs" "$destdir"
getPodLogs "e2e-agent" "$dest"

function getSystemCmdOutputs {
        dest="$1"
        shift

        if [ -n "$dest" ];
        then
                mkdir -p "$dest"
        fi

        kubectl get nodes -owide >& "$dest/node-list-with-ip"

        nodes=$(kubectl get nodes -o jsonpath='{ $.items[*].status.addresses[?(@.type=="InternalIP")].address }')
        for node in $nodes
        do
                timeout -s 9 30 curl --connect-timeout 5 -XPOST "http://$node:10012/nvmelist" >& "$dest/$node-nvme-list"
                timeout -s 9 30 curl --connect-timeout 5 -XPOST "http://$node:10012/findmnt" >& "$dest/$node-findmnt"
                timeout -s 9 30 curl --connect-timeout 5 -XPOST "http://$node:10012/lsblk" >& "$dest/$node-lsblk"
                timeout -s 9 120 curl --connect-timeout 5 -XPOST "http://$node:10012/syslog" >& "$dest/$node-syslog"
                timeout -s 9 120 curl --connect-timeout 5 -XPOST "http://$node:10012/dmesg" >& "$dest/$node-dmesg"
        done
}

if [ "$podlogs" -ne 0 ]; then
    getSystemCmdOutputs "$destdir"
fi

function getEvents {
    dest="$1"
    event_collector_dir="${SCRIPT_DIR}/../src/tools/e2e_event_collector/"
    if [ -d "${event_collector_dir}" ]; then
    	outfile_json=$(realpath "$dest/events.json")
    	outfile_yaml=$(realpath "$dest/events.yaml")
        pushd ${event_collector_dir} > /dev/null \
		&& go run e2e_event_collector.go -o json -f "${outfile_json}" \
		&& go run e2e_event_collector.go -o yaml -f "${outfile_yaml}"; \
		popd > /dev/null
    fi
}

getEvents "$destdir"
