#!/bin/sh

set -ex

MAYA_API_SERVER_NETWORK=$1

CONTAINER_IP_ADDR=$(ip -4 addr show scope global dev "${MAYA_API_SERVER_NETWORK}" | grep inet | awk '{print $2}' | cut -d / -f 1)

# Start apiserver service
exec /usr/local/bin/maya-apiserver.test start --bind="${CONTAINER_IP_ADDR}" -test.coverprofile=/tmp/maya-apiserver-coverage.cov 1>&2
