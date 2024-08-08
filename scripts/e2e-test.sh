#!/usr/bin/env bash

set -eu

SCRIPTDIR=$(dirname "$(realpath "$0")")
$SCRIPTDIR/exec_tests.py "$@"
