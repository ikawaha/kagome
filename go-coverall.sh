#!/bin/bash

COV_FILE=coverage.txt
COV_TMP_FILE=coverage_tmp.cov
ERROR=""

echo "mode: count" > $COV_FILE

for pkg in `go list ./...`
do
    touch $COV_TMP_FILE
    go test -covermode=count -coverprofile=$COV_TMP_FILE $pkg || ERROR="Error testing $pkg"
    tail -n +2 $COV_TMP_FILE >> $COV_FILE || (echo "Unable to append coverage for $pkg" && exit 1)
done

rm $COV_TMP_FILE

if [ ! -z "$ERROR" ]
then
    echo "Encountered error, last error was: $ERROR"
    exit 1
fi

GOPATH=$HOME/gopath
GOVERALLS=$(echo $GOPATH | tr ':' '\n' | head -n 1)/bin/goveralls
$GOVERALLS -coverprofile=$COV_FILE -service=travis-ci

