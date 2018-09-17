#!/bin/sh
service ssh start

# git the latest code
cd /root/gowac
git pull

while :
do
  go build

  ./gowac
  sleep .1
  rm -f gowac
done
