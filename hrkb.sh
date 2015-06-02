#!/bin/bash
# /etc/init.d/hrkb

# Path to hrkb 
hpath="/home/iworld/Development/goProjects/beego_test/src/github.com/fm/hrkb"
case "$1" in
  start)
    rm -f $hpath/.hrkb.pid
    nohup $hpath/./hrkb > /dev/null 2>&1 & echo $! > $hpath/hrkb.pid 
    echo "Service started via nohup, pid: " $!
    exit 0
    ;;
  stop)
    kill `cat $hpath/hrkb.pid`
    echo "Service stopped"
    exit 0
    ;;
  restart)
    kill `cat $hpath/hrkb.pid`
    echo "Service stopped"
    rm -f $hpath/.hrkb.pid
    nohup $hpath/./hrkb > /dev/null 2>&1 & echo $! > $hpath/hrkb.pid 
    echo "Service started via nohup, pid: " $!
    exit 0 
    ;;

  *)
    echo "Usage: /etc/init.d/hrkb {start|stop|restart}"
    exit 0 
    ;;
esac

exit 0
