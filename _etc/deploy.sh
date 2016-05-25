#!/bin/bash -xe

BUILD_DIR=`mktemp -d -t ahouhpuc`
BUILD_SHA=`git rev-parse HEAD`
git archive --format=tar HEAD | (cd $BUILD_DIR && tar xf -)
cd $BUILD_DIR
jekyll build
gnutar czf _site.tgz _site/
scp _site.tgz martin@37.59.112.124:ahouhpuc/
ssh -T martin@37.59.112.124 <<EOF
cd ahouhpuc
if [ ! -d $BUILD_SHA ]; then
  tar xzf _site.tgz
  mv _site $BUILD_SHA
  convmv -r -f utf8 -t utf8 --nfc --notest --replace $BUILD_SHA
  rm -f current && ln -s $BUILD_SHA current
fi
rm _site.tgz
EOF
rm _site.tgz

rm -f _etc/server
GOPATH=$(pwd)/_etc GOOS=linux GOARCH=amd64 go build -o _etc/server _etc/*.go
ssh -T root@37.59.112.124 <<EOF
/etc/init.d/ahouhpuc stop
EOF
sleep 1
scp _etc/server martin@37.59.112.124:ahouhpuc/server
rm _etc/server

ssh -T root@37.59.112.124 <<EOF
setcap cap_net_bind_service=+ep /home/martin/ahouhpuc/server
/etc/init.d/ahouhpuc start
sleep 1
/etc/init.d/ahouhpuc status
EOF
