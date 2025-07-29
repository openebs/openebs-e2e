#!/usr/bin/env bash

set -eu

SCRIPTDIR=$(dirname "$(realpath "$0")")
EXITV_INVALID_OPTION=1

help() {
  cat <<EOF
Usage: $0 [OPTIONS]

Options:
  --testplan                 Test plan[lvm, zfs, hostpath, selfci etc.]
  --openebs_chart_version   Openebs chart version to use for install test (default: latest)
  --oci                     OCI registry to use for openesb install (default: "false")
                            On true, Use OCI registry of install test
                            On false, Use Helm registry of install test 
Examples:
  $0 --testplan lvm
EOF
}

# Parse arguments
while [ "$#" -gt 0 ]; do
  case "$1" in
    -T|--testplan)
      shift
      testplan="$1"
      ;;
    --oci)
      shift
      case $1 in
          true)
             oci="$1"
             ;;
          false)
             oci="$1"
             ;;
          *)
              echo "Unknown boolean option oci  : $1"
              exit 1
              ;;
      esac
      ;;
    --openebs_chart_version)
        shift
        openebs_chart_version="$1"
        ;;
    *)
      echo "Unknown option: $1"
      help
      exit $EXITV_INVALID_OPTION
      ;;
  esac
  shift
done

echo "Testplan: $testplan"
# Get the array elements from the command-line argument
array_str=$(python3 $SCRIPTDIR/testlists.py --testplan $testplan --install)

# Split the string into an array
array=($array_str)

# Iterate through the array
for test in "${array[@]}"; do
    echo "Test: $test"
    $SCRIPTDIR/exec-tests.sh --tests $test --local true --product openebs --replicated_engine false
done