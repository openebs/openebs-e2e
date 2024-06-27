#!/usr/bin/env bash

# This script makes the best attempt to dump stuff
# so ignore fails and keep paddling.
# set -e

help() {
  cat <<EOF
This script generates a meta.json file with release info and other details.

Usage: $0 [OPTIONS]

Options:
  --destdir <path>   Location to store log files
  --release  <string> Release name to be stored in meta.json (ex: develop/2.x/3.x)
  --platform <string> Test Platform to be stored in meta.json (ex: hetzner/reading)
  --bundle <string> Install Bundle name to be stored in meta.json
EOF
}



# args destdir previous namespace podname containername
# destdir == "" -> stdout
# previous == "" -> current logs else -> previous logs
function createMetaJson {
    destdir=$1
    release=$2
    platform=$3
    bundle=$4

    if [ -z "$destdir" ] ; then
        echo "ERROR calling createMetaJson as destination directory not present"
        return
    fi

    if  [ -z "$release" ] ; then
        echo "ERROR calling createMetaJson as release information not present"
        return
    fi

    if  [ -z "$platform" ] ; then
        echo "ERROR calling createMetaJson as platform information not present"
        return
    fi

    if  [ -z "$bundle" ] ; then
        echo "ERROR calling createMetaJson as bundle information not present"
        return
    fi


    if [["$destdir" =~ [^a-zA-Z0-9] ]]; then
        echo "ERROR : directory name has non alphanumeric character"
	return
    fi

    mkdir -p "$destdir"
    JSON_FMT='{"release":"%s","platform":"%s","bundle":"%s"}\n'
    JSON_STRING=$(printf "$JSON_FMT" "$release" "$platform" "$bundle")
    echo "$JSON_STRING"
    echo "$JSON_STRING" > "$destdir/meta.json"
}

destdir=
release=
platform=
bundle=

# Parse arguments
while [ "$#" -gt 0 ]; do
  case "$1" in
    -d|--destdir)
      shift
      destdir="$1"
      ;;
    -r|--release)
      shift
      release="$1"
      ;;
    -p|--platform)
      shift
      platform="$1"
      ;;
    -b|--bundle)
      shift
      bundle="$1"
      ;;

    *)
      echo "Unknown option: $1"
      help
      exit 1
      ;;
  esac
  shift
done

createMetaJson "$destdir" "$release" "$platform" "$bundle"
