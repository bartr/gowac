#!/bin/sh
service ssh start

# git the latest code
cd /root/gowac
git pull
cd /root

while :
do
  go build gowac

  ./gowac
  sleep .1
  rm -f gowac
done