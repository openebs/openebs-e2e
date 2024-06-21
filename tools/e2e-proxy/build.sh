#!/usr/bin/env bash
# This script builds the image in the host docker registry with
# the latest tag
# if a registry is specified then it pushes the image to that
# registry with the configure TAG value see below
# if --also-tag-as-latest option is passed then the image is
# also pushed to the registry with the tagged as latest
# This works for legacy test runs and also other test frameworks
# as long as we do not make breaking changes.
set -e
IMAGE="mayadata/e2e-proxy"
TAG="v1.0.0"
registry=""
tag_as_latest=""

while [ "$#" -gt 0 ] ; do
    case "$1" in
        --registry)
            shift
            registry=$1
            ;;
        --also-tag-as-latest)
            tag_as_latest="Y"
            ;;
        *)
            echo "Unknown option: $1"
            help
            exit 1
            ;;
    esac
    shift
done

if docker build -t ${IMAGE} --build-arg GO_VERSION=1.19.3 . ; then
    if [ "${registry}" != "" ]; then
        if [ "${tag_as_latest}" == "Y" ];  then
            echo "tagging as latest and pushing image"
            docker tag ${IMAGE} $registry/${IMAGE}
            docker push $registry/${IMAGE}
        fi
        if [ "${TAG}" != "" ];  then
            echo "tagging as ${TAG} and pushing image"
            docker tag ${IMAGE} $registry/${IMAGE}:${TAG}
            docker push $registry/${IMAGE}:${TAG}
        else
            echo "TAG was not defined - image not retagged and pushed"
        fi
    fi
else
    exit 1
fi
