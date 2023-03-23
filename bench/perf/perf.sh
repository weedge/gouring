#!/bin/bash
set -x

make build

N=10000000
PERF_OPTS="-n $N -noti 5000000"

if [ $1 = "sqpoll-stat-cost" ];then
    ./perf -sqpoll $PERF_OPTS -pprofCpu pprof-sqpoll.cpu &
    echo "sqpoll pid: $!"
    sudo perf stat -a -ddd -p $!
    exit
elif [ $1 = "stat-cost" ];then
    ./perf $PERF_OPTS -pprofCpu pprof-nonsqpoll.cpu &
    echo "pid: $!"
    sudo perf stat -a -ddd -p $!
    exit
fi

if [ $1 = "sqpoll-syscall" ];then
    ./perf -sqpoll $PERF_OPTS -pprofCpu pprof-sqpoll.cpu &
    echo "sqpoll pid: $!"
    sudo perf trace -s -p $!
    exit
elif [ $1 = "syscall" ];then
    ./perf $PERF_OPTS -pprofCpu pprof-nonsqpoll.cpu &
    echo "pid: $!"
    sudo perf trace -s -p $!
    exit
fi

if [ $1 = "sqpoll-stat-event" ];then
    ./perf -sqpoll $PERF_OPTS -pprofCpu pprof-sqpoll.cpu &
    echo "sqpoll pid: $!"
    sudo perf stat -e io_uring:* -p $!
    exit
elif [ $1 = "stat-event" ];then
    ./perf $PERF_OPTS -pprofCpu pprof-nonsqpoll.cpu &
    echo "pid: $!"
    sudo perf stat -e io_uring:* -p $!
    exit
fi


