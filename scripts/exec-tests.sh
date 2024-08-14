#!/usr/bin/env bash

set -eu

SCRIPTDIR=$(dirname "$(realpath "$0")")
E2EROOT=$(realpath "$SCRIPTDIR/..")
TESTDIR=$(realpath "$SCRIPTDIR/../src")
ARTIFACTSDIR="${E2E_ARTIFACTSDIR:-$(realpath $SCRIPTDIR/../artifacts)}"

#exit values
EXITV_OK=0
EXITV_INVALID_OPTION=1
# EXITV_MISSING_OPTION=2
EXITV_FAILED=4
#EXITV_FILE_MISMATCH=5
#EXITV_CRD_GO_GEN=6
EXITV_VERSION_MISMATCH=7
# EXITV_MISSING_KUBECTL_PLUGIN=8
EXITV_FILE_MISSING=9
# EXITV_CONFIG_FILE_CREATE=10
EXITV_FAILED_CLUSTER_OK=255
# EXITV_UNSET_ENVVARS=254

# Global state variables
#  test configuration state variables
session="$(date +%Y%m%d-%H%M%S-)$(uuidgen -r)"

#  script state variables
tests=""
on_fail="stop"
uninstall_cleanup="n"
generate_logs=0
logsdir="$ARTIFACTSDIR/logs"
reportsdir="$ARTIFACTSDIR/reports"
coveragedir="$ARTIFACTSDIR/coverage/data"
policy_cleanup_before="${e2e_policy_cleanup_before:-false}"
test_list=""
product=
local=
replicated_engine=

help() {
  cat <<EOF
Usage: $0 [OPTIONS]

Options:
  --tests <list of tests>   Lists of tests to run, delimited by spaces (default: "$tests")
                            Note: the last test should be uninstall (if it is to be run)
  --reportsdir <path>       Path to use for junit xml test reports (default: repo root)
  --logs                    Generate logs and cluster state dump at the end of successful test run,
                            prior to uninstall.
  --logsdir <path>          Location to generate logs (default: emit to stdout).
  --onfail <stop|uninstall|reinstall|restart>
                            On fail, stop immediately,uninstall, reinstall and continue or restart and continue default($on_fail)
                            Behaviour for "uninstall" only differs if uninstall is in the list of tests (the default).
                            If set to "reinstall" on failure, all resources are cleaned up and mayastor is re-installed.
                            If set to "restart" on failure, all resources are cleaned up and mayastor pods are restarted by deleting.
  --uninstall_cleanup <y|n> On uninstall cleanup for reusable cluster. default($uninstall_cleanup)
  --config                  config name or configuration file default($config_file)
  --platform_config         test platform configuration file default($platform_config_file)
  --tag <name>              Docker image tag of mayastor images (default "$tag")
                            install files are retrieved from the CI registry using the appropriately
                            tagged docker image :- mayadata/install-images
  --mayastor                path to the mayastor source tree to use for testing.
                            If this is specified the install test uses the install yaml files from this tree
                            instead of the tagged image.
  --session                 session name, adds a subdirectory with session name to artifacts, logs and reports
                            directories to facilitate concurrent execution of test runs (default timestamp-uuid)
  --version                 Mayastor version, 0 => MOAC, > 1 => restful control plane
  --local                   This option will only exercised for local testing not in Jenkins pipeline.
                            On true, Creates namespace(mayastor/datacore) on cluster in case of install test
                            On false, namespace will not be created on cluster in case of install test
Examples:
  $0 --registry 127.0.0.1:5000 --tag a80ce0c --product openebs
EOF
}

# Parse arguments
while [ "$#" -gt 0 ]; do
  case "$1" in
    -T|--tests)
      shift
      test_list="$1"
      ;;
    -R|--reportsdir)
      shift
      reportsdir="$1"
      ;;
    -h|--help)
      help
      exit $EXITV_OK
      ;;
    --logs)
      generate_logs=1
      ;;
    --logsdir)
      shift
      logsdir="$1"
      if [[ "${logsdir:0:1}" == '.' ]]; then
          logsdir="$PWD/$logsdir"
      fi
      ;;
    --onfail)
        shift
        case $1 in
            uninstall)
                on_fail=$1
                ;;
            stop)
                on_fail=$1
                ;;
            reinstall|continue)
                on_fail="reinstall"
                policy_cleanup_before='true'
                ;;
            restart)
                on_fail=$1
                policy_cleanup_before='true'
                ;;
            *)
                echo "invalid option for --onfail"
                help
                exit $EXITV_INVALID_OPTION
        esac
      ;;
    --uninstall_cleanup)
        shift
        case $1 in
            y|n)
                uninstall_cleanup=$1
                ;;
            *)
                echo "invalid option for --uninstall_cleanup"
                help
                exit $EXITV_INVALID_OPTION
        esac
      ;;
    --config)
        shift
        config_file="$1"
        ;;
    --platform_config)
        shift
        platform_config_file="$1"
        ;;
    --product)
      shift
      case $1 in
          openebs)
             product="$1"
             ;;
          *)
              echo "Unknown product: $1"
              exit 1
              ;;
      esac
      ;;
    --session)
        shift
        session="$1"
        ;;
    --version)
        shift
            case "$1" in
                0|1)
                    mayastor_version=$1
                    ;;
                *)
                    echo "Unknown control plane: $1"
                    help
                    exit $EXITV_INVALID_OPTION
                    ;;
            esac
        ;;
    --replicated_engine)
      shift
      case $1 in
          true)
             replicated_engine="$1"
             ;;
          false)
             replicated_engine="$1"
             ;;
          *)
              echo "Unknown boolean option to replicated_engine : $1"
              exit 1
              ;;
      esac
      ;;
    --local)
      shift
      case $1 in
          true)
             local="$1"
             ;;
          false)
             local="$1"
             ;;
          *)
              echo "Unknown boolean option to local(Only used by Developer to test locally) : $1"
              exit 1
              ;;
      esac
      ;;
    --SPE)
        # short form for set product envvars
        :
        ;;
    *)
      echo "Unknown option: $1"
      help
      exit $EXITV_INVALID_OPTION
      ;;
  esac
  shift
done

if [ -z "$product" ]; then
    echo "defaulting product to openebs"
    product="openebs"
fi

export e2e_product="${product}"
export replicatedEngine="${replicated_engine}"

if [ -z "$session" ]; then
    sessiondir="$ARTIFACTSDIR"
else
    sessiondir="$ARTIFACTSDIR/sessions/$session"
    logsdir="$logsdir/$session"
    reportsdir="$reportsdir/$session"
    coveragedir="$coveragedir/$session"
fi

export e2e_session_dir=$sessiondir
export e2e_port_forwarding_enabled="True"


export e2e_root_dir="$E2EROOT"
export openebs_e2e_root_dir=$(realpath "$SCRIPTDIR/..")

tests="$test_list"

export e2e_reports_dir="$reportsdir"

if [ "$uninstall_cleanup" == 'n' ] ; then
    export e2e_uninstall_cleanup=0
else
    export e2e_uninstall_cleanup=1
fi

mkdir -p "$sessiondir"
mkdir -p "$reportsdir"
mkdir -p "$logsdir"

kubectl get nodes -o yaml > "$reportsdir/k8s_nodes.yaml"

test_failed=0

# run a go program in the src tree
# arguments <directory-relative-to-src> <main-go-file> [arguments...]
function runGo {
    pushd "src"
    go mod tidy
    if [ -z "$1" ] || [ ! -d "$1" ]; then
        echo "Unable to locate directory  $PWD/\"$1\""
        popd
        return 1
    fi
    cd "$1"
    shift
    main_go_file="$1"
    shift
    if ! go run "$main_go_file" "$@" ; then
        popd
        return 1
    fi

    popd
    return 0
}

# Run go test in directory specified as $1 (relative path)
# $2 is the fully qualified testname
# $3 is the name of the junit report file
# maximum test runtime is 120 minutes
function runGoTest {
    pushd "$1"
    if [ -f "go.mod" ]; then
        go mod tidy
    fi
    echo "Running go test in $PWD/\"$2\""
    if [ -z "$2" ] || [ ! -d "$2" ]; then
        echo "Unable to locate test directory  $PWD/\"$2\""
        popd
        echo "Finished go test in $PWD/\"$2\" result:notfound"
        return 1
    fi

    if [ -z "$3" ] ; then
        echo "name of the report file not specified"
        popd
        echo "Finished go test in $PWD/\"$2\" result:invalidarguments"
        return 1
    fi

    export e2etestlogdir="$logsdir/$2"
    cd "$2"

    # timeout test run after 1 hour
    runtimeout=60m
    echo "test run timeout=$runtimeout"
    if ! go test -v . -ginkgo.v -ginkgo.timeout=${runtimeout} -ginkgo.junit-report="${reportsdir}/e2e.$3.junit.xml" -timeout ${runtimeout}; then
        echo "Finished go test in $PWD/\"$2\" result:failed"
        popd
        return 1
    fi
    echo "Finished go test in $PWD/\"$2\" result:passed"
    popd
    return 0
}

function emitLogs {
    if [ -z "$1" ]; then
        logPath="$logsdir"
    else
        logPath="$logsdir/$1"
    fi
    mkdir -p "$logPath"
    if ! "$SCRIPTDIR/e2e-cluster-dump.sh" --destdir "$logPath" ; then
        # ignore failures in the dump script
        :
    fi
    unset logPath
}


# Check if $2 is in $1
contains() {
    [[ $1 =~ (^|[[:space:]])$2($|[[:space:]]) ]] && return 0  || return 1
}

export e2e_policy_cleanup_before="$policy_cleanup_before"
export e2e_product_config_file="$E2EROOT/configurations/product/mayastor_config.yaml"


#preprocess tests so that command line can use commas as delimiters
tests=${tests//,/ }

echo "Environment:"
echo "    e2e_session_dir=$e2e_session_dir"
echo "    e2e_root_dir=$e2e_root_dir"
echo "    openebs_e2e_root_dir=$openebs_e2e_root_dir"
echo "    e2e_product=$e2e_product"
echo "    e2e_reports_dir=$e2e_reports_dir"
echo "    e2e_uninstall_cleanup=$e2e_uninstall_cleanup"
echo "    e2e_policy_cleanup_before=$e2e_policy_cleanup_before"
echo ""
echo "Script control settings:"
echo "    on_fail=$on_fail"
echo "    uninstall_cleanup=$uninstall_cleanup"
echo "    generate_logs=$generate_logs"
echo "    logsdir=$logsdir"
echo ""
echo "list of tests: $tests"


for testname in $tests; do
  # Defer uninstall till after other tests have been run.
  if [ "$testname" != "uninstall" ]; then
    if [ -d "$TESTDIR/tests/$testname" ]; then
        testrootdir=$TESTDIR
        fqtestname="tests/$testname"
    elif [ -d "$TESTDIR/$testname" ]; then
        testrootdir=$TESTDIR
        fqtestname="$testname"
    else
        echo "test directory $testname not found under $TESTDIR/tests/ or $TESTDIR"
        exit $EXITV_FILE_MISSING
    fi

    if ! runGoTest "$testrootdir" "$fqtestname" "$testname"; then
        echo "Test \"$testname\" FAILED!"
        test_failed=1
        emitLogs "$testname"
    fi
  fi
done


if [ "$generate_logs" -ne 0 ]; then
    emitLogs ""
fi

if [ "$test_failed" -ne 0 ] && [ "$on_fail" == "stop" ]; then
    echo "At least one test FAILED!"
    exit $EXITV_FAILED
fi

# Always run uninstall test if specified
if contains "$tests" "uninstall"; then
    if ! runGoTest $TESTDIR "tests/uninstall" "uninstall"; then
        echo "Test \"uninstall\" FAILED!"
        test_failed=1
        emitLogs "uninstall"
    else
        if [ "$test_failed" -ne 0 ]; then
            # Tests failed, but uninstall was successful
            # So cluster is reusable
            echo "At least one test FAILED! Cluster is usable."
            exit $EXITV_FAILED_CLUSTER_OK
        fi
    fi
fi

if [ "$test_failed" -ne 0 ] ; then
    echo "At least one test FAILED!"
    exit $EXITV_FAILED
fi

echo "All tests have PASSED!"
exit $EXITV_OK
