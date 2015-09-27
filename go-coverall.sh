#!/bin/bash

COV_FILE=profile.cov.out
COV_TMP_FILE=profile_tmp.cov
ERROR=""

echo "mode: count" > $COV_FILE

for pkg in `go list ./...`
do
    touch $COV_TMP_FILE
    go test -covermode=count -coverprofile=$COV_TMP_FILE $pkg || ERROR="Error testing $pkg"
    tail -n +2 $COV_TMP_FILE >> $COV_FILE || (echo "Unable to append coverage for $pkg" && exit 1)
done

if [ ! -z "$ERROR" ]
then
    echo "Encountered error, last error was: $ERROR"
    exit 1
fi

GOVERALLS=$(echo $GOPATH | tr ':' '\n' | head -n 1)/bin/goveralls
$GOVERALLS -coverprofile=COV_FILE -service=travis-ci

