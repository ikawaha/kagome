#!/bin/bash

#export GOPATH=`pwd`:$GOPATH
echo "mode: count" > profile.cov

ERROR=""

for pkg in `go list ./...`
do
    touch profile_tmp.cov
    go test -v -covermode=count -coverprofile=profile_tmp.cov $pkg || ERROR="Error testing $pkg"
    tail -n +2 profile_tmp.cov >> profile.cov || (echo "Unable to append coverage for $pkg" && exit 1)
done

if [ ! -z "$ERROR" ]
then
    echo "Encountered error, last error was: $ERROR"
    exit 1
fi

GOVERALLS=$(echo $GOPATH | tr ':' '\n' | head -n 1)/bin/goveralls
$GOVERALLS -coverprofile=profile.cov -service=travis-ci

