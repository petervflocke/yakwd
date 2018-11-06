#!/bin/sh
SERVICE=/usr/local/bin/yakwd
if pgrep -f "$SERVICE" >/dev/null 2>&1 ; then
  echo "$SERVICE is RUNNING"
else
/etc/init.d/framework stop
/usr/bin/lipc-set-prop -- com.lab126.powerd preventScreenSaver 1
/usr/sbin/eips -s 5
/usr/local/bin/yakwd
/usr/bin/lipc-set-prop -- com.lab126.powerd preventScreenSaver 0
/etc/init.d/framework start
fi
