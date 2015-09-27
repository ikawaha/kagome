#!/bin/sh

function die() {
    echo $*
    exit 1
}

export GOPATH=`pwd`:$GOPATH
echo "mode: count" > profile.cov

ERROR=""

for pkg in `go list ./...`
do
    touch profile_tmp.cov
    go test -v -covermode=count -coverprofile=profile_tmp.cov $pkg || ERROR="Error testing $pkg"
    tail -n +2 profile_tmp.cov >> profile.cov || die "Unable to append coverage for $pkg"
done

if [ ! -z "$ERROR" ]
then
    die "Encountered error, last error was: $ERROR"
fi

GOVERALLS=$(echo $GOPATH | tr ':' '\n' | head -n 1)/bin/goveralls
$GOVERALLS -coverprofile=profile.cov -service=travis-ci

