#!/bin/sh
 
COV_PARTIAL_FILE=profile.cov.out
COV_FILE=profile-all.cov.out
COV_MODE=count
 
HEADER="mode: $COV_MODE"
 
if which gxargs > /dev/null; then
  XARGS=gxargs
else
  XARGS=xargs
fi
 
XARGS_ARG=DIR

CMD_COV="go test $XARGS_ARG -covermode=$COV_MODE -coverprofile=$XARGS_ARG/$COV_PARTIAL_FILE"
CMD_CHECKDIR="ls $XARGS_ARG/$COV_PARTIAL_FILE > /dev/null 2>&1"
CMD_CONCATCOV="cat $XARGS_ARG/$COV_PARTIAL_FILE | grep -v '$HEADER' >> $COV_FILE"
 
GOVERALLS=$(echo $GOPATH | tr ':' '\n' | head -n 1)/bin/goveralls
 
echo $HEADER > $COV_FILE
 
find ./*  -maxdepth 10 -name '*.go' -print0 | $XARGS -0 -L1 dirname | uniq | grep -v "/_" | $XARGS -I$XARGS_ARG sh -c "$CMD_COV && ($CMD_CHECKDIR && $CMD_CONCATCOV) || exit 0"
 
$GOVERALLS -coverprofile=$COV_FILE -service=travis-ci
