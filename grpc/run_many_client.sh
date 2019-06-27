#!/bin/bash

go build -o /tmp/client helloworld_client/main.go

num=$1
i=0
while [ $i -le $num ] ; do
    /tmp/client &
    pid=$!
    pids[${i}]=$!
    i=`expr $i + 1`
done
echo ${pids[@]:0}
