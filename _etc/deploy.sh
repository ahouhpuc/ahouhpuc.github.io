#!/bin/bash -xe

BUILD_DIR=`mktemp -d -t ahouhpucXXX`
BUILD_SHA=`git rev-parse HEAD`

TAR="tar"
if [ -x "$(command -v gnutar)" ]; then
  TAR="gnutar"
fi

git archive --format=tar HEAD | (cd $BUILD_DIR && $TAR xf -)
cd $BUILD_DIR
jekyll build
$TAR czf _site.tgz _site/
scp _site.tgz martin@ahouhpuc.fr:ahouhpuc/
ssh -T martin@ahouhpuc.fr <<EOF
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
cd _etc && GOOS=linux GOARCH=amd64 go build -o server
cd ..
ssh -T root@ahouhpuc.fr <<EOF
systemctl stop ahouhpuc.service
EOF
sleep 1
scp _etc/server martin@ahouhpuc.fr:ahouhpuc/server
rm _etc/server

ssh -T root@ahouhpuc.fr <<EOF
setcap cap_net_bind_service=+ep /home/martin/ahouhpuc/server
systemctl start ahouhpuc.service
sleep 1
systemctl status ahouhpuc.service
EOF
