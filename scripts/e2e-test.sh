#!/usr/bin/env bash

set -eu

SCRIPTDIR=$(dirname "$(realpath "$0")")
EXITV_INVALID_OPTION=1

help() {
  cat <<EOF
Usage: $0 [OPTIONS]

Options:
  --testplan                 Test plan[lvm, zfs, hostpath, selfci etc.]
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