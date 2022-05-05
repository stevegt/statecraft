#!/bin/bash

# used during development

pid=$1
wd=$(pwd)

retry() {
    cd $wd
    exec $0 $1
}

echo ==========
set -x

inotifywait -r -e modify .
sleep 1
[ -n "$pid" ] && kill $pid
go test -v || retry
go build || retry

# sc=$(realpath statecraft)
cd example/stoplight
# $sc car.statecraft car.dot || retry
# $sc stoplight.statecraft stoplight.dot || retry
go generate && go run . &
pid=$!
sleep 3

retry $pid
