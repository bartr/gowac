#!/bin/sh
service ssh start

# git the latest code
cd /root/go/src/gowac/app
git pull
cd /root

while :
do
  go build gowac/app -o ./gowac

  ./gowac
  sleep .1
  rm -f gowac
done
