#!/usr/bin/env bash

SCRIPTDIR=$(dirname "$(realpath "$0")")
GOSRCDIRAPPS=$(realpath "$SCRIPTDIR/../apps")
GOSRCDIRCOMMON=$(realpath "$SCRIPTDIR/../common")
GOSRCDIRE2EAGENT=$(realpath "$SCRIPTDIR/../tools/e2e-agent")
GOSRCDIRE2EPROXY=$(realpath "$SCRIPTDIR/../tools/e2e-proxy")
GOSRCDIRSRC=$(realpath "$SCRIPTDIR/../src")
exitv=0

reformat=1
if [ "$1" == "-n" ]; then
    reformat=0
    shift
fi

if [ -n "$*" ]; then
    # go fmt
    unformatted=$(gofmt -l "$@")
fi

if [ $reformat -ne 0 ]; then
    for file in $unformatted; do
        gofmt -w "$file"
        echo "Reformatted $file"
    done
else
    if [ -n "$unformatted" ]; then
        printf "Please run\n\tgofmt -w %s\n" "$unformatted"
        exitv=1
    fi
fi

if golangci-lint > /dev/null 2>&1 ; then
    cd "$GOSRCDIRAPPS" || exit 1
    echo "## linting apps ##"
    if ! golangci-lint run -v --allow-parallel-runners ; then
        exitv=1
    fi
    cd "$GOSRCDIRCOMMON" || exit 1
    echo ""
    echo "## linting common ##"
    if ! golangci-lint run -v --allow-parallel-runners ; then
        exitv=1
    fi
    cd "$GOSRCDIRE2EAGENT" || exit 1
    echo ""
    echo "## linting e2e-agent ##"
    if ! golangci-lint run -v --allow-parallel-runners ; then
        exitv=1
    fi
    cd "$GOSRCDIRE2EPROXY" || exit 1
    echo ""
    echo "## linting e2e-proxy ##"
    if ! golangci-lint run -v --allow-parallel-runners ; then
        exitv=1
    fi

    echo ""
    echo "## linting src  ##"
    cd "$GOSRCDIRSRC" || exit 1
    if ! golangci-lint run -v --allow-parallel-runners ; then
        exitv=1
    fi

else
    exitv=1
    echo "Please install golangci-lint"
fi

exit $exitv
