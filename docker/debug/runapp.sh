#!/bin/sh
service ssh start

# git the latest code
cd /root/gowac
git pull
cd /root

while :
do
  go build -o app gowac

  ./app
  sleep .1
  rm -f app
done
