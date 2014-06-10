#!/bin/bash -ex

BUILD_DIR=`mktemp -d -t ahouhpuc`
BUILD_SHA=`git rev-parse HEAD`
git archive --format=tar HEAD | (cd $BUILD_DIR && tar xf -)
cd $BUILD_DIR
jekyll build
gnutar czf _site.tgz _site/
scp _site.tgz martin@martinottenwaelter.fr:ahouhpuc/
ssh -T martin@martinottenwaelter.fr <<EOF
cd ahouhpuc
if [ -d $BUILD_SHA ]; then
  echo "$BUILD_SHA already exists. Aborting."
  exit 1
fi
tar xzf _site.tgz
mv _site $BUILD_SHA
convmv -r -f utf8 -t utf8 --nfc --notest --replace $BUILD_SHA
rm _site.tgz
rm -f current && ln -s $BUILD_SHA current
EOF
rm _site.tgz

scp _etc/nginx.conf martin@martinottenwaelter.fr:ahouhpuc/nginx.conf
ssh -T root@martinottenwaelter.fr <<EOF
/etc/init.d/nginx restart
EOF
