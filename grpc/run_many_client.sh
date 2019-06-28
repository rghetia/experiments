#!/bin/bash

num=$2
client=$1
name=`basename $client`
echo $name
i=0
while [ $i -le $num ] ; do
    log="/tmp/${name}_$i.log"
    rm -rf $log
    $client >& $log &
    pid=$!
    pids[${i}]=$!
    i=`expr $i + 1`
done
echo ${pids[@]:0}
