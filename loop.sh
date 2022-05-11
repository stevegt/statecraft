#!/bin/bash

# used during development

wd=$(pwd)

while true
do
    set +x
    echo ==========
    set -x
    cd $wd
    inotifywait -r -e modify .
    sleep 1
    [ -n "$pid" ] && kill $pid
    go test -v || continue 
    go build || continue

    # sc=$(realpath statecraft)
    cd example/stoplight
    # $sc car.statecraft car.dot || retry
    # $sc stoplight.statecraft stoplight.dot || retry
    go generate || continue
    go run . &
    pid=$!
done
