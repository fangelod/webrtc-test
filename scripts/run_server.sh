#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "Building server"
go build -race -o $GOPATH/bin/webrtc-test github.com/fangelod/webrtc-test/cmd/server

source $DIR/../configs/webrtc_test.env

$GOPATH/bin/webrtc-test

exit $?