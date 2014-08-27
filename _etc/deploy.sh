#!/bin/bash -e

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

scp _etc/server.go martin@37.59.112.124:ahouhpuc/
ssh -T martin@37.59.112.124 <<EOF
go build -o ahouhpuc/server ahouhpuc/server.go
EOF

scp _etc/nginx.conf root@37.59.112.124:/etc/nginx/sites-enabled/ahouhpuc.conf

ssh -T root@37.59.112.124 <<EOF
/etc/init.d/nginx restart
/etc/init.d/ahouhpuc restart
EOF
