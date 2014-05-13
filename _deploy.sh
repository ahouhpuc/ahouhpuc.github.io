#!/bin/bash

jekyll build
COPYFILE_DISABLE=true tar czf _site.tgz _site/
ssh martin@martinottenwaelter.fr "rm -rf public/ahouhpuc"
scp _site.tgz martin@martinottenwaelter.fr:public/
ssh martin@martinottenwaelter.fr "cd public && tar xzf _site.tgz && mv _site ahouhpuc && convmv -r -f utf8 -t utf8 --nfc --notest --replace ahouhpuc && rm _site.tgz"
rm _site.tgz
