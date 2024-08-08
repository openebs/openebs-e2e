#!/usr/bin/env bash

set -eu

SCRIPTDIR=$(dirname "$(realpath "$0")")
E2EROOT=$(realpath "$SCRIPTDIR/..")
TESTDIR=$(realpath "$SCRIPTDIR/../src")
ARTIFACTSDIR="${E2E_ARTIFACTSDIR:-$(realpath $SCRIPTDIR/../artifacts)}"


#exit values
EXITV_OK=0
EXITV_INVALID_OPTION=1
EXITV_MISSING_OPTION=2
EXITV_FAILED=4
#EXITV_FILE_MISMATCH=5
#EXITV_CRD_GO_GEN=6
EXITV_VERSION_MISMATCH=7
EXITV_MISSING_KUBECTL_PLUGIN=8
EXITV_FILE_MISSING=9
EXITV_CONFIG_FILE_CREATE=10
EXITV_FAILED_CLUSTER_OK=255
EXITV_UNSET_ENVVARS=254




platform_config_file="hetzner.yaml"
config_file="hcloudci_config.yaml"


session="$(date +%Y%m%d-%H%M%S-)$(uuidgen -r)"
tests=""
on_fail="stop"
uninstall_cleanup="n"
logsdir="$ARTIFACTSDIR/logs"
reportsdir="$ARTIFACTSDIR/reports"

policy_cleanup_before="${e2e_policy_cleanup_before:-false}"
test_list=""

product=
local=
maas_api_token=
maas_endpoint=


help() {
  cat <<EOF
Usage: $0 [OPTIONS]
Pre-requisites:
  environment variables:
    CI_REGISTRY must be set the docker CI registry for the product
    GCS_BUCKET must be set to the GCS bucket  for the product

Options:
  --build_number <number>   Build number, for use when sending Loki markers
  --loki_run_id <Loki run id>  ID string, for use when sending Loki markers
  --loki_test_label <Loki custom test label> Test label value, for use when sending Loki markers
  --registry <host[:port]>  Registry to pull the mayastor images from.
                            'dockerhub' means use DockerHub
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
  --grpc_code_gen <true|false>
                            On true, grpc server and clinet code will be generated
                            On false, grpc server and clinet code will not be generated
  --crd_code_gen <true|false>
                            On true, custom resource clinet code will be generated
                            On false, custom resource clinet code will not be generated
  --product                 Product key [mayastor, openebspro]
  --install_crds <true|false>
                            On true, product specific crds gets applied
                            On false, product specific crds will not be applied
  --licensegen              Path to the license-generator script
  --local                   This option will only exercised for local testing not in Jenkins pipeline.
                            On true, Creates namespace(mayastor/datacore) on cluster in case of install test
                            On false, namespace will not be created on cluster in case of install test
  --secret_config_file      This option will only exercised for local testing not in Jenkins pipeline with
                            local option. Provide docker credential yaml file path to create kubernetes secret
  --use_rest_api            Use rest api control plane instead of kubectl plugin
  --gen_rest_api            Generate and use rest api control plane instead of kubectl plugin
  --maas_api_token          This option will only exercised while running tests on maas platform.
                            Provide maas oauth api token which will used for injecting fault to nodes.
  --maas_endpoint           Provide maas endpoint with port. Ex: 127.0.0.0:8080
  --grpc_version            gRPC version
  --coverage                enable coverage
  --license_server          License server if prodct is openebspro
  --jenkins_username        Jenkins username used for authentication in case product is openebspro
                            This will be used to generate expiring license on demand
  --jenkins_api_token       Jenkins username used for authentication in case product is openebspro
                            This will be used to generate expiring license on demand
  --jenkins_job_buildtoken  Jenkins licnse expiry job build token to trigger job remotely
Examples:
  $0 --registry 127.0.0.1:5000 --tag a80ce0c --product openebspro
EOF
}

# Parse arguments
while [ "$#" -gt 0 ]; do
  case "$1" in
    -m|--mayastor)
      shift
      mayastor_root_dir=$1
      ;;
    -d|--device)
      shift
#      backward compatibility ignore -d|--device
#      device=$1
      ;;
    -r|--registry)
      shift
      if [[ "$1" == 'dockerhub' ]]; then
          registry=''
      else
          registry=$1
      fi
      ;;
    -t|--tag)
      shift
      tag=$1
      ;;
    -g|--grpc_code_gen)
      shift
      grpc_code_gen="$1"
      ;;
    -c|--crd_code_gen)
      shift
      crd_code_gen="$1"
      ;;
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
    --build_number) # TODO remove this option
      shift
      loki_run_id="$1"
      ;;
    --loki_run_id)
      shift
      loki_run_id="$1"
      ;;
    --loki_test_label)
      shift
      loki_test_label="$1"
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
    # --licensegen)
    #     shift
    #     licensegen=$(realpath "$1")
	# ;;
    --platform_config)
        shift
        platform_config_file="$1"
        ;;
    --product)
      shift
      case $1 in
          mayastor)
             product="$1"
             ;;
          openebspro)
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
    --ssh_identity)
        shift
        ssh_identity="$1"
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
    --install_crds)
      shift
      case $1 in
          true)
             install_crds="$1"
             ;;
          false)
             install_crds="$1"
             ;;
          *)
              echo "Unknown boolean option to install crds : $1"
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
    --secret_config_file)
        shift
        secret_config_file="$1"
        ;;
    --use_rest_api)
        shift
        e2e_control_plane_rest_api="true"
        ;;
    --gen_rest_api)
        shift
        e2e_control_plane_rest_api="true"
        gen_control_plane_rest_api=1
        ;;
    --maas_api_token)
        shift
        maas_api_token="$1"
        ;;
    --maas_endpoint)
        shift
        maas_endpoint="$1"
        ;;
    --grpc_version)
        shift
        e2e_grpc_version="$1"
        ;;
    --license_server)
        shift
        license_server="$1"
        ;;
    --jenkins_username)
        shift
        jenkins_username="$1"
        ;;
    --jenkins_api_token)
        shift
        jenkins_api_token="$1"
        ;;
    --jenkins_job_buildtoken)
        shift
        jenkins_job_buildtoken="$1"
        ;;
    --reactor_freeze_detect)
        reactor_freeze_detect=1
        ;;
    --coverage)
        coverage="true"
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

# export loki_run_id="$loki_run_id" # can be empty string
# export loki_test_label="$loki_test_label"

if [ -z "$product" ]; then
    echo "defaulting product to mayastor"
    product="mayastor"
fi

export e2e_product="${product}"

if [ -z "$session" ]; then
    sessiondir="$ARTIFACTSDIR"
else
    sessiondir="$ARTIFACTSDIR/sessions/$session"
    logsdir="$logsdir/$session"
    reportsdir="$reportsdir/$session"
    # coveragedir="$coveragedir/$session"
fi


export e2e_session_dir=$sessiondir
export e2e_maas_api_token=$maas_api_token
export e2e_maas_endpoint=$maas_endpoint
export e2e_grpc_version=$e2e_grpc_version
export e2e_port_forwarding_enabled="True"

if [ -n "$tag" ]; then
  export e2e_image_tag="$tag"
fi

# export e2e_docker_registry="$registry" # can be empty string
# export e2e_root_dir="$E2EROOT"
export openebs_e2e_root_dir=$(realpath "$SCRIPTDIR/../openebs-e2e")

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
        # for mayastor-e2e/src go mod file is here
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
    if [ -f "go.mod" ]; then
        # for 3rd party tests go mod file is here
        go mod tidy
    fi
    # timeout test run after 3 hours
    runtimeout=$($SCRIPTDIR/runtimeouts.py $2)
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



# export e2e_mayastor_version=$mayastor_version
export e2e_config_file="$config_file"
export e2e_platform_config_file="$platform_config_file"
export e2e_policy_cleanup_before="$policy_cleanup_before"



#preprocess tests so that command line can use commas as delimiters
tests=${tests//,/ }



# install and support bundle require product configuration which is
# embedded in the e2e_config package, render it readable as yaml
# by external entities.
e2e_product_config_yaml="$e2e_session_dir/product_config.yaml"
if ! runGo "tools/product-config" product-config.go -out "$e2e_product_config_yaml" ; then
    echo "failed to create $e2e_product_config_yaml"
    exit $EXITV_CONFIG_FILE_CREATE
fi



echo "Environment:"
echo "    e2e_session_dir=$e2e_session_dir"
echo "    e2e_mayastor_root_dir=$e2e_mayastor_root_dir"
echo "    loki_run_id=$loki_run_id"
echo "    loki_test_label=$loki_test_label"
echo "    e2e_root_dir=$e2e_root_dir"
echo "    openebs_e2e_root_dir=$openebs_e2e_root_dir"
echo "    e2e_product=$e2e_product"
echo "    e2e_image_tag=$e2e_image_tag"
echo "    e2e_docker_registry=$e2e_docker_registry"
echo "    e2e_reports_dir=$e2e_reports_dir"
echo "    e2e_uninstall_cleanup=$e2e_uninstall_cleanup"
echo "    e2e_config_file=$e2e_config_file"
echo "    e2e_platform_config_file=$e2e_platform_config_file"
echo "    e2e_policy_cleanup_before=$e2e_policy_cleanup_before"
echo "    e2e_mayastor_version=$e2e_mayastor_version"
echo "    e2e_install_crds=$e2e_install_crds"
echo "    licensegen=${licensegen}"
echo "    no_pool_install=${no_pool_install}"
echo "    install_license=${install_license}"
echo "    install_loki=${install_loki}"
echo "    e2e_grpc_version=${e2e_grpc_version}"
echo ""
echo "Script control settings:"
echo "    on_fail=$on_fail"
echo "    uninstall_cleanup=$uninstall_cleanup"
echo "    generate_logs=$generate_logs"
echo "    logsdir=$logsdir"
echo ""
echo "list of tests: $tests"

if [ "$coverage" == "true" ]; then
    if contains "$tests" "install" ; then
        if ! "$SCRIPTDIR/remote-coverage-files.py" --clear --identity_file "$ssh_identity" ; then
            echo "***************************** failed to clear coverage files"
        fi
    fi
fi

for testname in $tests; do
  # defer uninstall till after other tests have been run.
  if [ "$testname" != "uninstall" ] ;  then
      if [ -d "$TESTDIR/tests/$testname" ]; then
          testrootdir=$TESTDIR
          fqtestname="tests/$testname"
      else
          if [ -d "$TESTDIR3RDPARTY/$testname" ]; then
              testrootdir=$TESTDIR3RDPARTY
              fqtestname="$testname"
          else
              echo "test directory $testname not found under $TESTDIR/tests/ or $TESTDIR3RDPARTY"
              exit $EXITV_FILE_MISSING
          fi
      fi
      if [ "$testname" = "install" ] ;  then
        # Creating namespace only before install test
        if [ "$local" == "true" ]; then
            if ! "$SCRIPTDIR/e2e-helper.py"  --product_config_file "$e2e_product_config_yaml" --secret_config_file "$secret_config_file" ; then
                echo "Failed to create namespace and secret(in case of openebspro)"
            fi
        fi
      fi
      if ! runGoTest "$testrootdir" "$fqtestname" "$testname" ; then
          echo "Test \"$testname\" FAILED!"
          test_failed=1
#          emitLogs "$testname"
#          echo "Generating support bundle"
#          generateSupportBundle "$testname"
          if [ "$testname" != "install" ] ; then
              if [ "$on_fail" == "restart" ] ; then
                  echo "Attempting to continue by cleaning up and restarting mayastor pods........"
                  if ! runGo "tools/restart" "restart.go" ; then
                      echo "\"restart\" failed"
                      exit $EXITV_FAILED
                  fi
              elif [ "$on_fail" == "reinstall" ] ; then
                  echo "Attempting to continue by cleaning up and re-installing........"
                  runGo "tools/cleanup" "cleanup.go"
                  if ! runGoTest $TESTDIR "tests/uninstall" "uninstall"; then
                      echo "uninstall failed, abandoning attempt to continue"
                      exit $EXITV_FAILED
                  fi
                  if ! runGoTest $TESTDIR "tests/install" "install"; then
                      echo "(re)install failed, abandoning attempt to continue"
                      exit $EXITV_FAILED
                  fi
              else
                  break
              fi
          else
              break
          fi
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
if contains "$tests" "uninstall" ; then
    if ! runGoTest $TESTDIR "tests/uninstall" "uninstall" ; then
        echo "Test \"uninstall\" FAILED!"
        test_failed=1
        emitLogs "uninstall"
    else
        # if [ "$coverage" == "true" ]; then
        #     if ! "$SCRIPTDIR/remote-coverage-files.py" --get --path "$coveragedir" --identity_file "$ssh_identity" ; then
        #         echo "Failed to retrieve coverage files"
        #     fi
        # fi
        if  [ "$test_failed" -ne 0 ] ; then
            # tests failed, but uninstall was successful
            # so cluster is reusable
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
