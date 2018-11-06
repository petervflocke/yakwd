#!/bin/sh
SERVICE=yakwd
if pgrep "$SERVICE" >/dev/null 2>&1 ; then
  echo "$SERVICE is RUNNING"
else
/etc/init.d/framework stop
/usr/bin/lipc-set-prop -- com.lab126.powerd preventScreenSaver 1
/usr/sbin/eips -s 5
./yakwd
/usr/bin/lipc-set-prop -- com.lab126.powerd preventScreenSaver 0
/etc/init.d/framework start
fi
