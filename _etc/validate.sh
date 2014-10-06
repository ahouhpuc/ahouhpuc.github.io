#!/bin/bash

error=0

for file in _site/*.html; do
  if [[ $file != "_site/google6fc913931bb74ee6.html" ]]; then
    out=$(curl --silent --data-binary @$file --header "Content-Type: text/html" "http://html5.validator.nu/?out=gnu")
    if [[ -n $out ]]; then
      error=1
      echo $file
      echo $out
    fi
  fi
done

if [[ $error == 1 ]]; then
  echo "ERROR, one or more files not valid."
  exit 1
fi
