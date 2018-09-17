#!/bin/sh
service ssh start

# if you don't run from the app directory, make sure to copy www/
cd /root/go/src/gowac/app

# git the latest code
git pull

while :
do
  go build

  ./app
  sleep .1
  rm -f app
done
